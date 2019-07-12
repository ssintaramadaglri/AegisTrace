// Package proof creates portable Merkle and Ed25519 journal attestations.
package proof

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aegistrace/aegistrace/journal"
)

type Sibling struct {
	Hash string `json:"hash"`
	Left bool   `json:"left"`
}
type Attestation struct {
	Version    string    `json:"version"`
	EventIndex int       `json:"event_index"`
	EventCount int       `json:"event_count"`
	EventHash  string    `json:"event_hash"`
	Root       string    `json:"merkle_root"`
	Path       []Sibling `json:"path"`
	PublicKey  string    `json:"public_key"`
	Signature  string    `json:"signature"`
}

type KeyFile struct {
	Version string `json:"version"`
	Public  string `json:"public"`
	Private string `json:"private"`
}

func GenerateKey(path string) (ed25519.PublicKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate key: %w", err)
	}
	key := KeyFile{"aegis-key/v1", base64.StdEncoding.EncodeToString(pub), base64.StdEncoding.EncodeToString(priv)}
	data, _ := json.MarshalIndent(key, "", "  ")
	if err := os.WriteFile(path, append(data, '\n'), 0o600); err != nil {
		return nil, fmt.Errorf("write key: %w", err)
	}
	return pub, nil
}

func LoadPrivate(path string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read key: %w", err)
	}
	var key KeyFile
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("decode key: %w", err)
	}
	raw, err := base64.StdEncoding.DecodeString(key.Private)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	if len(raw) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key length")
	}
	return ed25519.PrivateKey(raw), nil
}

