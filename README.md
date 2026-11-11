# weave-codec-lz4

First-party LZ4 codec extension for [Weave](https://github.com/prulloac/weave).

Weave ships `none` and `zstd` only. This binary supplies the `lz4` codec through the raw
`weave-codec-<name>` protocol: `compress` reads bytes from stdin, `decompress` writes bytes to stdout,
and `--describe --json` advertises the capability manifest.

## Install

User-level installation, no admin or sudo required. Linux and macOS:

```sh
curl -fsSL https://raw.githubusercontent.com/prulloac/weave-codec-lz4/main/scripts/install.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/prulloac/weave-codec-lz4/main/scripts/install.ps1 | iex
```

The installer downloads the latest release, verifies its SHA-256 checksum and installs to
`~/.local/bin` (Linux/macOS) or `%LOCALAPPDATA%\Programs\weave\bin` (Windows), then adds that
directory to the user PATH.

For a managed environment:

```sh
curl -fsSL https://raw.githubusercontent.com/prulloac/weave-codec-lz4/main/scripts/install.sh \
  | sh -s -- --bin-dir "$HOME/bin" --no-modify-path
```

Options: `--version <tag>`, `--bin-dir <dir>`, `--no-modify-path`; environment equivalents
`WEAVE_CODEC_LZ4_VERSION`, `WEAVE_CODEC_LZ4_BIN_DIR` and `WEAVE_CODEC_LZ4_NO_MODIFY_PATH=1`.
The PowerShell script accepts `-Version`, `-BinDir` and `-NoPath`.

Put the binary on `$PATH` or in `<bundle>/.weave/bin/`; a bundle-local binary wins over `$PATH`.

## Use

```sh
weave codecs list ~/notes --json
weave ingest ~/notes --assets --asset-types application/pdf --codec lz4
weave show ~/notes <id> --out file.pdf
```

- Core stores the `lz4` name in `body_codec` and invokes this binary for compression and decoding.
- Existing LZ4 rows from Weave v1.0 require this extension to migrate to schema v7.
- `--codec lz4` fails before writing when the extension is absent.

## Protocol

```sh
weave-codec-lz4 --describe --json
weave-codec-lz4 compress < input > output.lz4
weave-codec-lz4 decompress < output.lz4 > input
```

Manifest:

```json
{
  "v": "weave.extension/v1",
  "name": "lz4",
  "kind": "codec",
  "version": "v0.8.0",
  "directions": ["compress", "decompress"],
  "frame": true
}
```

Input and output are bounded to 64 MiB. A failed decode writes no partial plaintext.

## Development

```sh
./scripts/verify.sh
go build -o weave-codec-lz4 .
```

Behavior examples live in [`docs/bdd/codec.md`](docs/bdd/codec.md). The authoritative contract is
[`codec-extensions.md`](https://github.com/prulloac/weave-docs/blob/main/specs/codec-extensions.md).

## Issues

This repository does not accept issues. Report first-party defects and documentation errors in
[`prulloac/weave`](https://github.com/prulloac/weave/issues/new/choose); security reports use private
vulnerability reporting there.

## License

MIT. See [`LICENSE`](LICENSE).
