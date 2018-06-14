# Release Process

1. Ensure the worktree is clean and CI passes on the release commit.
2. Run `go mod verify`, `go vet ./...`, `go test -race ./...`, and the CLI smoke flow.
3. Confirm persisted-format changes are backward compatible or versioned.
4. Update `CHANGELOG.md` and the CLI version.
5. Create a signed semantic-version tag.
