package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

var (
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrNegativeOffset        = errors.New("offset cannot be negative")
	ErrNegativeLimit         = errors.New("limit cannot be negative")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if offset < 0 {
		return ErrNegativeOffset
	}
	if limit < 0 {
		return ErrNegativeLimit
	}

	var (
		source, dest *os.File
		err          error
	)

	source, err = os.Open(fromPath)
	defer func(source io.Closer) {
		_ = source.Close()
	}(source)
	if err != nil {
		return fmt.Errorf("cannot open file %s: %w", fromPath, err)
	}
	var sizeLeft int64
	if sizeLeft, err = seekOffset(source, offset); err != nil {
		return fmt.Errorf("error during seek: %w", err)
	}

	dest, err = os.Create(toPath)
	defer func(dest io.Closer) {
		_ = dest.Close()
	}(dest)
	if err != nil {
		return fmt.Errorf("cannot create file %s: %w", toPath, err)
	}

	if limit == 0 || limit > sizeLeft {
		limit = sizeLeft
	}
	if err = copyContents(source, dest, limit); err != nil {
		return fmt.Errorf("error in copy: %w", err)
	}
	return nil
}

func seekOffset(file *os.File, offset int64) (int64, error) {
	stat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("cannot access file details %s: %w", file.Name(), err)
	}
	size := stat.Size()
	if size < offset {
		return 0, ErrOffsetExceedsFileSize
	}
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, fmt.Errorf("cannot seek to offset %d: %w", offset, err)
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
	return nil
}
