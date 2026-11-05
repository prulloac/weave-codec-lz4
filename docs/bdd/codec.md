# BDD: LZ4 Codec Extension

Behavior cases for the first-party `weave-codec-lz4` binary.

**Spec reference:** [codec-extensions.md](https://github.com/prulloac/weave-docs/blob/main/specs/codec-extensions.md)
**Go tests:** `main_test.go`

---

## Feature: Manifest

### Scenario: JSON capability manifest

- **When** I run `weave-codec-lz4 --describe --json`
- **Then** it returns `weave.extension/v1`
- **And** the name is `lz4` and the kind is `codec`
- **And** it advertises `compress` and `decompress`
- **And** `frame` is `true`

### Scenario: Version output

- **When** I run `weave-codec-lz4 --version`
- **Then** it prints the binary version and `weave.extension/v1`

---

## Feature: Raw Transport

### Scenario: Round-trip

- **Given** arbitrary bytes on stdin
- **When** I run `compress` and then `decompress`
- **Then** the original bytes are returned
- **And** stdout contains only codec bytes

### Scenario: Empty input round-trips

- **Given** empty stdin
- **When** I run `compress` and then `decompress`
- **Then** the output is empty

### Scenario: Corrupt input writes no partial plaintext

- **Given** bytes that are not an LZ4 frame
- **When** I run `decompress`
- **Then** the command exits `1`
- **And** stdout is empty
- **And** stderr contains a diagnostic

### Scenario: Unknown command fails

- **When** I run an unknown subcommand
- **Then** the command exits `1`
- **And** stdout is empty
