# Protocol

AegisTrace does not expose a network daemon in the current release. This document defines the boundaries used by adapters and persisted execution records.

## Provider request

A provider receives a context deadline and a request containing run ID, step ID, optional action, resolved input, and remaining token budget. It returns output, observed token usage, and non-sensitive metadata.

Providers must not retain request credentials in response metadata. They should return errors for non-2xx responses, malformed payloads, or absent expected output.

## MCP

The MCP adapter sends JSON-RPC 2.0 method `tools/call` with the admitted workflow action as `name` and resolved step input as `arguments`. Bearer authentication is optional and supplied only from `AEGIS_MCP_TOKEN`.

## Replay

Replay performs no provider call. It parses all journal events, verifies sequence and hash links, and reconstructs completed step outputs in journal order. A recorded terminal failure is returned as a replay error.

## Stability

Persisted format identifiers are stable within major format versions. Provider response parsing may expand compatibly. A future distributed protocol will be versioned independently and will not redefine artifact or proof semantics.
