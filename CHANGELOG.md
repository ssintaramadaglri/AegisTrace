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

### Fixed

- Capability scope comparison is exact, not prefix-based
- Expired tokens are rejected before the provider boundary

## [0.7.0] - 2023-02-10

### Added

- Merkle inclusion proofs (`aegis-proof/v1`) with Ed25519 signatures
- `aegis keygen`, `aegis prove`, and `aegis verify-proof` commands

## [0.6.1] - 2022-03-22

### Fixed

- Cycle detection covers `on_error` edges
- Budget exhaustion emits a terminal `run.failed` event

## [0.6.0] - 2021-03-05

### Added

- MCP adapter behind an explicit endpoint and token
- Ollama and deterministic echo adapters
- Provider response size cap (8 MiB)

## [0.5.0] - 2020-01-17

### Added

- Append-only JSONL journal with a SHA-256 hash chain
- Offline replay that verifies the chain and never calls a provider

## [0.4.0] - 2019-02-21

### Added

- Budgeted runtime: step, tool-call, token, and timeout limits
- `on_error` edges and terminal failure events

## [0.3.0] - 2018-04-09

### Added

- Static verifier: reachability, reference, and tool-permission checks
- Scoped capability tokens signed with Ed25519

## [0.2.0] - 2017-06-15

### Added

- Strict YAML compiler with unknown-field rejection
- SHA-256 workflow digest sealing

## [0.1.0] - 2016-12-02

### Added

- Initial `aegis` CLI: `compile`, `verify`, and `run`
- Portable workflow schema prototype
