package runtime

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aegistrace/aegistrace/checkpoint"
	"github.com/aegistrace/aegistrace/journal"
	"github.com/aegistrace/aegistrace/provider"
	"github.com/aegistrace/aegistrace/workflow"
)

type Approver interface {
	Approve(context.Context, workflow.Step, map[string]any) error
}
type ApproverFunc func(context.Context, workflow.Step, map[string]any) error

func (f ApproverFunc) Approve(ctx context.Context, step workflow.Step, input map[string]any) error {
	return f(ctx, step, input)
}

type Observer interface{ Event(string, map[string]any) }
type ObserverFunc func(string, map[string]any)

func (f ObserverFunc) Event(name string, fields map[string]any) { f(name, fields) }

type Engine struct {
	Artifact    workflow.Artifact
	Providers   *provider.Registry
	Journal     *journal.Writer
	Checkpoints checkpoint.Store
	Approver    Approver
	Observer    Observer
}

type Result struct {
	RunID     string         `json:"run_id"`
	Output    any            `json:"output"`
	Outputs   map[string]any `json:"outputs"`
	Steps     int            `json:"steps"`
	ToolCalls int            `json:"tool_calls"`
	Tokens    int            `json:"tokens"`
	Duration  time.Duration  `json:"duration"`
}

func (e *Engine) Run(ctx context.Context, runID string, input map[string]any) (Result, error) {
	result := Result{RunID: runID, Outputs: make(map[string]any)}
	if e.Providers == nil || e.Journal == nil {
		return result, fmt.Errorf("providers and journal are required")
	}
	budget := e.Artifact.Workflow.Budget
	ctx, cancel := context.WithTimeout(ctx, budget.Timeout)
	defer cancel()
	started := time.Now()
	_, _ = e.Journal.Append("run.started", "", map[string]any{"artifact": e.Artifact.Digest, "input": input})
	steps := make(map[string]workflow.Step, len(e.Artifact.Workflow.Steps))
	for _, step := range e.Artifact.Workflow.Steps {
		steps[step.ID] = step
	}
	current := e.Artifact.Workflow.Entry
	for current != "" {
		if err := ctx.Err(); err != nil {
			return result, e.fail("budget timeout", err)
		}
		if result.Steps >= budget.MaxSteps {
			return result, e.fail("step budget exhausted", nil)
		}
		step, ok := steps[current]
		if !ok {
			return result, e.fail("missing admitted step "+current, nil)
		}
		resolved, err := resolveMap(step.Input, input, result.Outputs)
		if err != nil {
			return result, e.stepFailure(step, err)
		}
		_, _ = e.Journal.Append("step.started", step.ID, map[string]any{"kind": step.Kind, "input": resolved})
		if e.Observer != nil {
			e.Observer.Event("step.started", map[string]any{"run_id": runID, "step_id": step.ID, "kind": step.Kind})
		}
		output, tokens, err := e.execute(ctx, runID, step, resolved, result)
		result.Steps++
		if err != nil {
