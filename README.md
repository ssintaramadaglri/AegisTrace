<div align="center">

# AegisTrace

**Verifiable execution for durable AI agents.**  
Compile policy-bound workflows, record every effect, replay without the network, and prove individual events.

[![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)](go.mod)
[![CI](https://github.com/aegistrace/aegistrace/actions/workflows/ci.yml/badge.svg)](https://github.com/aegistrace/aegistrace/actions/workflows/ci.yml)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)
[![Format](https://img.shields.io/badge/artifact-aegis%2Fv1-6f42c1)](SPECIFICATION.md)
[![Python](https://img.shields.io/badge/Python-not%20required-success)](go.mod)

</div>

---

AegisTrace is a verifiable runtime for durable AI workflows. It compiles readable YAML into digest-sealed artifacts, rejects unsafe workflow graphs before execution, records every state transition in a hash-chained journal, and produces independently verifiable Merkle and Ed25519 proofs.

The runtime treats model and tool calls as effects. Live runs invoke configured providers; replay reads recorded results and never contacts a provider. This gives operators a stable audit record and developers a reproducible debugging path without pretending that external models are deterministic.

## Status

AegisTrace is an early, working release. The local runtime, compiler, verifier, journal, replay, checkpoints, proof system, capability tokens, HTTP model adapters, MCP adapter, and CLI are implemented. Distributed scheduling, database backends, and OpenTelemetry export are roadmap items.

## How it works

<div align="center">
  <img src="docs/assets/pipeline.svg" alt="AegisTrace pipeline: YAML source is compiled and statically verified into a sealed artifact, then executed while recording a hash-chained journal, checkpoints, and a signed proof, with a separate offline replay path" width="760" />
</div>

Source is compiled with unknown-field rejection, statically verified, and sealed into a digest-locked artifact. At run time the runtime drives every external effect through the provider boundary, writes each transition into the hash-chained journal, and can produce a signed Merkle proof. Replay reconstructs the run from the journal and never reaches a provider.

| Command | Purpose |
|---|---|
| `aegis compile` | Parse, normalize, verify, and seal a workflow |
| `aegis verify` | Recheck artifact digest and static invariants |
| `aegis run` | Execute with budgets, approvals, checkpoints, and journaling |
| `aegis replay` | Validate and reconstruct a run without network access |
| `aegis inspect` | Read a verified event or complete journal |
| `aegis keygen` | Generate a local Ed25519 attestation key |
| `aegis prove` | Create a signed Merkle inclusion proof |
| `aegis verify-proof` | Verify a proof independently of its journal |

## Features

- Strict YAML compiler with unknown-field rejection
- SHA-256 sealed, portable JSON artifacts
- Static validation of graph reachability, cycles, budgets, references, and tool permissions
- Budgeted execution with model, tool, approval, checkpoint, and emit steps
- Atomic local checkpoints for crash-safe state capture
- Append-only JSONL journal with a SHA-256 hash chain
- Offline replay with chain verification
- Merkle inclusion proofs signed with Ed25519
- Short-lived Ed25519 capability tokens with scoped permissions
- OpenAI-compatible, Anthropic, Gemini, Ollama, MCP, and deterministic local adapters
- Zero Python runtime dependency

## What AegisTrace proves

AegisTrace separates claims it can verify from properties that remain outside the runtime boundary.

<div align="center">
  <img src="docs/assets/proof-chain.svg" alt="AegisTrace journal chain and Merkle proof: an altered event breaks the chain, a valid inclusion path verifies without the journal, replay contacts no provider, and an expired capability is rejected" width="700" />
</div>

| Question | Evidence |
|---|---|
