# Specification

## Artifact `aegis/v1`

An artifact is JSON containing `magic`, `format`, `digest`, `compiled_at`, and `workflow`. `magic` is `AEGISTRACE`. `digest` is `sha256:` followed by lowercase hexadecimal SHA-256 over Go JSON encoding of the normalized workflow. `compiled_at` is informational and excluded from the digest.

A loader must reject an invalid magic value, unknown format, digest mismatch, unknown JSON field, or workflow that no longer passes static verification.

## Workflow

A workflow has a name, version, entry step, budget, permissions, and steps. Step IDs are unique. Valid kinds are `model`, `tool`, `approval`, `checkpoint`, and `emit`.
