package gateway

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/flynn/noise"
	"golang.org/x/crypto/curve25519"
)

// KeyPrefix labels an encoded public key with its algorithm.
//
// The keys here are Curve25519 Diffie-Hellman keys (noise.DH25519), not
// ed25519 signing keys. Noise uses them for key agreement, never for
// signatures. Some deployments label these "ed25519:" out of habit; that
// label is wrong, and a wrong label on a key is the kind of small untruth
// that costs somebody an afternoon. ParsePublicKey therefore accepts any
// prefix, but everything this package writes is labeled honestly.
const KeyPrefix = "x25519:"

// PublicKeySize is the length in bytes of a Curve25519 public key.
const PublicKeySize = 32

// EncodePublicKey renders a raw public key as "x25519:<base64>".
func EncodePublicKey(pub []byte) string {
	return KeyPrefix + base64.StdEncoding.EncodeToString(pub)
}

// ParsePublicKey decodes a "<type>:<base64>" public key into raw bytes.
//
// The algorithm label is accepted but not trusted: the length check below is
// what actually rejects a key that could not possibly work. A 32-byte check
// here turns a confusing handshake failure into a clear configuration error.
func ParsePublicKey(s string) ([]byte, error) {
	encoded := s
	if _, rest, found := strings.Cut(s, ":"); found {
		encoded = rest
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode public key %q: %w", s, err)
	}
	if len(data) != PublicKeySize {
		return nil, fmt.Errorf("public key %q is %d bytes, want %d", s, len(data), PublicKeySize)
	}
	return data, nil
}

// LoadOrCreateKeypair reads a keypair from path, generating and saving a new
// one if the file does not exist.
//
// Every binary that speaks to the gateway needs a stable identity across
// restarts, and the gateway config refers to agents by public key. Generating
// a fresh key on every start would silently lock an agent out of its own
// config, so the key has to live on disk.
func LoadOrCreateKeypair(path string) (noise.DHKey, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		priv, decErr := base64.StdEncoding.DecodeString(strings.TrimSpace(string(data)))
		if decErr != nil {
			return noise.DHKey{}, fmt.Errorf("decode key file %s: %w", path, decErr)
		}
		if len(priv) != PublicKeySize {
			return noise.DHKey{}, fmt.Errorf("key file %s holds %d bytes, want %d", path, len(priv), PublicKeySize)
		}
		// Derive the public half from the stored private half so the pair is
		// always consistent with what is on disk.
		pub, pubErr := curve25519.X25519(priv, curve25519.Basepoint)
		if pubErr != nil {
			return noise.DHKey{}, fmt.Errorf("derive public key from %s: %w", path, pubErr)
		}
		return noise.DHKey{Private: priv, Public: pub}, nil
	}
	if !os.IsNotExist(err) {
		return noise.DHKey{}, fmt.Errorf("read key file %s: %w", path, err)
	}

	key, err := GenerateKeypair()
	if err != nil {
		return noise.DHKey{}, err
	}
	if mkErr := os.MkdirAll(filepath.Dir(path), 0o700); mkErr != nil {
		return noise.DHKey{}, mkErr
	}
	encoded := base64.StdEncoding.EncodeToString(key.Private)
	if wErr := os.WriteFile(path, []byte(encoded), 0o600); wErr != nil {
		return noise.DHKey{}, fmt.Errorf("write key file %s: %w", path, wErr)
	}
	return key, nil
}
