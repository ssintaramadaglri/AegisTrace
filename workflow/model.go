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
		MaxToolCalls int    `yaml:"max_tool_calls"`
		MaxTokens    int    `yaml:"max_tokens"`
		Timeout      string `yaml:"timeout"`
	}
	if err := unmarshal(&raw); err != nil {
		return err
	}
	b.MaxSteps, b.MaxToolCalls, b.MaxTokens = raw.MaxSteps, raw.MaxToolCalls, raw.MaxTokens
	if raw.Timeout == "" {
		b.Timeout = time.Minute
		return nil
	}
	d, err := time.ParseDuration(raw.Timeout)
	if err != nil {
		return fmt.Errorf("invalid timeout %q: %w", raw.Timeout, err)
	}
	b.Timeout = d
	return nil
}

type Permission struct {
	Tool   string   `json:"tool" yaml:"tool"`
	Scopes []string `json:"scopes,omitempty" yaml:"scopes,omitempty"`
}

type Step struct {
	ID               string         `json:"id" yaml:"id"`
	Kind             string         `json:"kind" yaml:"kind"`
	Provider         string         `json:"provider,omitempty" yaml:"provider,omitempty"`
	Action           string         `json:"action,omitempty" yaml:"action,omitempty"`
