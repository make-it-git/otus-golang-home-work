package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrOffsetSeekFailed      = errors.New("offset seek failed")
)

func Copy(fromPath, toPath string, offset, limit int64) (err error) {
	src, err := os.OpenFile(fromPath, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return err
	}

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}
	if !stat.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	ret, err := src.Seek(offset, 0)
	if err != nil {
		return err
	}

	if ret != offset {
		return ErrOffsetSeekFailed
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer dst.Close()

	writtenBytes := int64(0)
	hasLimit := limit > 0

	buf := make([]byte, 1024)
	totalToWrite := stat.Size() - offset
	if hasLimit && totalToWrite > limit {
		totalToWrite = limit
	} else {
		limit = totalToWrite
	}

	pBar := pb.New(int(totalToWrite))
	pBar.SetWriter(os.Stdout)
	pBar.Start()
	barReader := pBar.NewProxyReader(src)

	for {
		nRead, err := barReader.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if nRead == 0 {
			break
		}

		nTake := nRead
		if hasLimit {
			leftToWrite := limit - writtenBytes
			if leftToWrite < int64(nRead) {
				nTake = int(leftToWrite)
			}
		}

		nWrite, err := dst.Write(buf[:nTake])
		if err != nil {
			return err
		}

		writtenBytes += int64(nWrite)

		if hasLimit && writtenBytes >= limit {
			break
		}
	}

	return nil
}
