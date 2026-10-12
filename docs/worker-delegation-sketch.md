# Signed worker delegation - sketch

A worker receives a delegated capability signed by the orchestrator, scoped to the workflow and an expiry. The worker journals locally and ships the sealed segment back; the orchestrator verifies the chain link before merging segments.

Open questions: whether delegation chains stay one level deep, and how revocation interacts with an offline worker.