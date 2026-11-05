# AGENTS.md

`weave-codec-lz4` is the first-party LZ4 codec extension for Weave. It is a standalone
`weave-codec-<name>` binary, not a library: it reads raw bytes from stdin, writes raw bytes to stdout
and never imports `weave` or touches a database.

## Commands

- Build: `go build -o weave-codec-lz4 .`
- All tests: `go test ./...`
- Race: `go test -race -count=1 -shuffle=on ./...`
- Vet / format: `go vet ./...`; `gofmt -l .` must print nothing
- Local describe: `go run . --describe --json`
- Version: `go run . --version`
- Full gate: `./scripts/verify.sh`
- User-level install: `./scripts/install.sh --help` (or `./scripts/install.ps1` on Windows)

## Protocol

- `weave-codec-lz4 --describe --json` returns `weave.extension/v1` with name `lz4`, kind `codec`,
  directions `compress` and `decompress`, and `frame: true`.
- `weave-codec-lz4 compress` reads raw bytes on stdin and writes an LZ4 frame to stdout.
- `weave-codec-lz4 decompress` reads an LZ4 frame on stdin and writes raw bytes to stdout.
- Diagnostics go to stderr; stdout contains only codec bytes.
- Exit `0` on success and `1` on failure. No partial output on a failed decode.
- Input and output are bounded to 64 MiB.

## Invariants

1. **Raw transport:** do not add a JSON envelope, base64 or framing around codec bytes.
2. **No database or bundle access:** the codec is pure data in, data out.
3. **Deterministic output:** the same input and version produce the same LZ4 frame.
4. **No secrets in output:** errors never include input bytes.
5. **No core imports:** the binary builds and tests from this repository alone.
6. **User-level install:** installers must never require `sudo` or administrator rights.

## Testing quirks

- Tests call `run` directly with in-memory stdin/stdout; they never spawn the built binary.
- `TestCorruptInputFailsWithoutOutput` pins the no-partial-output rule.
- `TestVersion` pins `weave-codec-lz4 <tag> (weave.extension/v1)`.

## Release

- Tags use `vX.Y.Z`; `.github/workflows/release.yml` builds six static archives, checksums and a
  draft GitHub Release.
- Keep the binary name `weave-codec-lz4`; Weave discovers codecs by the `weave-codec-` prefix.
- The user-level installers default to `~/.local/bin` or `%LOCALAPPDATA%\Programs\weave\bin`.

## Documentation split

- This repository owns install/usage docs and codec behavior examples.
- The codec contract lives in
  [`weave-docs/specs/codec-extensions.md`](https://github.com/prulloac/weave-docs/blob/main/specs/codec-extensions.md).
