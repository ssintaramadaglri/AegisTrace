package proof_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"path/filepath"
	"testing"

	"github.com/aegistrace/aegistrace/journal"
	"github.com/aegistrace/aegistrace/proof"
)

func TestAttestation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.atlog")
	writer, err := journal.Create(path, "proof-run")
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		if _, err := writer.Append("event", "", map[string]any{"number": i}); err != nil {
			t.Fatal(err)
		}
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	events, err := journal.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	_, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	attestation, err := proof.Create(events, 3, private)
	if err != nil {
		t.Fatal(err)
	}
	if err := proof.Verify(attestation); err != nil {
		t.Fatal(err)
	}
	attestation.EventHash = events[2].Hash
	if err := proof.Verify(attestation); err == nil {
		t.Fatal("tampered proof verified")
	}
}

func TestKeyFileRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "key.json")
	pub, err := proof.GenerateKey(path)
	if err != nil {
		t.Fatal(err)
	}
	private, err := proof.LoadPrivate(path)
	if err != nil {
		t.Fatal(err)
	}
	if !private.Public().(ed25519.PublicKey).Equal(pub) {
		t.Fatal("public key mismatch")
	}
}
