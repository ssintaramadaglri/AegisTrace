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
		if !allowedKinds[s.Kind] {
			add(path+".kind", "must be model, tool, approval, checkpoint, or emit")
		}
		if (s.Kind == "model" || s.Kind == "tool") && s.Provider == "" {
			add(path+".provider", "is required")
		}
		if s.Kind == "tool" {
			if s.Action == "" {
				add(path+".action", "is required")
			}
			if !permissions[s.Action] && !permissions["*"] {
				add(path+".action", "tool is not permitted")
			}
		}
	}
	if _, ok := steps[w.Entry]; !ok {
		add("entry", "does not reference a step")
	}
	for i, s := range w.Steps {
		for field, target := range map[string]string{"next": s.Next, "on_error": s.OnError} {
			if target != "" {
				r.Edges++
				if _, ok := steps[target]; !ok {
					add(fmt.Sprintf("steps[%d].%s", i, field), "unknown target "+target)
				}
			}
		}
	}
	if len(r.Issues) > 0 {
		return r
	}
	verifyGraph(w.Entry, steps, &r)
	return r
}

func verifyGraph(entry string, steps map[string]workflow.Step, r *Report) {
	state := make(map[string]uint8)
	reachable := make(map[string]bool)
	var visit func(string)
	visit = func(id string) {
		if state[id] == 1 {
			r.Issues = append(r.Issues, Issue{"steps", "cycle detected at " + id})
			return
		}
		if state[id] == 2 {
			return
		}
		state[id], reachable[id] = 1, true
		s := steps[id]
		if s.Next != "" {
			visit(s.Next)
		}
		if s.OnError != "" {
			visit(s.OnError)
		}
		state[id] = 2
	}
	visit(entry)
	for id := range steps {
		if !reachable[id] {
			r.Issues = append(r.Issues, Issue{"steps." + id, "is unreachable"})
		}
	}
}
