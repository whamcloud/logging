// Copyright (c) 2021 DDN. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package logging_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/whamcloud/logging"
	"github.com/whamcloud/logging/alert"
	"github.com/whamcloud/logging/audit"
)

func genTestFile(nameOnly bool) (*os.File, error) {
	f, err := ioutil.TempFile("", "logtest-")
	if err != nil {
		return nil, err
	}

	if nameOnly {
		f.Close()
		return f, os.Remove(f.Name())
	}

	return f, nil
}

func TestCreateWriterNewFile(t *testing.T) {
	f, err := genTestFile(true)
	if err != nil {
		t.Fatal(err)
	}

	_, err = logging.CreateWriter(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(f.Name())
}

func TestCreateWriterExistingFile(t *testing.T) {
	// Just test that we append properly
	f, err := genTestFile(false)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte("line1\n"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	w, err := logging.CreateWriter(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.Write([]byte("line2\n"))
	w.(*os.File).Close()

	buf, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("line1\nline2\n")
	if string(buf) != string(expected) {
		t.Fatalf("expected %s, got %s", expected, buf)
	}
	os.Remove(f.Name())
}

func TestSetWriter(t *testing.T) {
	var buf bytes.Buffer
	var testInputs = map[int]string{
		0: "this is an alert!",
		1: "no big deal",
	}

	logging.SetWriter(&buf)

	alert.Warn(testInputs[0])
	audit.Log(testInputs[1])

	lines := strings.Split(buf.String(), "\n")
	lines = lines[:len(lines)-1] // Don't want the empty line
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, found: %q", lines)
	}
	if !strings.HasPrefix(lines[0], "ALERT") || !strings.HasSuffix(lines[0], testInputs[0]) {
		t.Fatalf("alert logging didn't work: %s", lines[0])
	}
	if !strings.HasSuffix(lines[1], testInputs[1]) {
		t.Fatalf("audit logging didn't work: %s", lines[1])
	}
}
