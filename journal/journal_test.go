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
