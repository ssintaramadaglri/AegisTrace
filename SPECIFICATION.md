# Specification

## Artifact `aegis/v1`

An artifact is JSON containing `magic`, `format`, `digest`, `compiled_at`, and `workflow`. `magic` is `AEGISTRACE`. `digest` is `sha256:` followed by lowercase hexadecimal SHA-256 over Go JSON encoding of the normalized workflow. `compiled_at` is informational and excluded from the digest.

A loader must reject an invalid magic value, unknown format, digest mismatch, unknown JSON field, or workflow that no longer passes static verification.

## Workflow

A workflow has a name, version, entry step, budget, permissions, and steps. Step IDs are unique. Valid kinds are `model`, `tool`, `approval`, `checkpoint`, and `emit`.

`model` and `tool` require a provider. `tool` also requires an action admitted by an exact permission or the wildcard permission. Edges use `next` and `on_error`. Every target must exist, every step must be reachable from entry, and cycles are forbidden.

Input values may reference `${input.path}` or `${steps.id}`. References are resolved immediately before execution. Missing references fail closed.

## Journal

A journal is UTF-8 JSONL. Event sequence starts at zero. Each event contains sequence, UTC time, run ID, optional step ID, type, optional data, previous hash, and hash.

The hash is lowercase hexadecimal SHA-256 over the JSON encoding of all event fields except `hash`. The first event has an empty previous hash. Every later event must name the preceding event hash.

## Merkle proof

Leaves are the 32-byte journal event hashes. Internal nodes are SHA-256 of byte `0x01`, the left child, and the right child. An odd final node is duplicated. The root is signed over UTF-8 bytes `aegistrace-proof-v1\n<hex-root>` using Ed25519.

## Capability

A capability has base64url JSON claims and an Ed25519 signature separated by a period. The signature input is `aegis-cap/v1.` followed by the raw claim bytes. Claims include issuer, subject, scopes, issue time, expiry, and random nonce. Maximum accepted issuance lifetime is 24 hours.

## Checkpoint

A checkpoint contains the artifact digest, next step, accumulated outputs, counters, run ID, save time, and format version. Checkpoints are confidential runtime state, not cryptographic proofs.
