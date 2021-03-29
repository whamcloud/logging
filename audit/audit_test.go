// Copyright (c) 2021 DDN. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package audit_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/whamcloud/logging/audit"
)

var testInputs = map[int]string{
	0: "line1",
	1: "line2",
	2: "line3",
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	a := audit.NewLogger(&buf)

	a.Log(testInputs[0])
	writer := a.Writer().Prefix("writer: ")
	writer.Write([]byte(testInputs[1]))

	lines := strings.Split(buf.String(), "\n")
	lines = lines[:len(lines)-1] // Don't want the empty line
	for i, line := range lines {
		if !strings.HasSuffix(line, testInputs[i]) {
			t.Fatalf("line %d: expected %s, found %s", i-1, testInputs[i], line)
		}
	}
}

func TestSetOutput(t *testing.T) {
	var bufA bytes.Buffer
	var bufB bytes.Buffer
	a := audit.NewLogger(&bufA)

	a.Log(testInputs[0])
	writer := a.Writer().Prefix("writer: ")

	a.SetOutput(&bufB)
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

func TestWriterWithLogger(t *testing.T) {
	var buf bytes.Buffer
	a := audit.NewLogger(&buf)

	a.Log(testInputs[0])
	writer := a.Writer().Prefix("writer: ")
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
