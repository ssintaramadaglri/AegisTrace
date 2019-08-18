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
