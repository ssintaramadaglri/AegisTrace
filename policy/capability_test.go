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
