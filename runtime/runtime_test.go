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
		{ID: "model", Kind: "model", Provider: "test", Input: map[string]any{"prompt": "${input.message}"}, Next: "save"},
		{ID: "save", Kind: "checkpoint", Next: "tool"},
		{ID: "tool", Kind: "tool", Provider: "test", Action: "task.create", Input: map[string]any{"body": "${steps.model}"}},
	}}
	artifact, err := workflow.NewArtifact(w, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	calls := 0
	registry := provider.NewRegistry()
	if err := registry.Register("test", provider.Func(func(_ context.Context, request provider.Request) (provider.Response, error) {
