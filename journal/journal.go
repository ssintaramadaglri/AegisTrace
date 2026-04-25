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
	w.mu.Lock()
	defer w.mu.Unlock()
	e := Event{Sequence: w.next, Time: time.Now().UTC(), RunID: w.runID, StepID: stepID, Type: kind, Data: data, Previous: w.previous}
	hash, err := hashEvent(e)
	if err != nil {
		return e, err
	}
	e.Hash = hash
	encoded, err := json.Marshal(e)
	if err != nil {
		return e, fmt.Errorf("encode event: %w", err)
	}
	if _, err = w.file.Write(append(encoded, '\n')); err != nil {
		return e, fmt.Errorf("append journal: %w", err)
	}
	if err = w.file.Sync(); err != nil {
		return e, fmt.Errorf("sync journal: %w", err)
	}
	w.next++
	w.previous = hash
	return e, nil
}

func (w *Writer) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	err := w.file.Close()
	w.file = nil
	return err
}

func Load(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open journal: %w", err)
	}
	defer f.Close()
	var events []Event
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 64*1024), 4*1024*1024)
	for line := 1; s.Scan(); line++ {
		var e Event
		if err := json.Unmarshal(s.Bytes(), &e); err != nil {
			return nil, fmt.Errorf("journal line %d: %w", line, err)
		}
		events = append(events, e)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("read journal: %w", err)
	}
	return events, nil
}

func Verify(events []Event) error {
	previous := ""
	for i, e := range events {
		if e.Sequence != uint64(i) {
			return fmt.Errorf("event %d has sequence %d", i, e.Sequence)
		}
		if e.Previous != previous {
			return fmt.Errorf("event %d previous hash mismatch", i)
		}
		have := e.Hash
		e.Hash = ""
		want, err := hashEvent(e)
		if err != nil {
			return err
		}
		if have != want {
			return fmt.Errorf("event %d hash mismatch", i)
		}
		previous = have
	}
	return nil
}

func hashEvent(e Event) (string, error) {
	e.Hash = ""
	data, err := json.Marshal(struct {
		Sequence uint64         `json:"sequence"`
		Time     time.Time      `json:"time"`
		RunID    string         `json:"run_id"`
		StepID   string         `json:"step_id,omitempty"`
		Type     string         `json:"type"`
		Data     map[string]any `json:"data,omitempty"`
		Previous string         `json:"previous,omitempty"`
	}{e.Sequence, e.Time, e.RunID, e.StepID, e.Type, e.Data, e.Previous})
	if err != nil {
		return "", fmt.Errorf("canonicalize event: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
