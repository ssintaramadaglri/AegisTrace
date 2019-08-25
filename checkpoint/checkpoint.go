// Package checkpoint persists resumable workflow state atomically.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var safeID = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

type State struct {
