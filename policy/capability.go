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
	}
	if ttl <= 0 || ttl > 24*time.Hour {
		return "", fmt.Errorf("ttl must be between 1ns and 24h")
	}
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	claims := Claims{issuer, subject, append([]string(nil), scopes...), now.UTC(), now.Add(ttl).UTC(), hex.EncodeToString(nonce)}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	signature := ed25519.Sign(private, append([]byte("aegis-cap/v1."), payload...))
	return base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}

