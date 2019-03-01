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
