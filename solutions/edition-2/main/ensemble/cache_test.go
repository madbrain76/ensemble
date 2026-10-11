package ensemble_test

import (
	"bytes"
	"context"
	"encoding/json"
	"ensemble/ensemble"
	"ensemble/internal/common"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestCacheLensMainBoundary drives real HTTP and owned diagnostics. Auxiliary
// prompts must not replace the preceding main-conversation request.
func TestCacheLensMainBoundary(t *testing.T) {
	for _, vendor := range []ensemble.Vendor{ensemble.Anthropic, ensemble.OpenAI, ensemble.Gemini} {
		t.Run(vendor.String(), func(t *testing.T) {
			peer := &recallPeer{vendor: vendor, judgment: "[]"}
			server := httptest.NewServer(http.HandlerFunc(peer.serve))
			defer server.Close()
			var logs bytes.Buffer
			a, err := ensemble.New(&logs).NewAgent(ensemble.Config{Vendor: vendor, Model: "fake-model", BaseURL: server.URL, DataDir: t.TempDir()})
			if err != nil {
				t.Fatal(err)
			}
			defer a.Shutdown()
			for _, q := range []string{"red", "green"} {
				if _, err := a.Ask(context.Background(), q); err != nil {
					t.Fatal(err)
				}
			}
			lens := a.Engine().CacheLens()
			if lens == nil || lens.Engine() != a.Engine() {
				t.Fatal("lens ownership")
			}
			current := filepath.Join(a.Config().DataDir, "cache-current.txt")
			before, err := os.ReadFile(current)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := os.Stat(filepath.Join(a.Config().DataDir, "cache-prior.txt")); err != nil {
				t.Fatal(err)
			}
			if !strings.Contains(logs.String(), ":appended") || strings.Contains(logs.String(), "cache prefix changed") || !strings.Contains(logs.String(), "cache grant:") {
				t.Fatalf("verdicts: %s", logs.String())
			}
			a.Engine().Render(a.History().Context())
			if _, err := a.Engine().Judge(context.Background(), "zero-based candidate indices"); err != nil {
				t.Fatal(err)
			}
			after, _ := os.ReadFile(current)
			if !bytes.Equal(before, after) {
				t.Fatal("auxiliary request rotated lens")
			}
			body, _ := a.Engine().Render(a.History().Context())
			body = bytes.ReplaceAll(body, []byte("helpful assistant"), []byte("different assistant"))
			lens.ObserveRequest(vendor.String()+"/fake-model", body)
			if !strings.Contains(logs.String(), "cache prefix changed") {
				t.Fatal("stable-prefix rewrite unreported")
			}
		})
	}
}

// TestCacheLensSerializedMarkers checks the actual renderer in compact and
// indented forms; marker stripping must preserve unrelated escaping and bytes.
func TestCacheLensSerializedMarkers(t *testing.T) {
	a, err := ensemble.New(io.Discard).NewAgent(ensemble.Config{Model: "fake", DataDir: t.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	defer a.Shutdown()
	body, err := a.Engine().Render(&common.Context{Dialogue: []common.Entry{{Actor: common.Human, Parts: []common.Part{{Type: "text", Text: "keep <escaped> bytes"}}}}})
	if err != nil {
		t.Fatal(err)
	}
	for _, pretty := range []bool{false, true} {
		marked := bytes.Replace(body, []byte(`"type":"text"`), []byte(`"type":"text","cache_control":{"type":"ephemeral"}`), 1)
		plain := body
		if pretty {
			var x, y bytes.Buffer
			json.Indent(&x, marked, "", "  ")
			json.Indent(&y, plain, "", "  ")
			marked, plain = x.Bytes(), y.Bytes()
		}
		lens := a.Engine().CacheLens()
		lens.ObserveRequest("fixture", marked)
		lens.ObserveRequest("fixture", plain)
		previous, _ := os.ReadFile(filepath.Join(a.Config().DataDir, "cache-prior.txt"))
		current, _ := os.ReadFile(filepath.Join(a.Config().DataDir, "cache-current.txt"))
		if !bytes.Equal(previous, current) {
			t.Fatalf("marker changed serialized artifact:\n%s\n%s", previous, current)
		}
		if !bytes.Contains(current, []byte(`\u003cescaped\u003e`)) {
			t.Fatal("escaped request bytes rewritten")
		}
	}
}
