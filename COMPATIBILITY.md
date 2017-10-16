# Compatibility

AegisTrace requires Go 1.22 or newer to build.

Artifact, proof, capability, and checkpoint documents carry explicit format versions. Readers reject unknown artifact fields and versions. Additive workflow source fields require a format and verifier review because strict YAML parsing intentionally rejects unknown configuration.

