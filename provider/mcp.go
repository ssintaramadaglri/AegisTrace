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
