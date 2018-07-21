# Security Policy

## Supported versions

Security fixes are provided for the latest minor release on the default branch. Pre-release formats may change with migration notes.

## Reporting

Do not open a public issue for a suspected vulnerability. Send a minimal reproduction and impact assessment to `security@aegistrace.dev`. Do not include production credentials, private prompts, or customer journal data.

Reports should identify the affected version, operating system, expected boundary, observed behavior, and reproduction steps. Maintainers will acknowledge a complete report within five business days.

## Key handling

`aegis keygen` writes private keys with owner-only permissions where the operating system supports them. Production deployments should use an HSM or managed signing service; the local JSON key format is intended for development and offline signing. Never commit signing keys, provider credentials, capability tokens, journals, or checkpoints.

## Journal confidentiality

Hash chaining provides integrity evidence, not encryption. Journals may contain prompts, model output, tool arguments, and tool results. Protect them with filesystem controls and encrypt storage or transport as required by the deployment.

## Cryptography
