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

