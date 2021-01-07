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
