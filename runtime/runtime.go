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
			_, _ = e.Journal.Append("step.failed", step.ID, map[string]any{"error": err.Error()})
			if step.OnError != "" {
				current = step.OnError
				continue
			}
			return result, e.fail("step "+step.ID, err)
		}
		if step.Kind == "tool" {
			result.ToolCalls++
		}
		result.Tokens += tokens
		if budget.MaxTokens > 0 && result.Tokens > budget.MaxTokens {
			return result, e.fail("token budget exhausted", nil)
		}
		result.Output, result.Outputs[step.ID] = output, output
		_, _ = e.Journal.Append("step.completed", step.ID, map[string]any{"output": output, "tokens": tokens})
		if step.Kind == "checkpoint" && e.Checkpoints != nil {
			state := checkpoint.State{RunID: runID, WorkflowDigest: e.Artifact.Digest, NextStep: step.Next, Outputs: result.Outputs, Steps: result.Steps, ToolCalls: result.ToolCalls, Tokens: result.Tokens}
			if err := e.Checkpoints.Save(state); err != nil {
				return result, e.fail("save checkpoint", err)
			}
			_, _ = e.Journal.Append("checkpoint.saved", step.ID, map[string]any{"next": step.Next})
		}
		current = step.Next
	}
	result.Duration = time.Since(started)
	_, err := e.Journal.Append("run.completed", "", map[string]any{"output": result.Output, "steps": result.Steps, "tools": result.ToolCalls, "tokens": result.Tokens})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (e *Engine) execute(ctx context.Context, runID string, step workflow.Step, input map[string]any, result Result) (any, int, error) {
	if step.RequiresApproval || step.Kind == "approval" {
		if e.Approver == nil {
			return nil, 0, fmt.Errorf("approval required")
		}
		if err := e.Approver.Approve(ctx, step, input); err != nil {
			return nil, 0, err
		}
		if step.Kind == "approval" {
			return map[string]any{"approved": true}, 0, nil
		}
	}
	switch step.Kind {
	case "model", "tool":
		if step.Kind == "tool" {
			if result.ToolCalls >= e.Artifact.Workflow.Budget.MaxToolCalls {
				return nil, 0, fmt.Errorf("tool-call budget exhausted")
			}
		}
		p, err := e.Providers.Get(step.Provider)
		if err != nil {
			return nil, 0, err
		}
		response, err := p.Execute(ctx, provider.Request{RunID: runID, StepID: step.ID, Action: step.Action, Input: input, MaxTokens: e.Artifact.Workflow.Budget.MaxTokens - result.Tokens})
		return response.Output, response.UsageTokens, err
	case "checkpoint":
		return map[string]any{"checkpoint": step.ID}, 0, nil
	case "emit":
		return input, 0, nil
	default:
		return nil, 0, fmt.Errorf("unsupported admitted step kind %q", step.Kind)
	}
}

func (e *Engine) fail(message string, cause error) error {
	text := message
	if cause != nil {
		text += ": " + cause.Error()
	}
	_, _ = e.Journal.Append("run.failed", "", map[string]any{"error": text})
	return fmt.Errorf("%s", text)
}
func (e *Engine) stepFailure(step workflow.Step, cause error) error {
	return e.fail("resolve step "+step.ID, cause)
}

func Replay(events []journal.Event) (Result, error) {
	var result Result
	if err := journal.Verify(events); err != nil {
		return result, err
	}
	result.Outputs = make(map[string]any)
	for _, event := range events {
		if result.RunID == "" {
			result.RunID = event.RunID
		}
		switch event.Type {
		case "step.completed":
			result.Steps++
			result.Output = event.Data["output"]
			result.Outputs[event.StepID] = result.Output
		case "run.completed":
			result.Output = event.Data["output"]
			result.ToolCalls = eventInt(event.Data["tools"])
			result.Tokens = eventInt(event.Data["tokens"])
		case "run.failed":
			return result, fmt.Errorf("recorded run failed: %v", event.Data["error"])
		}
	}
	return result, nil
}

func eventInt(value any) int {
	switch number := value.(type) {
	case int:
		return number
	case int64:
		return int(number)
	case float64:
		return int(number)
	default:
		return 0
	}
}

func resolveMap(values map[string]any, input, outputs map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(values))
	for key, value := range values {
		resolved, err := resolve(value, input, outputs)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", key, err)
		}
		out[key] = resolved
	}
	return out, nil
}
func resolve(value any, input, outputs map[string]any) (any, error) {
	switch v := value.(type) {
	case string:
		if !strings.HasPrefix(v, "${") || !strings.HasSuffix(v, "}") {
			return v, nil
		}
		path := strings.TrimSuffix(strings.TrimPrefix(v, "${"), "}")
		if strings.HasPrefix(path, "input.") {
			value, ok := lookup(input, strings.TrimPrefix(path, "input."))
			if !ok {
				return nil, fmt.Errorf("unknown input %q", path)
			}
			return value, nil
		}
		if strings.HasPrefix(path, "steps.") {
			value, ok := outputs[strings.TrimPrefix(path, "steps.")]
			if !ok {
				return nil, fmt.Errorf("unknown step output %q", path)
			}
			return value, nil
		}
		return nil, fmt.Errorf("unsupported expression %q", path)
	case map[string]any:
		return resolveMap(v, input, outputs)
	case []any:
		out := make([]any, len(v))
		for i := range v {
			item, err := resolve(v[i], input, outputs)
			if err != nil {
				return nil, err
			}
			out[i] = item
		}
		return out, nil
	default:
		return value, nil
	}
}
func lookup(values map[string]any, path string) (any, bool) {
	var current any = values
	for _, part := range strings.Split(path, ".") {
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = object[part]
		if !ok {
			return nil, false
		}
	}
	return current, true
}
