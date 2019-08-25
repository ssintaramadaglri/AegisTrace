// Package checkpoint persists resumable workflow state atomically.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var safeID = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type State struct {
	Version        string         `json:"version"`
	RunID          string         `json:"run_id"`
	WorkflowDigest string         `json:"workflow_digest"`
	NextStep       string         `json:"next_step"`
	Outputs        map[string]any `json:"outputs"`
	Steps          int            `json:"steps"`
	ToolCalls      int            `json:"tool_calls"`
	Tokens         int            `json:"tokens"`
	SavedAt        time.Time      `json:"saved_at"`
}

type Store interface {
	Save(State) error
	Load(string) (State, error)
}
type FileStore struct{ Directory string }

