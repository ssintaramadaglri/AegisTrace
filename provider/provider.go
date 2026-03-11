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
	mu        sync.RWMutex
	providers map[string]Provider
}

func NewRegistry() *Registry { return &Registry{providers: make(map[string]Provider)} }
func (r *Registry) Register(name string, p Provider) error {
	if name == "" || p == nil {
		return fmt.Errorf("provider name and implementation are required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = p
	return nil
}
func (r *Registry) Get(name string) (Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %q is not registered", name)
	}
	return p, nil
}
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Echo is deterministic and useful for local workflows and tests.
type Echo struct{}

func (Echo) Execute(_ context.Context, request Request) (Response, error) {
	return Response{Output: map[string]any{"action": request.Action, "input": request.Input}, Metadata: map[string]string{"provider": "echo"}}, nil
}
