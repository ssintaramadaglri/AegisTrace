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
