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
