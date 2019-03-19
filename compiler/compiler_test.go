package compiler_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/aegistrace/aegistrace/compiler"
)

const validWorkflow = `
name: test
version: "1"
budget:
  max_steps: 4
  max_tool_calls: 1
  max_tokens: 50
  timeout: 2s
permissions:
