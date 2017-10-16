# Compatibility

AegisTrace requires Go 1.22 or newer to build.

Artifact, proof, capability, and checkpoint documents carry explicit format versions. Readers reject unknown artifact fields and versions. Additive workflow source fields require a format and verifier review because strict YAML parsing intentionally rejects unknown configuration.

The command-line interface follows semantic versioning after 1.0. Before 1.0, breaking changes are documented in the changelog. Journal event data may gain fields, while event sequence, type, previous hash, and hash semantics remain stable within the journal generation.
