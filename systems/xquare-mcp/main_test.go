package main

import (
	"bytes"
	"testing"
)

func TestRunWritesSystemName(t *testing.T) {
	t.Parallel()

	var output bytes.Buffer
	if err := run(&output); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	const want = systemName + "\n"
	if got := output.String(); got != want {
		t.Errorf("run() output = %q, want %q", got, want)
	}
}
