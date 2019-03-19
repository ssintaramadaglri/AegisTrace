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
