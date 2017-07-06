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
