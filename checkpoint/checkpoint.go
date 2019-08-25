// Package checkpoint persists resumable workflow state atomically.
package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
