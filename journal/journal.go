// Package journal provides tamper-evident JSONL execution records.
package journal

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Event struct {
	Sequence uint64         `json:"sequence"`
	Time     time.Time      `json:"time"`
	RunID    string         `json:"run_id"`
	StepID   string         `json:"step_id,omitempty"`
	Type     string         `json:"type"`
	Data     map[string]any `json:"data,omitempty"`
	Previous string         `json:"previous,omitempty"`
	Hash     string         `json:"hash"`
}

type Writer struct {
	mu       sync.Mutex
	file     *os.File
	runID    string
	next     uint64
	previous string
}

func Create(path, runID string) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return nil, fmt.Errorf("create journal: %w", err)
	}
	return &Writer{file: f, runID: runID}, nil
}

func (w *Writer) Append(kind, stepID string, data map[string]any) (Event, error) {
