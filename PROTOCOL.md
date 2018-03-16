# Protocol

AegisTrace does not expose a network daemon in the current release. This document defines the boundaries used by adapters and persisted execution records.

## Provider request

A provider receives a context deadline and a request containing run ID, step ID, optional action, resolved input, and remaining token budget. It returns output, observed token usage, and non-sensitive metadata.

Providers must not retain request credentials in response metadata. They should return errors for non-2xx responses, malformed payloads, or absent expected output.

