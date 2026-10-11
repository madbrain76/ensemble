package gateway

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestFormatProvenanceShape(t *testing.T) {
	got := FormatProvenance(Source{
		Type:        "webhook",
		ChannelID:   "webhook-game",
		UserID:      "server",
		DisplayName: "Game Server",
		Trust:       TrustTrusted,
	})
	want := "[GATEWAY source=webhook-game user=server trust=trusted]"
	if got != want {
		t.Fatalf("provenance header\n got: %s\nwant: %s", got, want)
	}
}

// TestWithProvenancePrependsAboveContent pins the thing the agent actually
// relies on: the header is above the body, separated from it, and the body
// is unmodified.
func TestWithProvenancePrependsAboveContent(t *testing.T) {
	env := Envelope{
		Content: "The goblin attacks.",
		Source: Source{
			ChannelID: "webhook-game",
			UserID:    "server",
			Trust:     TrustTrusted,
		},
	}
	got := env.WithProvenance()

	if !strings.HasPrefix(got, "[GATEWAY ") {
		t.Errorf("provenance must come first, got %q", got)
	}
	if !strings.HasSuffix(got, "The goblin attacks.") {
		t.Errorf("content must be preserved verbatim, got %q", got)
	}
	header, body, found := strings.Cut(got, "\n\n")
	if !found {
		t.Fatalf("header and body must be separated by a blank line, got %q", got)
	}
	if strings.Contains(header, "goblin") {
		t.Errorf("content leaked into the header: %q", header)
	}
	if body != "The goblin attacks." {
		t.Errorf("body = %q, want the original content", body)
	}
}

func TestTrustValid(t *testing.T) {
	for _, ok := range []Trust{TrustOwner, TrustTrusted, TrustExternal} {
		if !ok.Valid() {
			t.Errorf("%q should be valid", ok)
		}
	}
	// The case that matters: an unrecognized level must not pass as valid,
	// because the safe reading of "I do not know this trust level" is to
	// refuse, not to assume the lowest one and carry on.
	for _, bad := range []Trust{"", "admin", "OWNER", "root"} {
		if bad.Valid() {
			t.Errorf("%q should not be valid", bad)
		}
	}
}

// TestEnvelopeJSONRoundTrip guards the wire contract. The envelope crosses a
// process boundary, so a renamed JSON tag is a protocol break, not a
// refactor.
func TestEnvelopeJSONRoundTrip(t *testing.T) {
	in := Envelope{
		MessageID:      NewMessageID(),
		Seq:            42,
		AgentID:        "ensemble",
		Timestamp:      time.Now().UTC().Truncate(time.Second),
		ConversationID: "interview-1",
		Content:        "hello",
		ResponseMode:   ResponseAsync,
		Source: Source{
			Type:        "agent",
			ChannelID:   "interviewer",
			UserID:      "interviewer",
			DisplayName: "Interviewer",
			Trust:       TrustTrusted,
		},
	}
	data, err := json.Marshal(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out Envelope
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out != in {
		t.Fatalf("round trip changed the envelope:\n got %+v\nwant %+v", out, in)
	}

	for _, field := range []string{
		"message_id", "seq", "agent_id", "source", "timestamp",
		"conversation_id", "content", "response_mode", "trust",
	} {
		if !strings.Contains(string(data), `"`+field+`"`) {
			t.Errorf("wire form is missing field %q: %s", field, data)
		}
	}
}

// TestNewMessageIDIsUniqueAndOrdered checks both properties the queue depends
// on: IDs do not collide, and they sort chronologically.
func TestNewMessageIDIsUniqueAndOrdered(t *testing.T) {
	const n = 1000
	seen := make(map[string]bool, n)
	ids := make([]string, 0, n)
	for i := 0; i < n; i++ {
		id := NewMessageID()
		if seen[id] {
			t.Fatalf("duplicate message ID %q after %d draws", id, i)
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if len(ids[0]) != 36 {
		t.Errorf("message ID %q is %d chars, want 36", ids[0], len(ids[0]))
	}
	if ids[0][14] != '7' {
		t.Errorf("message ID %q is not version 7", ids[0])
	}

	// Generated in order, so a lexical sort must be a no-op across the
	// millisecond boundaries the batch spans.
	batch := []string{}
	for i := 0; i < 5; i++ {
		batch = append(batch, NewMessageID())
		time.Sleep(2 * time.Millisecond)
	}
	sorted := append([]string(nil), batch...)
	sort.Strings(sorted)
	for i := range batch {
		if batch[i] != sorted[i] {
			t.Fatalf("UUIDv7 IDs are not time ordered:\n got %v\nsorted %v", batch, sorted)
		}
	}
}
