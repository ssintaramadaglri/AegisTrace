# Specification

## Artifact `aegis/v1`

An artifact is JSON containing `magic`, `format`, `digest`, `compiled_at`, and `workflow`. `magic` is `AEGISTRACE`. `digest` is `sha256:` followed by lowercase hexadecimal SHA-256 over Go JSON encoding of the normalized workflow. `compiled_at` is informational and excluded from the digest.
