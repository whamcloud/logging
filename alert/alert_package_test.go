// Copyright (c) 2016 DDN. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package alert_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/intel-hpdd/logging/alert"
)

func TestPackageWriter(t *testing.T) {
	var buf bytes.Buffer
	alert.SetOutput(&buf)

	alert.Warn(testInputs[0])
	writer := alert.Writer().Prefix("writer: ")
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
	alert.SetOutput(&bufA)

	alert.Warn(testInputs[0])
	writer := alert.Writer().Prefix("writer: ")

	alert.SetOutput(&bufB)
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
	alert.SetOutput(&buf)

	alert.Warn(testInputs[0])
	writer := alert.Writer().Prefix("writer: ")
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
