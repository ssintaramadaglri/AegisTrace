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

func (s FileStore) Save(state State) error {
	if !safeID.MatchString(state.RunID) {
		return fmt.Errorf("invalid run id")
	}
	if err := os.MkdirAll(s.Directory, 0o700); err != nil {
		return err
	}
	state.Version, state.SavedAt = "aegis-checkpoint/v1", time.Now().UTC()
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(s.Directory, state.RunID+".json")
	temp := path + ".tmp"
	if err := os.WriteFile(temp, append(data, '\n'), 0o600); err != nil {
		return err
	}
	if err := os.Rename(temp, path); err != nil {
		_ = os.Remove(temp)
		return fmt.Errorf("commit checkpoint: %w", err)
	}
	return nil
}

func (s FileStore) Load(runID string) (State, error) {
	var state State
	if !safeID.MatchString(runID) {
		return state, fmt.Errorf("invalid run id")
	}
	data, err := os.ReadFile(filepath.Join(s.Directory, runID+".json"))
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return state, err
	}
	if state.Version != "aegis-checkpoint/v1" || state.RunID != runID {
		return state, fmt.Errorf("invalid checkpoint")
	}
	return state, nil
}
