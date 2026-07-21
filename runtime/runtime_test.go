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
		calls++
		return provider.Response{Output: request.Input, UsageTokens: 2}, nil
	})); err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "run.atlog")
	writer, err := journal.Create(path, "integration-run")
	if err != nil {
		t.Fatal(err)
	}
	engine := aegisruntime.Engine{Artifact: artifact, Providers: registry, Journal: writer, Checkpoints: checkpoint.FileStore{Directory: filepath.Join(dir, "checkpoints")}}
	result, err := engine.Run(context.Background(), "integration-run", map[string]any{"message": "hello"})
	if err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	if calls != 2 || result.Steps != 3 || result.ToolCalls != 1 || result.Tokens != 4 {
		t.Fatalf("unexpected result: calls=%d result=%+v", calls, result)
	}
	state, err := (checkpoint.FileStore{Directory: filepath.Join(dir, "checkpoints")}).Load("integration-run")
	if err != nil {
		t.Fatal(err)
	}
	if state.NextStep != "tool" {
		t.Fatalf("unexpected checkpoint: %+v", state)
	}
	events, err := journal.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	replayed, err := aegisruntime.Replay(events)
	if err != nil {
		t.Fatal(err)
	}
	if replayed.Steps != 3 || replayed.RunID != "integration-run" || replayed.ToolCalls != 1 || replayed.Tokens != 4 {
		t.Fatalf("unexpected replay: %+v", replayed)
	}
	if calls != 2 {
		t.Fatal("replay invoked provider")
	}
}
