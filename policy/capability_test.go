package policy_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"testing"
	"time"

	"github.com/aegistrace/aegistrace/policy"
)

func TestCapabilityScopeAndExpiry(t *testing.T) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 8, 31, 2, 0, 0, 0, time.UTC)
