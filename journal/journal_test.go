package journal_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/aegistrace/aegistrace/journal"
)

func TestJournalRoundTripAndTamperDetection(t *testing.T) {
	path := filepath.Join(t.TempDir(), "run.atlog")
	writer, err := journal.Create(path, "run-1")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Append("run.started", "", map[string]any{"input": "safe"}); err != nil {
		t.Fatal(err)
	}
	if _, err := writer.Append("step.completed", "one", map[string]any{"output": "ok"}); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	events, err := journal.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 2 {
		t.Fatalf("got %d events", len(events))
	}
	if err := journal.Verify(events); err != nil {
		t.Fatal(err)
	}
	events[0].Data["input"] = "tampered"
	if err := journal.Verify(events); err == nil || !strings.Contains(err.Error(), "hash mismatch") {
		t.Fatalf("expected tamper failure, got %v", err)
	}
}
