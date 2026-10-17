package reviewchecks

import (
	"bytes"
	"encoding/json"
	"ensemble"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// Only synthetic local markers; no real provider payload or capability row.
const reviewReturnedModel = "review-r1-returned"

var reviewFrom = ensemble.Provenance{Vendor: ensemble.Gemini, Surface: ensemble.GenerateContent, Model: reviewReturnedModel}

type reviewBoundary struct {
	name   string
	native []string
	frames [][]string
	want   []ensemble.Part
}

func reviewText(s, sig string) ensemble.Part {
	p := ensemble.Part{Type: "text", Text: s}
	if sig != "" {
		p.From = reviewFrom
		p.Opaque = json.RawMessage(fmt.Sprintf("%q", sig))
	}
	return p
}
func reviewOpaque(raw string) ensemble.Part {
	return ensemble.Part{Type: "opaque", From: reviewFrom, Data: json.RawMessage(raw)}
}
func reviewNormalize(x any) any {
	b, _ := json.Marshal(x)
	var v any
	_ = json.Unmarshal(b, &v)
	return v
}
func reviewWire(parts []string, last bool) string {
	end := ""
	if last {
		end = `,"finishReason":"STOP"`
	}
	return `{"modelVersion":"` + reviewReturnedModel + `","candidates":[{"content":{"parts":[` + strings.Join(parts, ",") + `]}` + end + `}],"usageMetadata":{"promptTokenCount":7,"cachedContentTokenCount":2,"candidatesTokenCount":5,"thoughtsTokenCount":1}}`
}
func reviewAgent(t *testing.T, root *ensemble.Ensemble, v ensemble.Vendor, resolved string) ensemble.Agent {
	t.Helper()
	model := "gemini-3.8-flash"
	if v == ensemble.Anthropic {
		model = "claude-sonnet-4-6"
	}
	a, e := root.NewAgent(ensemble.Config{Vendor: v, Model: model, ResolvedModel: resolved, APIKey: "review-fixture-only", System: "review", MaxTokens: 128, DisableStreaming: true})
	if e != nil {
		t.Fatal(e)
	}
	return a
}
func TestReviewerGeminiExplicitBoundariesAndReplay(t *testing.T) {
	cases := []reviewBoundary{
		{name: "opaque-first", native: []string{`{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`, `{"text":"visible"}`}, want: []ensemble.Part{reviewOpaque(`{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`), reviewText("visible", "")}},
		{name: "signed-then-unsigned", native: []string{`{"text":"signed ","thoughtSignature":"REVIEW_SIG_ONE"}`, `{"text":"plain"}`}, want: []ensemble.Part{reviewText("signed ", "REVIEW_SIG_ONE"), reviewText("plain", "")}},
		{name: "two-signed", native: []string{`{"text":"first ","thoughtSignature":"REVIEW_SIG_ONE"}`, `{"text":"second","thoughtSignature":"REVIEW_SIG_TWO"}`}, want: []ensemble.Part{reviewText("first ", "REVIEW_SIG_ONE"), reviewText("second", "REVIEW_SIG_TWO")}},
		{name: "thought-opaque-text", native: []string{`{"text":"review-hidden","thought":true,"thoughtSignature":"REVIEW_SIG_THOUGHT"}`, `{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`, `{"text":"visible"}`}, want: []ensemble.Part{reviewOpaque(`{"text":"review-hidden","thought":true,"thoughtSignature":"REVIEW_SIG_THOUGHT"}`), reviewOpaque(`{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`), reviewText("visible", "")}},
		{name: "text-continuation-and-closure", native: []string{`{"text":"joined","thoughtSignature":"REVIEW_SIG_TAIL"}`, `{"text":" separate"}`}, frames: [][]string{{`{"text":"join"}`}, {`{"text":"ed"}`}, {`{"thoughtSignature":"REVIEW_SIG_TAIL"}`}, {`{"text":" separate"}`}}, want: []ensemble.Part{reviewText("joined", "REVIEW_SIG_TAIL"), reviewText(" separate", "")}},
		{name: "thought-continuation", native: []string{`{"text":"review-hidden","thought":true,"thoughtSignature":"REVIEW_SIG_THOUGHT"}`, `{"text":"visible"}`}, frames: [][]string{{`{"text":"review-","thought":true}`}, {`{"text":"hidden","thought":true}`}, {`{"thoughtSignature":"REVIEW_SIG_THOUGHT"}`}, {`{"text":"visible"}`}}, want: []ensemble.Part{reviewOpaque(`{"text":"review-hidden","thought":true,"thoughtSignature":"REVIEW_SIG_THOUGHT"}`), reviewText("visible", "")}},
		{name: "tool-signature-continuation", native: []string{`{"text":"explanation"}`, `{"functionCall":{"id":"review-call","name":"review_lookup","args":{"n":1}},"thoughtSignature":"REVIEW_SIG_TOOL"}`}, frames: [][]string{{`{"text":"explanation"}`, `{"functionCall":{"id":"review-call","name":"review_lookup","args":{"n":1}}}`}, {`{"thoughtSignature":"REVIEW_SIG_TOOL"}`}}, want: []ensemble.Part{reviewText("explanation", ""), {Type: "tool_call", CallID: "review-call", Name: "review_lookup", Args: json.RawMessage(`{"n":1}`), From: reviewFrom, Opaque: json.RawMessage(`"REVIEW_SIG_TOOL"`)}}},
		{name: "explicit-unsigned-parts", native: []string{`{"text":"first "}`, `{"text":"second"}`}, frames: [][]string{{`{"text":"first "}`, `{"text":"second"}`}}, want: []ensemble.Part{reviewText("first ", ""), reviewText("second", "")}},
		{name: "explicit-opaque-after-open-text", native: []string{`{"text":"visible"}`, `{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`}, frames: [][]string{{`{"text":"visible"}`, `{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`}}, want: []ensemble.Part{reviewText("visible", ""), reviewOpaque(`{"thoughtSignature":"REVIEW_SIG_OPAQUE"}`)}},
	}
	for _, c := range cases {
		layouts := map[string][][]string{"continuations": c.frames}
		if c.frames == nil {
			split := [][]string{}
			for _, p := range c.native {
				split = append(split, []string{p})
			}
			layouts = map[string][][]string{"one-frame": {c.native}, "split-frames": split}
		}
		for layout, frames := range layouts {
			t.Run(c.name+"/"+layout, func(t *testing.T) {
				root := ensemble.New(log.New(io.Discard, "", 0))
				t.Cleanup(func() { _ = root.Shutdown() })
				a := reviewAgent(t, root, ensemble.Gemini, reviewReturnedModel)
				want := ensemble.ResponseEvent{From: reviewFrom, Parts: c.want, Usage: ensemble.Usage{Input: 5, CacheRead: 2, Output: 6}}
				var got ensemble.ResponseEvent
				for repeat := 0; repeat < 2; repeat++ {
					for _, stream := range []bool{false, true} {
						body := reviewWire(c.native, true)
						ct := "application/json"
						if stream {
							ct = "text/event-stream"
							var b strings.Builder
							for i, p := range frames {
								fmt.Fprintf(&b, "data: %s\n\n", reviewWire(p, i == len(frames)-1))
							}
							body = b.String()
						}
						events := 0
						finals := map[uint64]ensemble.Part{}
						delta := map[uint64]map[ensemble.DeltaKind]string{}
						err := a.Engine().Parse(&http.Response{StatusCode: 200, Header: http.Header{"Content-Type": []string{ct}}, Body: io.NopCloser(strings.NewReader(body))}, ensemble.StreamCallbacks{
							OnEvent: func(e ensemble.Event) {
								if e.Response != nil {
									got = *e.Response
									events++
								}
							},
							OnPartFinal: func(id uint64, p ensemble.Part) {
								if _, exists := finals[id]; exists {
									t.Errorf("duplicate final ID %d", id)
								}
								finals[id] = p
							},
							OnDelta: func(id uint64, k ensemble.DeltaKind, s string) {
								if strings.Contains(s, "REVIEW_SIG_") {
									t.Error("opaque signature exposed as delta")
								}
								if delta[id] == nil {
									delta[id] = map[ensemble.DeltaKind]string{}
								}
								delta[id][k] += s
							},
						})
						if err != nil {
							t.Fatalf("stream=%v parse: %v", stream, err)
						}
						if events != 1 || !reflect.DeepEqual(reviewNormalize(got), reviewNormalize(want)) {
							t.Fatalf("stream=%v canonical facts mismatch\ngot=%s\nwant=%s", stream, mustJSONReview(got), mustJSONReview(want))
						}
						if len(finals) != len(c.want) {
							t.Fatalf("final count %d want %d", len(finals), len(c.want))
						}
						for i, p := range c.want {
							id := uint64(i)
							f, ok := finals[id]
							if !ok || !reflect.DeepEqual(reviewNormalize(f), reviewNormalize(p)) {
								t.Fatalf("local ID %d final mismatch", id)
							}
							expect := map[ensemble.DeltaKind]string{}
							switch p.Type {
							case "text":
								expect[ensemble.DeltaText] = p.Text
							case "tool_call":
								expect[ensemble.DeltaToolCall] = p.Name + string(p.Args)
							case "opaque":
								var native map[string]any
								_ = json.Unmarshal(p.Data, &native)
								if native["thought"] == true {
									expect[ensemble.DeltaThinking], _ = native["text"].(string)
								}
							}
							d := delta[id]
							if d == nil {
								d = map[ensemble.DeltaKind]string{}
							}
							if !reflect.DeepEqual(d, expect) {
								t.Fatalf("stream=%v ID=%d delta=%v want=%v", stream, id, d, expect)
							}
						}
						for id := range delta {
							if _, ok := finals[id]; !ok {
								t.Fatalf("orphan delta %d", id)
							}
						}
					}
				}
				if e := a.Log().Append(ensemble.Event{Type: "response_ended", Response: &got}); e != nil {
					t.Fatal(e)
				}
				ctx := a.Log().Context()
				wire, e := a.Engine().Render(&ctx)
				if e != nil {
					t.Fatal(e)
				}
				var doc struct {
					Contents []struct {
						Role  string            `json:"role"`
						Parts []json.RawMessage `json:"parts"`
					} `json:"contents"`
				}
				if e = json.Unmarshal(wire, &doc); e != nil {
					t.Fatal(e)
				}
				native := []json.RawMessage{}
				for _, m := range doc.Contents {
					if m.Role == "model" {
						native = append(native, m.Parts...)
					}
				}
				expectNative := []json.RawMessage{}
				for _, p := range c.native {
					expectNative = append(expectNative, json.RawMessage(p))
				}
				if !reflect.DeepEqual(reviewNormalize(native), reviewNormalize(expectNative)) {
					t.Fatalf("same-source replay differs: %s", wire)
				}
				for _, gate := range []string{"model", "vendor", "surface"} {
					t.Run("drop-"+gate, func(t *testing.T) {
						v := ensemble.Gemini
						model := reviewReturnedModel
						if gate == "model" {
							model = "review-other-model"
						}
						if gate == "vendor" {
							v = ensemble.Anthropic
						}
						target := reviewAgent(t, root, v, model)
						copyCtx := a.Log().Context()
						if gate == "surface" {
							for i := range copyCtx.Dialogue {
								for j := range copyCtx.Dialogue[i].Parts {
									p := &copyCtx.Dialogue[i].Parts[j]
									if p.From.Vendor == ensemble.Gemini {
										p.From.Surface = ensemble.Messages
									}
								}
							}
						}
						out, e := target.Engine().Render(&copyCtx)
						if e != nil {
							t.Fatal(e)
						}
						if bytes.Contains(out, []byte("REVIEW_SIG_")) || bytes.Contains(out, []byte("review-hidden")) {
							t.Fatal("private replay crossed source boundary")
						}
						for _, p := range c.want {
							if p.Type == "text" && !bytes.Contains(out, []byte(p.Text)) {
								t.Fatalf("ordinary text disappeared across %s", gate)
							}
						}
					})
				}
			})
		}
	}
}
func mustJSONReview(v any) string { b, _ := json.Marshal(v); return string(b) }
