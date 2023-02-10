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
| Was this workflow altered after compilation? | Artifact digest validation |
| Could the graph run forever? | Cycles are rejected before admission |
| Was a tool undeclared? | Exact tool permission verification |
| Was an event changed, removed, or reordered? | Journal sequence and hash-chain validation |
| Is one event included in an attested run? | Signed Merkle inclusion proof |
| Did replay contact an external provider? | Replay has no provider execution path |
| Can this capability invoke the requested scope now? | Ed25519 signature, scope, and expiry checks |

AegisTrace does not certify that a model response is correct, that an allowed tool is safe, or that a compromised host is trustworthy. Those boundaries are explicit so attestations remain meaningful.

## Failure behavior

The runtime fails closed at each boundary:

- Unknown YAML fields stop compilation.
- Modified or unknown artifact formats stop loading.
- Missing inputs and prior-step outputs stop the affected step.
- Provider response bodies larger than 8 MiB are rejected.
- Workflow deadlines cancel in-flight HTTP requests.
- Exhausted step, tool, or token budgets terminate the run.
- Missing approval prevents guarded execution.
- Journal-chain divergence prevents replay and proof creation.
- Invalid signatures, roots, paths, scopes, and expirations prevent verification.

Provider errors may follow an explicitly verified `on_error` edge. Without that edge, the journal records a terminal `run.failed` event.

## Operational model

For local development, one `aegis` process owns the workflow, provider registry, journal writer, and checkpoint directory. This keeps the trusted execution path small and makes the binary easy to inspect and deploy.

For production use, run AegisTrace with a dedicated service identity, separate credentials by provider, keep MCP permissions narrow, encrypt journals and checkpoints at rest, and move attestation keys into a managed signer. The artifact and proof formats do not depend on a particular storage platform, so future distributed workers can preserve the same verification model.

A typical promotion flow is:

```text
source review -> compile -> verify -> sign/release artifact
                                      |
                                      v
                               execute with policy
                                      |
                         journal + checkpoints + proof
                                      |
                                      v
                              audit / replay / CI
```

## Install

AegisTrace requires Go 1.22 or newer.

```sh
go install github.com/aegistrace/aegistrace/cmd/aegis@latest
```

From a checkout:

```sh
go build -o aegis ./cmd/aegis
go test ./...
```

## Quick start

Compile and verify the example:

```sh
aegis compile examples/echo-agent.yaml -o echo.aegis
aegis verify echo.aegis
```

Run it and retain the audit journal:

```sh
aegis run echo.aegis \
  --input '{"message":"Investigate checkout latency"}' \
  --journal run.atlog
```

Replay and inspect without invoking tools:

```sh
aegis replay run.atlog
aegis inspect run.atlog --event 2
```

Create a signed inclusion proof:

```sh
aegis keygen -o signing-key.json
aegis prove run.atlog --event 2 --key signing-key.json -o event.proof.json
aegis verify-proof event.proof.json
```

## Workflow example

```yaml
name: issue-triage
version: "1"
entry: classify
budget:
  max_steps: 8
  max_tool_calls: 1
  max_tokens: 1000
  timeout: 30s
permissions:
  - tool: ticket.create
steps:
  - id: classify
    kind: model
    provider: echo
    input:
      prompt: ${input.message}
  - id: save
    kind: checkpoint
  - id: create
    kind: tool
    provider: echo
    action: ticket.create
    requires_approval: true
    input:
      classification: ${steps.classify}
  - id: result
    kind: emit
    input:
      ticket: ${steps.create}
```

Blank `next` fields are normalized to the following step. Explicit `next` and `on_error` edges are supported. Workflows must form a finite DAG; unbounded loops are rejected.

## Providers

`echo` is always available and makes no network calls. Other adapters are enabled only when their endpoint environment variable is set.

| Provider | Endpoint variable | Credential variable |
|---|---|---|
| OpenAI-compatible | `AEGIS_OPENAI_ENDPOINT` | `OPENAI_API_KEY` |
| Anthropic | `AEGIS_ANTHROPIC_ENDPOINT` | `ANTHROPIC_API_KEY` |
| Gemini | `AEGIS_GEMINI_ENDPOINT` | `GEMINI_API_KEY` |
| Ollama | `AEGIS_OLLAMA_ENDPOINT` | none |
| MCP | `AEGIS_MCP_ENDPOINT` | `AEGIS_MCP_TOKEN` |

Provider responses are capped at 8 MiB. HTTP clients use request deadlines inherited from the workflow budget.

## Security model

A valid artifact proves structural properties; it does not prove that a model response is truthful or a permitted tool is harmless. Provider credentials remain process-local and are never inserted by AegisTrace into workflow artifacts. Journal files can contain sensitive inputs and outputs and are created with owner-only permissions where supported.

Read [SECURITY.md](SECURITY.md) and [THREAT_MODEL.md](THREAT_MODEL.md) before production use.

## Formats and compatibility

- Workflow source: YAML
- Artifact: `aegis/v1` JSON with a SHA-256 workflow digest
- Journal: newline-delimited JSON with per-event chaining
- Proof: `aegis-proof/v1` JSON
- Capability: signed compact token with `aegis-cap/v1` domain separation
- Checkpoint: `aegis-checkpoint/v1` JSON

Versioned formats are documented in [SPECIFICATION.md](SPECIFICATION.md). Unknown artifact fields are rejected to prevent ambiguous interpretation.

## Repository layout

```text
cmd/aegis/      command-line interface
workflow/       portable schema and artifact seal
compiler/       strict source compiler
verifier/       static admission checks
runtime/        budgeted execution and replay
journal/        tamper-evident event log
proof/          Merkle and Ed25519 attestations
policy/         scoped capability tokens
provider/       model, MCP, and local adapters
checkpoint/     atomic state persistence
examples/       runnable workflows
```

## Design principles

1. Verify before execution.
2. Put every external effect behind a provider boundary.
3. Record results, not secrets required to obtain them.
4. Make tampering detectable with standard cryptography.
5. Keep formats small, explicit, and versioned.
6. Prefer reproducibility and clear failure over hidden recovery.

## Development

```sh
go fmt ./...
go vet ./...
go test -race ./...
go build ./cmd/aegis
```

Every source file must remain below 500 lines. Contributions should add focused packages rather than expanding central files indefinitely.

## License

MIT License. See [LICENSE](LICENSE).

<!-- draft note 438 -->
