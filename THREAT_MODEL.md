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
