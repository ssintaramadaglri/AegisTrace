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
