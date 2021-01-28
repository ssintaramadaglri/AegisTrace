package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type APIStyle string

const (
	OpenAI    APIStyle = "openai"
	Anthropic APIStyle = "anthropic"
	Gemini    APIStyle = "gemini"
)

type HTTPModel struct {
	Name     string
	Style    APIStyle
	Endpoint string
	APIKey   string
	Model    string
	Client   *http.Client
}

func (p *HTTPModel) Execute(ctx context.Context, request Request) (Response, error) {
	if p.Endpoint == "" || p.Model == "" {
		return Response{}, fmt.Errorf("provider %s is not configured", p.Name)
	}
	body := p.requestBody(request)
	encoded, err := json.Marshal(body)
	if err != nil {
		return Response{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.Endpoint, bytes.NewReader(encoded))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	switch p.Style {
	case Anthropic:
		req.Header.Set("x-api-key", p.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	case Gemini:
		if p.APIKey != "" {
			req.Header.Set("x-goog-api-key", p.APIKey)
		}
	default:
		if p.APIKey != "" {
			req.Header.Set("Authorization", "Bearer "+p.APIKey)
		}
	}
	client := p.Client
	if client == nil {
		client = &http.Client{Timeout: 90 * time.Second}
	}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, fmt.Errorf("%s request: %w", p.Name, err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Response{}, fmt.Errorf("%s returned %s: %s", p.Name, resp.Status, strings.TrimSpace(string(raw)))
	}
	return p.parseResponse(raw)
}

func (p *HTTPModel) requestBody(r Request) map[string]any {
	prompt := promptFrom(r.Input)
	switch p.Style {
	case Anthropic:
		return map[string]any{"model": p.Model, "max_tokens": max(1, r.MaxTokens), "messages": []any{map[string]any{"role": "user", "content": prompt}}}
	case Gemini:
		return map[string]any{"contents": []any{map[string]any{"role": "user", "parts": []any{map[string]any{"text": prompt}}}}, "generationConfig": map[string]any{"maxOutputTokens": max(1, r.MaxTokens)}}
	default:
		return map[string]any{"model": p.Model, "messages": []any{map[string]any{"role": "user", "content": prompt}}, "max_tokens": max(1, r.MaxTokens)}
	}
}

func (p *HTTPModel) parseResponse(raw []byte) (Response, error) {
