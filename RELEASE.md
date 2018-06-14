# Release Process

1. Ensure the worktree is clean and CI passes on the release commit.
2. Run `go mod verify`, `go vet ./...`, `go test -race ./...`, and the CLI smoke flow.
