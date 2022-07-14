package main

import (
	"errors"
	"io/ioutil"
	"syscall"
	"testing"
)

func TestCopy(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Errorf("%s", err)
	}

	defer syscall.Unlink(tmpFile.Name())

	err = ioutil.WriteFile(tmpFile.Name(), []byte("abcdef"), 0644)
	if err != nil {
		t.Errorf("%s", err)
	}

	dst, err := ioutil.TempFile("", "dst")

	t.Run("Expect offset more than file size to fail", func(t *testing.T) {
		err = Copy(tmpFile.Name(), dst.Name(), 100, 0)
		if !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Errorf("%s", err)
		}
	})

	t.Run("Expect invalid file type not to process", func(t *testing.T) {
		err = Copy("/dev/urandom", dst.Name(), 0, 0)
		if !errors.Is(err, ErrUnsupportedFile) {
			t.Errorf("%s", err)
		}
	})
}
