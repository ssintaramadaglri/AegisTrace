package verifier_test

import (
	"strings"
	"testing"
	"time"

	"github.com/aegistrace/aegistrace/verifier"
	"github.com/aegistrace/aegistrace/workflow"
)

func base() workflow.Workflow {
	return workflow.Workflow{Name: "test", Version: "1", Entry: "a", Budget: workflow.Budget{MaxSteps: 4, MaxToolCalls: 1, Timeout: time.Second}, Steps: []workflow.Step{{ID: "a", Kind: "emit"}}}
}

func TestVerifyAcceptsFiniteWorkflow(t *testing.T) {
	report := verifier.Verify(base())
	if err := report.Error(); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyRejectsCycleAndUnreachableStep(t *testing.T) {
	w := base()
	w.Steps = []workflow.Step{{ID: "a", Kind: "emit", Next: "b"}, {ID: "b", Kind: "emit", Next: "a"}, {ID: "orphan", Kind: "emit"}}
	err := verifier.Verify(w).Error()
	if err == nil || !strings.Contains(err.Error(), "cycle detected") || !strings.Contains(err.Error(), "unreachable") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerifyRejectsUnpermittedTool(t *testing.T) {
	w := base()
	w.Steps[0] = workflow.Step{ID: "a", Kind: "tool", Provider: "echo", Action: "system.delete"}
	err := verifier.Verify(w).Error()
	if err == nil || !strings.Contains(err.Error(), "not permitted") {
		t.Fatalf("unexpected error: %v", err)
	}
}
