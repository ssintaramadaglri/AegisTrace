# Specification

## Artifact `aegis/v1`

An artifact is JSON containing `magic`, `format`, `digest`, `compiled_at`, and `workflow`. `magic` is `AEGISTRACE`. `digest` is `sha256:` followed by lowercase hexadecimal SHA-256 over Go JSON encoding of the normalized workflow. `compiled_at` is informational and excluded from the digest.

A loader must reject an invalid magic value, unknown format, digest mismatch, unknown JSON field, or workflow that no longer passes static verification.

## Workflow

A workflow has a name, version, entry step, budget, permissions, and steps. Step IDs are unique. Valid kinds are `model`, `tool`, `approval`, `checkpoint`, and `emit`.

`model` and `tool` require a provider. `tool` also requires an action admitted by an exact permission or the wildcard permission. Edges use `next` and `on_error`. Every target must exist, every step must be reachable from entry, and cycles are forbidden.

Input values may reference `${input.path}` or `${steps.id}`. References are resolved immediately before execution. Missing references fail closed.

