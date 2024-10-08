package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestPrintVersion(t *testing.T) {
	buildVersion = "1.2.3"
	buildDate = "2023-02-15"
	buildCommit = "abcdefg"
	expectedOutput := "Build version: 1.2.3\nBuild date: 2023-02-15\nBuild commit: abcdefg\n"
	testPrintVersion(t, expectedOutput)
}

func testPrintVersion(t *testing.T, expectedOutput string) {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printVersion()

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	os.Stdout = oldStdout
	w.Close()
	out := <-outC
	if out != expectedOutput {
		t.Errorf("printVersion() = %q, want %q", out, expectedOutput)
	}
}
