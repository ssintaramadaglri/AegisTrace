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
