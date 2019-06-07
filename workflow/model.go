// Package workflow defines AegisTrace's portable workflow and artifact formats.
package workflow

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

const (
	ArtifactMagic   = "AEGISTRACE"
	ArtifactVersion = "aegis/v1"
)

type Budget struct {
	MaxSteps     int           `json:"max_steps" yaml:"max_steps"`
