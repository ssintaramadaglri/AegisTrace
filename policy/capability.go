// Package policy implements short-lived Ed25519 capability tokens.
package policy

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Claims struct {
	Issuer    string    `json:"iss"`
	Subject   string    `json:"sub"`
	Scopes    []string  `json:"scopes"`
	IssuedAt  time.Time `json:"iat"`
	ExpiresAt time.Time `json:"exp"`
	Nonce     string    `json:"nonce"`
}

func Issue(private ed25519.PrivateKey, issuer, subject string, scopes []string, ttl time.Duration, now time.Time) (string, error) {
	if issuer == "" || subject == "" {
		return "", fmt.Errorf("issuer and subject are required")
