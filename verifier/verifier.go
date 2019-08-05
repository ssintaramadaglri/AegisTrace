// Package verifier statically admits safe, finite AegisTrace workflows.
package verifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/aegistrace/aegistrace/workflow"
)

type Issue struct{ Path, Message string }

func (i Issue) Error() string { return i.Path + ": " + i.Message }

type Report struct {
	Workflow string  `json:"workflow"`
	Steps    int     `json:"steps"`
	Edges    int     `json:"edges"`
	Issues   []Issue `json:"-"`
}

func (r Report) OK() bool { return len(r.Issues) == 0 }
func (r Report) Error() error {
	if r.OK() {
		return nil
	}
	parts := make([]string, len(r.Issues))
	for i := range r.Issues {
		parts[i] = r.Issues[i].Error()
	}
	return fmt.Errorf("workflow rejected with %d issue(s): %s", len(parts), strings.Join(parts, "; "))
}

func Verify(w workflow.Workflow) Report {
	r := Report{Workflow: w.Name, Steps: len(w.Steps)}
	add := func(path, message string) { r.Issues = append(r.Issues, Issue{path, message}) }
	if strings.TrimSpace(w.Name) == "" {
		add("name", "is required")
	}
	if strings.TrimSpace(w.Version) == "" {
		add("version", "is required")
	}
	if w.Budget.MaxSteps <= 0 {
		add("budget.max_steps", "must be positive")
	}
	if w.Budget.MaxToolCalls < 0 {
		add("budget.max_tool_calls", "cannot be negative")
	}
	if w.Budget.MaxTokens < 0 {
		add("budget.max_tokens", "cannot be negative")
	}
	if w.Budget.Timeout <= 0 || w.Budget.Timeout > 24*time.Hour {
		add("budget.timeout", "must be between 1ns and 24h")
	}
	if len(w.Steps) == 0 {
		add("steps", "at least one step is required")
		return r
	}

	steps := make(map[string]workflow.Step, len(w.Steps))
	permissions := make(map[string]bool)
	for _, p := range w.Permissions {
		if p.Tool == "" {
			add("permissions", "tool cannot be empty")
		}
		permissions[p.Tool] = true
	}
	allowedKinds := map[string]bool{"model": true, "tool": true, "approval": true, "checkpoint": true, "emit": true}
	for i, s := range w.Steps {
		path := fmt.Sprintf("steps[%d]", i)
		if s.ID == "" {
			add(path+".id", "is required")
			continue
		}
		if _, exists := steps[s.ID]; exists {
			add(path+".id", "duplicate id "+s.ID)
		}
		steps[s.ID] = s
