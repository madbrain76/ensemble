package gateway

import (
	"bytes"
	"path/filepath"
	"testing"
)

func TestPublicKeyEncodeRoundTrip(t *testing.T) {
	key, err := GenerateKeypair()
	if err != nil {
		t.Fatalf("GenerateKeypair: %v", err)
	}
	encoded := EncodePublicKey(key.Public)
	decoded, err := ParsePublicKey(encoded)
	if err != nil {
		t.Fatalf("ParsePublicKey(%q): %v", encoded, err)
	}
	if !bytes.Equal(decoded, key.Public) {
		t.Fatalf("round trip changed the key")
	}
}

func TestParsePublicKeyRejectsWrongLength(t *testing.T) {
	// A positive control: a well-formed key of the right length parses.
	good := EncodePublicKey(make([]byte, PublicKeySize))
	if _, err := ParsePublicKey(good); err != nil {
		t.Fatalf("positive control failed: %v", err)
	}
	// The actual assertion: a short key is rejected rather than accepted and
	// then failing later as an inscrutable handshake error.
	short := EncodePublicKey(make([]byte, 16))
	if _, err := ParsePublicKey(short); err == nil {
		t.Fatalf("expected a 16-byte key to be rejected")
	}
}

// TestLoadOrCreateKeypairStableAcrossRestart is the one that matters: the
// gateway config names an agent by public key, so if reloading a key file
// produced a different public half, the agent would be locked out of its own
// config with no error until the handshake failed.
func TestLoadOrCreateKeypairStableAcrossRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sub", "agent.key")

	first, err := LoadOrCreateKeypair(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	second, err := LoadOrCreateKeypair(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}

	if !bytes.Equal(first.Private, second.Private) {
		t.Errorf("private key changed across reload")
	}
	if !bytes.Equal(first.Public, second.Public) {
		t.Errorf("public key changed across reload: generated %x, derived %x",
			first.Public, second.Public)
	}
}

// TestLoadedKeypairActuallyHandshakes proves the derived public half is not
// merely stable but cryptographically correct: a peer told only the public
// key can complete a KK handshake with the reloaded private key.
func TestLoadedKeypairActuallyHandshakes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "agent.key")
	if _, err := LoadOrCreateKeypair(path); err != nil {
		t.Fatalf("create: %v", err)
	}
	agent, err := LoadOrCreateKeypair(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	peer, err := GenerateKeypair()
	if err != nil {
		t.Fatalf("peer key: %v", err)
	}

	initiator, err := NewNoiseKK("initiator", agent, peer.Public, []byte("test"))
	if err != nil {
		t.Fatalf("initiator: %v", err)
	}
	responder, err := NewNoiseKK("responder", peer, agent.Public, []byte("test"))
	if err != nil {
		t.Fatalf("responder: %v", err)
	}

	msg1, err := initiator.WriteHandshakeMessage(nil)
	if err != nil {
		t.Fatalf("msg1: %v", err)
	}
	if _, err := responder.ReadHandshakeMessage(msg1); err != nil {
		t.Fatalf("responder read msg1: %v", err)
	}
	msg2, err := responder.WriteHandshakeMessage(nil)
	if err != nil {
		t.Fatalf("msg2: %v", err)
	}
	if _, err := initiator.ReadHandshakeMessage(msg2); err != nil {
		t.Fatalf("initiator read msg2: %v", err)
	}
	if !initiator.IsHandshakeComplete() || !responder.IsHandshakeComplete() {
		t.Fatalf("handshake incomplete after two messages")
	}
}
