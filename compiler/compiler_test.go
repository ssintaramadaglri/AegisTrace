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
	loaded, err := compiler.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Digest != artifact.Digest {
		t.Fatal("digest changed after load")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data = []byte(strings.Replace(string(data), `"name": "test"`, `"name": "changed"`, 1))
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := compiler.Load(path); err == nil || !strings.Contains(err.Error(), "digest mismatch") {
		t.Fatalf("expected digest failure, got %v", err)
	}
}

func TestCompileRejectsUnknownField(t *testing.T) {
	_, _, err := compiler.Compile([]byte("name: x\nunknown: true\n"), time.Now())
	if err == nil || !strings.Contains(err.Error(), "field unknown") {
		t.Fatalf("expected strict parse error, got %v", err)
	}
}
