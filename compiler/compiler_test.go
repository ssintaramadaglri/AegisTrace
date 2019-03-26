package compiler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aegistrace/aegistrace/compiler"
)

const validWorkflow = `
name: test
version: "1"
budget:
  max_steps: 4
  max_tool_calls: 1
  max_tokens: 50
  timeout: 2s
permissions:
  - tool: task.create
steps:
  - id: first
    kind: emit
    input:
      value: ok
  - id: second
    kind: tool
    provider: echo
    action: task.create
`

func TestCompileNormalizesAndSeals(t *testing.T) {
	now := time.Date(2026, 8, 31, 1, 2, 3, 0, time.UTC)
	artifact, report, err := compiler.Compile([]byte(validWorkflow), now)
	if err != nil {
		t.Fatal(err)
	}
	if report.Steps != 2 || artifact.Workflow.Entry != "first" {
		t.Fatalf("unexpected report or entry: %#v %#v", report, artifact.Workflow)
	}
	if artifact.Workflow.Steps[0].Next != "second" {
		t.Fatal("sequential edge was not normalized")
	}
	if err := artifact.ValidateSeal(); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(t.TempDir(), "test.aegis")
	if err := compiler.Save(path, artifact); err != nil {
		t.Fatal(err)
	}
