package runtime_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/aegistrace/aegistrace/checkpoint"
	"github.com/aegistrace/aegistrace/journal"
	"github.com/aegistrace/aegistrace/provider"
	aegisruntime "github.com/aegistrace/aegistrace/runtime"
	"github.com/aegistrace/aegistrace/workflow"
)

func TestRunCheckpointAndReplay(t *testing.T) {
	w := workflow.Workflow{Name: "integration", Version: "1", Entry: "model", Budget: workflow.Budget{MaxSteps: 5, MaxToolCalls: 1, MaxTokens: 20, Timeout: time.Second}, Permissions: []workflow.Permission{{Tool: "task.create"}}, Steps: []workflow.Step{
