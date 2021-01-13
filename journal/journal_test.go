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
