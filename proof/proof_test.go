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
