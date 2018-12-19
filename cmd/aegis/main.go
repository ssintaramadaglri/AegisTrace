package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aegistrace/aegistrace/checkpoint"
	"github.com/aegistrace/aegistrace/compiler"
	"github.com/aegistrace/aegistrace/journal"
	"github.com/aegistrace/aegistrace/proof"
	"github.com/aegistrace/aegistrace/provider"
	aegisruntime "github.com/aegistrace/aegistrace/runtime"
	"github.com/aegistrace/aegistrace/verifier"
	"github.com/aegistrace/aegistrace/workflow"
)

const version = "0.9.1"

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "compile":
		err = compileCommand(os.Args[2:])
	case "verify":
		err = verifyCommand(os.Args[2:])
	case "run":
		err = runCommand(os.Args[2:])
	case "replay":
		err = replayCommand(os.Args[2:])
	case "inspect":
		err = inspectCommand(os.Args[2:])
	case "keygen":
		err = keygenCommand(os.Args[2:])
	case "prove":
		err = proveCommand(os.Args[2:])
	case "verify-proof":
		err = verifyProofCommand(os.Args[2:])
	case "version", "--version", "-v":
		fmt.Println("aegis", version)
		return
	case "help", "--help", "-h":
		usage()
		return
	default:
		usage()
		err = fmt.Errorf("unknown command %q", os.Args[1])
	}
