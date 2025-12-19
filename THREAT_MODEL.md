# Threat Model

## Protected assets

- Workflow intent and admitted permissions
- Integrity and ordering of execution records
- Signing keys and short-lived capabilities
- Provider credentials
- Checkpoint and journal contents

## Adversaries

AegisTrace considers malicious workflow authors, compromised or misleading model providers, untrusted MCP/tool servers, users who can modify stored artifacts or journals, and network attackers handled by the configured HTTPS transport.

## Guarantees

- Modified artifacts fail digest validation.
- Invalid graph structure and undeclared tool actions fail admission.
- Deleted, reordered, inserted, or modified journal events are detectable when the chain is verified.
- A valid proof binds an event hash to a signed journal Merkle root.
- Capabilities cannot be altered without invalidating their signature and are checked for expiry and required scope.
- Provider calls respect workflow context deadlines and bounded response reads.

## Non-guarantees

- Journal confidentiality without encrypted storage
- Truthfulness or safety of model output
- Safety of a tool explicitly permitted by policy
- Protection after signing-key compromise
- Availability of external services
- Deterministic live output from external providers
- Host integrity when the operating system is compromised

## Primary controls

Artifacts are verified both at compilation and load. Tool permissions are exact by default. Runtime effects cross a narrow provider interface. Provider credentials come from environment configuration rather than artifacts. Journals use owner-only creation permissions where supported. Proof and capability signatures use protocol-specific domain separation.

## Deployment guidance

Run workers with least-privilege service identities, separate model and tool credentials, deny unrestricted MCP tools, isolate untrusted tools, redact data before provider calls, encrypt journal storage, rotate signing keys, cap workflow budgets, and require human approval before irreversible effects.
