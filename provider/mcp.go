package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type MCP struct {
	Endpoint, Token string
	Client          *http.Client
}

type rpcRequest struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}
type rpcResponse struct {
	Result any `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func (m *MCP) Execute(ctx context.Context, request Request) (Response, error) {
	if m.Endpoint == "" {
		return Response{}, fmt.Errorf("MCP endpoint is not configured")
	}
	payload := rpcRequest{"2.0", 1, "tools/call", map[string]any{"name": request.Action, "arguments": request.Input}}
	body, err := json.Marshal(payload)
	if err != nil {
		return Response{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, m.Endpoint, bytes.NewReader(body))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if m.Token != "" {
		req.Header.Set("Authorization", "Bearer "+m.Token)
	}
	client := m.Client
