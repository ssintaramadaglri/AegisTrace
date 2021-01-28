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
