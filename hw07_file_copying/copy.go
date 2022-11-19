package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("offset cannot be negative")
	ErrNegativeLimit         = errors.New("limit cannot be negative")
	ErrFileIsDirectory       = errors.New("file is a directory")
)

func Copy(fromPath, toPath string, offset, limit int64) (finalErr error) {
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}

	source, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = source.Close()
	}()
	var sizeLeft int64
	if sizeLeft, err = positionFile(source, offset); err != nil {
		return err
	}

	dest, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		err := dest.Close()
		if err != nil && finalErr == nil {
			finalErr = err
		}
	}()

	if limit == 0 || limit > sizeLeft {
		limit = sizeLeft
	}
	if err = copyContents(source, dest, limit); err != nil {
		return err
	}
	return nil
}

func positionFile(file *os.File, offset int64) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	if stat.IsDir() {
		return 0, ErrFileIsDirectory
	}
	if !stat.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}

	size := stat.Size()
	if size < offset {
		return 0, ErrOffsetExceedsFileSize
	}
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}
	return size - offset, nil
}

func copyContents(source io.Reader, dest io.Writer, limit int64) error {
	const bufferSize = 1024
	buf := make([]byte, bufferSize)
	var total int64
	var prevProgress int

	for {
		shouldExit := false
		bytesToWrite, err := source.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				shouldExit = true
			} else {
				return err
			}
		}

		total += int64(bytesToWrite)
		if total >= limit {
			bytesToWrite -= int(total - limit)
			shouldExit = true
			total = limit
		}

		if bytesToWrite > 0 {
			_, err = dest.Write(buf[:bytesToWrite])
			if err != nil {
				return err
			}
		}

		progress := int(float32(total) / float32(limit) * 100)
		if progress > prevProgress {
			fmt.Printf("%d%%...", progress)
			prevProgress = progress
		}

		if shouldExit {
			break
		}
	}
	fmt.Println()
	return nil
}
