// Package workflow defines AegisTrace's portable workflow and artifact formats.
package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ArtifactMagic   = "AEGISTRACE"
	ArtifactVersion = "aegis/v1"
)

type Budget struct {
	MaxSteps     int           `json:"max_steps" yaml:"max_steps"`
	MaxToolCalls int           `json:"max_tool_calls" yaml:"max_tool_calls"`
	MaxTokens    int           `json:"max_tokens" yaml:"max_tokens"`
	Timeout      time.Duration `json:"timeout" yaml:"timeout"`
}

func (b *Budget) UnmarshalYAML(unmarshal func(any) error) error {
	var raw struct {
		MaxSteps     int    `yaml:"max_steps"`
