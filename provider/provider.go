// Package provider isolates model and tool effects behind a narrow interface.
package provider

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type Request struct {
	RunID     string         `json:"run_id"`
	StepID    string         `json:"step_id"`
	Action    string         `json:"action,omitempty"`
	Input     map[string]any `json:"input,omitempty"`
	MaxTokens int            `json:"max_tokens,omitempty"`
}

type Response struct {
	Output      any               `json:"output"`
	UsageTokens int               `json:"usage_tokens,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type Provider interface {
	Execute(context.Context, Request) (Response, error)
}
type Func func(context.Context, Request) (Response, error)

func (f Func) Execute(ctx context.Context, r Request) (Response, error) { return f(ctx, r) }

type Registry struct {
