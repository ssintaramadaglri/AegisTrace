# Contributing

Contributions should preserve the small trusted core and explicit format boundaries.

## Workflow

1. Open an issue for format or security-boundary changes.
2. Create a focused branch.
3. Add tests for behavior and failure paths.
4. Run `go fmt ./...`, `go vet ./...`, and `go test -race ./...`.
5. Keep every source file below 500 lines.
6. Update format and security documentation when invariants change.

Use standard-library facilities unless a dependency removes substantial security or interoperability risk. Dependencies must be pinned, maintained, and justified in the pull request.

Public APIs require clear names, stable error behavior, and package documentation. Avoid hidden global state, unchecked type assertions, unbounded reads, and logging secrets. Cryptographic changes require test vectors or independent interoperability tests.

Commits should be reviewable and describe why a change is needed. Generated files must identify their generator and should not be edited manually.
