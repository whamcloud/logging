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

func TestPackageWriter(t *testing.T) {
	var buf bytes.Buffer
	audit.SetOutput(&buf)

	audit.Log(testInputs[0])
	writer := audit.Writer().Prefix("writer: ")
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
	var bufA bytes.Buffer
	var bufB bytes.Buffer
	audit.SetOutput(&bufA)

	audit.Log(testInputs[0])
	writer := audit.Writer().Prefix("writer: ")

	audit.SetOutput(&bufB)
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
	var buf bytes.Buffer
	audit.SetOutput(&buf)

	audit.Log(testInputs[0])
	writer := audit.Writer().Prefix("writer: ")
	log := log.New(writer, "2nd log: ", 0)
	log.Print(testInputs[1])

	lines := strings.Split(buf.String(), "\n")
	lines = lines[:len(lines)-1] // Don't want the empty line
	for i, line := range lines {
		if !strings.HasSuffix(line, testInputs[i]) {
			t.Fatalf("line %d: expected %s, found %s", i-1, testInputs[i], line)
		}
	}
}
