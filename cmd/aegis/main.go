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
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func compileCommand(args []string) error {
	fs := flag.NewFlagSet("compile", flag.ContinueOnError)
	output := fs.String("o", "", "output artifact")
	if err := fs.Parse(interspersed(args)); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: aegis compile workflow.yaml [-o workflow.aegis]")
	}
	source, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	artifact, report, err := compiler.Compile(source, time.Now())
	if err != nil {
		return err
	}
	if *output == "" {
		*output = fs.Arg(0) + ".aegis"
	}
	if err := compiler.Save(*output, artifact); err != nil {
		return err
	}
	fmt.Printf("compiled %s\nverified %d steps, %d edges\ndigest   %s\n", *output, report.Steps, report.Edges, artifact.Digest)
	return nil
}

func verifyCommand(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: aegis verify workflow.aegis")
	}
	artifact, err := compiler.Load(args[0])
	if err != nil {
		return err
	}
	report := verifier.Verify(artifact.Workflow)
	if err := report.Error(); err != nil {
		return err
	}
	fmt.Printf("verified %s: %d steps, %d edges, %s\n", artifact.Workflow.Name, report.Steps, report.Edges, artifact.Digest)
	return nil
}

func runCommand(args []string) error {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	journalPath := fs.String("journal", "run.atlog", "journal path")
	inputJSON := fs.String("input", "{}", "JSON input")
	approve := fs.Bool("approve", false, "approve guarded steps")
	if err := fs.Parse(interspersed(args)); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("usage: aegis run workflow.aegis [options]")
	}
	artifact, err := compiler.Load(fs.Arg(0))
	if err != nil {
		return err
	}
	var input map[string]any
	if err := json.Unmarshal([]byte(*inputJSON), &input); err != nil {
		return fmt.Errorf("input: %w", err)
	}
	runID, err := randomID()
	if err != nil {
		return err
	}
	writer, err := journal.Create(*journalPath, runID)
	if err != nil {
		return err
	}
	defer writer.Close()
	registry := provider.NewRegistry()
