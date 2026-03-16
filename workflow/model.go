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
	Input            map[string]any `json:"input,omitempty" yaml:"input,omitempty"`
	Next             string         `json:"next,omitempty" yaml:"next,omitempty"`
	OnError          string         `json:"on_error,omitempty" yaml:"on_error,omitempty"`
	RequiresApproval bool           `json:"requires_approval,omitempty" yaml:"requires_approval,omitempty"`
}

type Workflow struct {
	Name        string            `json:"name" yaml:"name"`
	Version     string            `json:"version" yaml:"version"`
	Description string            `json:"description,omitempty" yaml:"description,omitempty"`
	Entry       string            `json:"entry" yaml:"entry"`
	Budget      Budget            `json:"budget" yaml:"budget"`
	Permissions []Permission      `json:"permissions,omitempty" yaml:"permissions,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Steps       []Step            `json:"steps" yaml:"steps"`
}

type Artifact struct {
	Magic      string    `json:"magic"`
	Format     string    `json:"format"`
	Digest     string    `json:"digest"`
	CompiledAt time.Time `json:"compiled_at"`
	Workflow   Workflow  `json:"workflow"`
}

func Digest(w Workflow) (string, error) {
	data, err := json.Marshal(w)
	if err != nil {
		return "", fmt.Errorf("canonicalize workflow: %w", err)
	}
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func NewArtifact(w Workflow, now time.Time) (Artifact, error) {
	digest, err := Digest(w)
	if err != nil {
		return Artifact{}, err
	}
	return Artifact{Magic: ArtifactMagic, Format: ArtifactVersion, Digest: digest, CompiledAt: now.UTC(), Workflow: w}, nil
}

func (a Artifact) ValidateSeal() error {
	if a.Magic != ArtifactMagic || a.Format != ArtifactVersion {
		return fmt.Errorf("unsupported artifact %q/%q", a.Magic, a.Format)
	}
	digest, err := Digest(a.Workflow)
	if err != nil {
		return err
	}
	if digest != a.Digest {
		return fmt.Errorf("artifact digest mismatch: have %s, computed %s", a.Digest, digest)
	}
	return nil
}
