// Package compiler parses, normalizes, verifies, and seals workflow artifacts.
package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aegistrace/aegistrace/verifier"
	"github.com/aegistrace/aegistrace/workflow"
	"gopkg.in/yaml.v3"
)

func Parse(source []byte) (workflow.Workflow, error) {
	var w workflow.Workflow
	dec := yaml.NewDecoder(bytes.NewReader(source))
	dec.KnownFields(true)
	if err := dec.Decode(&w); err != nil {
		return w, fmt.Errorf("parse workflow: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); err != io.EOF {
		return w, fmt.Errorf("workflow must contain one YAML document")
	}
	normalize(&w)
	return w, nil
}

func normalize(w *workflow.Workflow) {
	if w.Version == "" {
		w.Version = "1"
	}
	if w.Entry == "" && len(w.Steps) > 0 {
		w.Entry = w.Steps[0].ID
	}
	if w.Budget.MaxSteps == 0 {
		w.Budget.MaxSteps = max(32, len(w.Steps)*2)
	}
	if w.Budget.Timeout == 0 {
		w.Budget.Timeout = time.Minute
	}
	for i := 0; i+1 < len(w.Steps); i++ {
		if w.Steps[i].Next == "" {
			w.Steps[i].Next = w.Steps[i+1].ID
		}
	}
}

func Compile(source []byte, now time.Time) (workflow.Artifact, verifier.Report, error) {
	w, err := Parse(source)
	if err != nil {
		return workflow.Artifact{}, verifier.Report{}, err
	}
	report := verifier.Verify(w)
	if err := report.Error(); err != nil {
		return workflow.Artifact{}, report, err
	}
	artifact, err := workflow.NewArtifact(w, now)
	return artifact, report, err
}

func Save(path string, artifact workflow.Artifact) error {
	data, err := json.MarshalIndent(artifact, "", "  ")
	if err != nil {
		return fmt.Errorf("encode artifact: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write artifact: %w", err)
	}
	return nil
}

func Load(path string) (workflow.Artifact, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return workflow.Artifact{}, fmt.Errorf("read artifact: %w", err)
	}
	var artifact workflow.Artifact
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&artifact); err != nil {
		return artifact, fmt.Errorf("decode artifact: %w", err)
	}
	if err := artifact.ValidateSeal(); err != nil {
		return artifact, err
	}
	if err := verifier.Verify(artifact.Workflow).Error(); err != nil {
		return artifact, err
	}
	return artifact, nil
}
