// Command weave-codec-lz4 is the first-party LZ4 codec extension for Weave.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/pierrec/lz4/v4"
)

const (
	protocolVersion = "weave.extension/v1"
	name            = "lz4"
	kind            = "codec"
	maxInput        = 64 << 20
	maxOutput       = 64 << 20
)

var version = "dev"

type manifest struct {
	V          string   `json:"v"`
	Name       string   `json:"name"`
	Kind       string   `json:"kind"`
	Version    string   `json:"version"`
	Directions []string `json:"directions"`
	Frame      bool     `json:"frame"`
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		return fail(stderr, "usage: weave-codec-lz4 --describe|compress|decompress")
	}
	switch args[0] {
	case "--version", "-v", "version":
		fmt.Fprintf(stdout, "weave-codec-lz4 %s (%s)\n", version, protocolVersion)
		return 0
	case "--describe", "describe":
		return describe(args[1:], stdout, stderr)
	case "compress":
		return compress(stdin, stdout, stderr)
	case "decompress":
		return decompress(stdin, stdout, stderr)
	default:
		return fail(stderr, "unknown command: "+args[0])
	}
}

func describe(args []string, stdout, stderr io.Writer) int {
	jsonMode := false
	for _, arg := range args {
		if arg == "--json" {
			jsonMode = true
		}
	}
	value := manifest{
		V: protocolVersion, Name: name, Kind: kind, Version: version,
		Directions: []string{"compress", "decompress"}, Frame: true,
	}
	if jsonMode {
		if err := json.NewEncoder(stdout).Encode(value); err != nil {
			return fail(stderr, "cannot encode manifest")
		}
		return 0
	}
	fmt.Fprintf(stdout, "name: %s\nversion: %s\ndirections: compress, decompress\nframe: true\n", name, version)
	return 0
}

func compress(stdin io.Reader, stdout, stderr io.Writer) int {
	input, err := readInput(stdin)
	if err != nil {
		return fail(stderr, err.Error())
	}
	var buffer bytes.Buffer
	writer := lz4.NewWriter(&buffer)
	if err := writer.Apply(
		lz4.BlockSizeOption(lz4.Block64Kb),
		lz4.ChecksumOption(false),
	); err != nil {
		return fail(stderr, "cannot configure lz4")
	}
	if _, err := writer.Write(input); err != nil {
		return fail(stderr, "compression failed")
	}
	if err := writer.Close(); err != nil {
		return fail(stderr, "compression failed")
	}
	if _, err := stdout.Write(buffer.Bytes()); err != nil {
		return fail(stderr, "cannot write output")
	}
	return 0
}

func decompress(stdin io.Reader, stdout, stderr io.Writer) int {
	input, err := readInput(stdin)
	if err != nil {
		return fail(stderr, err.Error())
	}
	reader := lz4.NewReader(bytes.NewReader(input))
	output, err := io.ReadAll(io.LimitReader(reader, maxOutput+1))
	if err != nil {
		return fail(stderr, "decompression failed")
	}
	if len(output) > maxOutput {
		return fail(stderr, "decompressed output exceeded the limit")
	}
	if _, err := stdout.Write(output); err != nil {
		return fail(stderr, "cannot write output")
	}
	return 0
}

func readInput(stdin io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(stdin, maxInput+1))
	if err != nil {
		return nil, fmt.Errorf("cannot read input")
	}
	if len(data) > maxInput {
		return nil, fmt.Errorf("input exceeded the limit")
	}
	return data, nil
}

func fail(stderr io.Writer, message string) int {
	fmt.Fprintf(stderr, "weave-codec-lz4: %s\n", message)
	return 1
}
