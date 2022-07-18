package main

import (
	"bytes"
	"testing"
)

func TestRunCmd(t *testing.T) {
	env := Environment{}

	tests := []struct {
		cmd          []string
		expectedCode int
		expectedOut  string
		expectedErr  string
	}{
		{[]string{"/bin/echo", "ABC"}, 0, "ABC\n", ""},
		{[]string{"invalidcommand"}, -1, "", ""},
		{[]string{"/bin/ls", "/notexists"}, 2, "", "/bin/ls: cannot access '/notexists': No such file or directory\n"},
	}

	for _, test := range tests {
		bufOut := bytes.Buffer{}
		bufErr := bytes.Buffer{}
		bufIn := bytes.Buffer{}
		code := RunCmd(test.cmd, env, &bufOut, &bufErr, &bufIn)
		if code != test.expectedCode {
			t.Errorf("Expected %d, got %d", test.expectedCode, code)
			t.FailNow()
		}
		if bufOut.String() != test.expectedOut {
			t.Errorf("Stdout expected '%v', got '%v'", test.expectedOut, bufOut.String())
			t.FailNow()
		}
		if bufErr.String() != test.expectedErr {
			t.Errorf("Stderr expected '%v', got '%v'", test.expectedErr, bufErr.String())
			t.FailNow()
		}
	}
}
