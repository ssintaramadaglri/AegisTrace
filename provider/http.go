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
	var doc map[string]any
	if err := json.Unmarshal(raw, &doc); err != nil {
		return Response{}, fmt.Errorf("decode %s response: %w", p.Name, err)
	}
	var text string
	switch p.Style {
	case Anthropic:
		text = nestedText(doc, "content", "text")
	case Gemini:
		text = geminiText(doc)
	default:
		text = openAIText(doc)
	}
	if text == "" {
		return Response{}, fmt.Errorf("%s response contained no text", p.Name)
	}
	return Response{Output: text, UsageTokens: usage(doc), Metadata: map[string]string{"provider": p.Name, "model": p.Model}}, nil
}

func promptFrom(input map[string]any) string {
	if v, ok := input["prompt"].(string); ok {
		return v
	}
	raw, _ := json.Marshal(input)
	return string(raw)
}
func openAIText(d map[string]any) string {
	choices, _ := d["choices"].([]any)
	if len(choices) == 0 {
		return ""
	}
	c, _ := choices[0].(map[string]any)
	m, _ := c["message"].(map[string]any)
	s, _ := m["content"].(string)
	return s
}
func nestedText(d map[string]any, list, field string) string {
	xs, _ := d[list].([]any)
	if len(xs) == 0 {
		return ""
	}
	x, _ := xs[0].(map[string]any)
	s, _ := x[field].(string)
	return s
}
func geminiText(d map[string]any) string {
	cs, _ := d["candidates"].([]any)
	if len(cs) == 0 {
		return ""
	}
	c, _ := cs[0].(map[string]any)
	content, _ := c["content"].(map[string]any)
	parts, _ := content["parts"].([]any)
	if len(parts) == 0 {
		return ""
	}
	p, _ := parts[0].(map[string]any)
	s, _ := p["text"].(string)
	return s
}
func usage(d map[string]any) int {
	for _, key := range []string{"usage", "usageMetadata"} {
		if u, ok := d[key].(map[string]any); ok {
			for _, k := range []string{"total_tokens", "totalTokenCount", "output_tokens"} {
				if n, ok := u[k].(float64); ok {
					return int(n)
				}
			}
		}
	}
	return 0
}
