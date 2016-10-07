// Copyright (c) 2016 Intel Corporation. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package debug_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.intel.com/hpdd/logging/debug"
)

func TestPackageDisable(t *testing.T) {
	// Would be nice to have a BeforeAll()...
	debug.SetOutput(os.Stderr)
	debug.Disable()

	var buf bytes.Buffer
	debug.SetOutput(&buf)

	// Should be disabled by default
	debug.Print(testInputs[0])
	debug.Enable()
	debug.Print(testInputs[1])
	debug.Disable()
	debug.Print(testInputs[3])

	lines := strings.Split(buf.String(), "\n")
	lines = lines[:len(lines)-1] // Don't want the empty line
	if len(lines) != 1 || !strings.HasSuffix(lines[0], testInputs[1]) {
		t.Fatalf("Expected only %s to be logged (found %q)", testInputs[1], lines)
	}
}

func TestPackageWriter(t *testing.T) {
	debug.SetOutput(os.Stderr)
	debug.Disable()

	var buf bytes.Buffer
	debug.SetOutput(&buf)
	debug.Enable()

	debug.Print(testInputs[0])
	writer := debug.Writer().Prefix("writer: ")
	writer.Write([]byte(testInputs[1]))

	lines := strings.Split(buf.String(), "\n")
	lines = lines[:len(lines)-1] // Don't want the empty line
	for i, line := range lines {
		if !strings.HasSuffix(line, testInputs[i]) {
			t.Fatalf("line %d: expected %s, found %s", i-1, testInputs[i], line)
		}
	}
}

func TestPackageSetOutput(t *testing.T) {
	debug.SetOutput(os.Stderr)
	debug.Disable()

	var bufA bytes.Buffer
	var bufB bytes.Buffer
	debug.SetOutput(&bufA)
	debug.Enable()

	debug.Print(testInputs[0])
	writer := debug.Writer().Prefix("writer: ")

	debug.SetOutput(&bufB)
	writer.Write([]byte(testInputs[1]))

	linesA := strings.Split(bufA.String(), "\n")
	linesB := strings.Split(bufB.String(), "\n")

	if len(linesA) < 1 || !strings.HasSuffix(linesA[0], testInputs[0]) {
		t.Fatalf("Output didn't make it to first writer")
	}

	if len(linesB) < 1 || !strings.HasSuffix(linesB[0], testInputs[1]) {
		t.Fatalf("Output didn't make it to second writer")
	}
}

func TestPackageWriterWithLogger(t *testing.T) {
	debug.SetOutput(os.Stderr)
	debug.Disable()

	var buf bytes.Buffer
	debug.SetOutput(&buf)
	debug.Enable()

	debug.Print(testInputs[0])
	writer := debug.Writer().Prefix("writer: ")
	log := log.New(writer, "2nd log: ", 0)
	log.Print(testInputs[1])

	lines := strings.Split(buf.String(), "\n")
	lines = lines[:len(lines)-1] // Don't want the empty line
	for i, line := range lines {
		if !strings.HasSuffix(line, testInputs[i]) {
			t.Fatalf("line %d: expected %s, found %s", i-1, testInputs[i], line)
		}
	}

	prefixedOutput := "writer: 2nd log: " + testInputs[1]
	if !strings.HasSuffix(lines[1], prefixedOutput) {
		t.Fatalf("prefixes wrong: %s", lines[1])
	}
}
