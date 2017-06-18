# Architecture

## Execution pipeline

```text
YAML source -> strict parser -> normalizer -> verifier -> sealed artifact
                                                        |
                                                        v
input -> runtime -> provider boundary -> external model/tool
             |              |
             v              v
       checkpoint       hash-chained journal -> replay / proof
```

The compiler parses one YAML document with unknown-field rejection. Normalization supplies deterministic defaults and sequential edges. The verifier admits only finite, reachable graphs with valid budgets and explicit tool permissions. The artifact stores the normalized workflow and its SHA-256 digest.

The runtime accepts only artifacts whose digest and static properties validate again at load time. It resolves explicit input and prior-step references, enforces step, tool, token, and wall-clock budgets, and invokes all external effects through the provider interface.

## Packages

- `workflow`: versioned schemas and artifact digest
- `compiler`: strict parsing, normalization, artifact persistence
- `verifier`: graph and policy admission
- `runtime`: execution, budgets, expression resolution, replay
- `provider`: isolated effect implementations
- `journal`: append-only event persistence and chain verification
- `checkpoint`: atomic resumable state snapshots
- `proof`: Merkle inclusion paths and Ed25519 signatures
- `policy`: scoped, expiring capability tokens
- `cmd/aegis`: user-facing orchestration

## Determinism

External model and tool services are nondeterministic. A live AegisTrace run therefore cannot guarantee equal outputs across repeated network calls. Instead, the journal makes the state transition history deterministic after recording. Replay validates the complete hash chain and reconstructs outputs without network access.
