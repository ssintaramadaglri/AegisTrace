# Contributing

Contributions should preserve the small trusted core and explicit format boundaries.

## Workflow

1. Open an issue for format or security-boundary changes.
2. Create a focused branch.
3. Add tests for behavior and failure paths.
4. Run `go fmt ./...`, `go vet ./...`, and `go test -race ./...`.
5. Keep every source file below 500 lines.
6. Update format and security documentation when invariants change.
