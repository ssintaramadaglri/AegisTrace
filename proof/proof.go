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

func Create(events []journal.Event, index int, private ed25519.PrivateKey) (Attestation, error) {
	if err := journal.Verify(events); err != nil {
		return Attestation{}, err
	}
	if index < 0 || index >= len(events) {
		return Attestation{}, fmt.Errorf("event index out of range")
	}
	leaves, err := eventLeaves(events)
	if err != nil {
		return Attestation{}, err
	}
	path, root := merklePath(leaves, index)
	message := []byte("aegistrace-proof-v1\n" + hex.EncodeToString(root))
	pub := private.Public().(ed25519.PublicKey)
	return Attestation{
		Version: "aegis-proof/v1", EventIndex: index, EventCount: len(events), EventHash: events[index].Hash,
		Root: hex.EncodeToString(root), Path: path, PublicKey: base64.StdEncoding.EncodeToString(pub),
		Signature: base64.StdEncoding.EncodeToString(ed25519.Sign(private, message)),
	}, nil
}

func Verify(a Attestation) error {
	if a.Version != "aegis-proof/v1" || a.EventCount <= 0 || a.EventIndex < 0 || a.EventIndex >= a.EventCount {
		return fmt.Errorf("invalid proof header")
	}
	leaf, err := hex.DecodeString(a.EventHash)
	if err != nil || len(leaf) != sha256.Size {
		return fmt.Errorf("invalid event hash")
	}
	current := leaf
	for _, sibling := range a.Path {
		raw, err := hex.DecodeString(sibling.Hash)
		if err != nil || len(raw) != sha256.Size {
			return fmt.Errorf("invalid sibling hash")
		}
		if sibling.Left {
			current = nodeHash(raw, current)
		} else {
			current = nodeHash(current, raw)
		}
	}
	if hex.EncodeToString(current) != a.Root {
		return fmt.Errorf("Merkle root mismatch")
	}
	pub, err := base64.StdEncoding.DecodeString(a.PublicKey)
	if err != nil || len(pub) != ed25519.PublicKeySize {
		return fmt.Errorf("invalid public key")
	}
	sig, err := base64.StdEncoding.DecodeString(a.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature encoding")
	}
	if !ed25519.Verify(pub, []byte("aegistrace-proof-v1\n"+a.Root), sig) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func eventLeaves(events []journal.Event) ([][]byte, error) {
	out := make([][]byte, len(events))
	for i, e := range events {
		raw, err := hex.DecodeString(e.Hash)
		if err != nil || len(raw) != sha256.Size {
			return nil, fmt.Errorf("event %d has invalid hash", i)
		}
		out[i] = raw
	}
	return out, nil
}

func merklePath(leaves [][]byte, index int) ([]Sibling, []byte) {
	level := leaves
	position := index
	var path []Sibling
	for len(level) > 1 {
		if len(level)%2 == 1 {
			level = append(level, level[len(level)-1])
		}
		sibling := position ^ 1
		path = append(path, Sibling{hex.EncodeToString(level[sibling]), sibling < position})
		next := make([][]byte, 0, len(level)/2)
		for i := 0; i < len(level); i += 2 {
			next = append(next, nodeHash(level[i], level[i+1]))
		}
		level, position = next, position/2
	}
	return path, level[0]
}

func nodeHash(left, right []byte) []byte {
	sum := sha256.Sum256(append(append([]byte{1}, left...), right...))
	return sum[:]
}
