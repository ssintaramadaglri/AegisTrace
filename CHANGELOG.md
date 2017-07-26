# Changelog

All notable changes follow Keep a Changelog. Versions use semantic versioning.

## [Unreleased]

### Planned

- Resume execution from validated checkpoints
- Encrypted journal envelope format
- OpenTelemetry exporter

## [0.9.1] - 2026-08-31

### Added

- Go-first workflow compiler and static verifier
- Digest-sealed `aegis/v1` artifacts
- Budgeted local runtime with approval and checkpoint steps
- Hash-chained JSONL execution journal and offline replay
- Merkle inclusion proofs with Ed25519 signatures
- Scoped, expiring Ed25519 capability tokens
- OpenAI-compatible, Anthropic, Gemini, Ollama, MCP, and echo providers
- `compile`, `verify`, `run`, `replay`, `inspect`, `keygen`, `prove`, and `verify-proof` commands

### Fixed

- Replay refuses journals whose chain diverges before the requested event
- MCP responses larger than 8 MiB are rejected instead of buffered

## [0.9.0] - 2026-03-27

### Added

- Artifact digest verification on every load path
- Per-step token budgeting and deadline propagation to HTTP providers
- `aegis inspect --event` for reading a single verified event

### Changed

- Unknown artifact fields are rejected as hard errors

## [0.8.1] - 2025-08-08

### Fixed

- Checkpoint writes are atomic on Windows (write-then-rename)
- Journal sequence numbers survive interrupted runs

## [0.8.0] - 2025-01-15

### Added

- `aegis-checkpoint/v1` checkpoint format with crash-safe state capture
- Approval steps that block guarded execution until confirmed

## [0.7.1] - 2024-01-26
