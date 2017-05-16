# Architecture

## Execution pipeline

```text
YAML source -> strict parser -> normalizer -> verifier -> sealed artifact
                                                        |
                                                        v
input -> runtime -> provider boundary -> external model/tool
             |              |
             v              v
       checkpoint       hash-chained journal -> replay / proof
