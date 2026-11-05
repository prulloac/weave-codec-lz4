package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"testing"
)

func invoke(t *testing.T, args []string, input []byte) (int, []byte, string) {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := run(args, bytes.NewReader(input), &stdout, &stderr)
	return code, stdout.Bytes(), stderr.String()
}

func TestDescribeJSON(t *testing.T) {
	code, output, stderr := invoke(t, []string{"--describe", "--json"}, nil)
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	var value manifest
	if err := json.Unmarshal(output, &value); err != nil {
		t.Fatal(err)
	}
	if value.V != protocolVersion || value.Name != "lz4" || value.Kind != "codec" ||
		!value.Frame || len(value.Directions) != 2 {
		t.Fatalf("unexpected manifest: %+v", value)
	}
}

func TestRoundTrip(t *testing.T) {
	input := make([]byte, 1<<20)
	rand.New(rand.NewSource(42)).Read(input)
	code, compressed, stderr := invoke(t, []string{"compress"}, input)
	if code != 0 || stderr != "" {
		t.Fatalf("compress exit=%d stderr=%q", code, stderr)
	}
	code, output, stderr := invoke(t, []string{"decompress"}, compressed)
	if code != 0 || stderr != "" {
		t.Fatalf("decompress exit=%d stderr=%q", code, stderr)
	}
	if !bytes.Equal(output, input) {
		t.Fatal("round trip mismatch")
	}
}

func TestEmptyRoundTrip(t *testing.T) {
	code, compressed, stderr := invoke(t, []string{"compress"}, nil)
	if code != 0 || stderr != "" {
		t.Fatalf("compress exit=%d stderr=%q", code, stderr)
	}
	code, output, stderr := invoke(t, []string{"decompress"}, compressed)
	if code != 0 || stderr != "" || len(output) != 0 {
		t.Fatalf("decompress exit=%d stderr=%q output=%d", code, stderr, len(output))
	}
}

func TestCorruptInputFailsWithoutOutput(t *testing.T) {
	code, output, stderr := invoke(t, []string{"decompress"}, []byte("not lz4"))
	if code != 1 || len(output) != 0 || stderr == "" {
		t.Fatalf("exit=%d output=%d stderr=%q", code, len(output), stderr)
	}
}

func TestUnknownCommand(t *testing.T) {
	code, output, stderr := invoke(t, []string{"bogus"}, nil)
	if code != 1 || len(output) != 0 || stderr == "" {
		t.Fatalf("exit=%d output=%d stderr=%q", code, len(output), stderr)
	}
}

func TestVersion(t *testing.T) {
	original := version
	version = "v9.9.9"
	t.Cleanup(func() { version = original })
	code, output, stderr := invoke(t, []string{"--version"}, nil)
	if code != 0 || stderr != "" {
		t.Fatalf("exit=%d stderr=%q", code, stderr)
	}
	if got, want := string(output), "weave-codec-lz4 v9.9.9 (weave.extension/v1)\n"; got != want {
		t.Fatalf("version = %q, want %q", got, want)
	}
}
