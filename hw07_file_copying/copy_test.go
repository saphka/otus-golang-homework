package main

import (
	"errors"
	"io"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("bad offset", func(t *testing.T) {
		err := Copy("dummy", "dummy", -1, 0)
		require.Equal(t, err, ErrNegativeOffset)
	})
	t.Run("bad limit", func(t *testing.T) {
		err := Copy("dummy", "dummy", 0, -1)
		require.Equal(t, err, ErrNegativeLimit)
	})

	t.Run("full copy", func(t *testing.T) {
		largeBuf := make([]byte, 1024*1024)
		f, err := generateFile(largeBuf)
		require.Nil(t, err)

		dest := path.Join(os.TempDir(), "dest_full.txt")
		err = Copy(f, dest, 0, 0)
		require.Nil(t, err)

		destF, err := os.Open(dest)
		require.Nil(t, err)
		contents, err := io.ReadAll(destF)
		require.Nil(t, err)
		err = destF.Close()
		require.Nil(t, err)
		require.Equal(t, largeBuf, contents)
	})

	t.Run("offset copy", func(t *testing.T) {
		largeBuf := make([]byte, 1024*1024)
		f, err := generateFile(largeBuf)
		require.Nil(t, err)

		dest := path.Join(os.TempDir(), "dest_offset.txt")
		err = Copy(f, dest, 35, 0)
		require.Nil(t, err)

		destF, err := os.Open(dest)
		require.Nil(t, err)
		contents, err := io.ReadAll(destF)
		require.Nil(t, err)
		err = destF.Close()
		require.Nil(t, err)
		require.Equal(t, largeBuf[35:], contents)
	})

	t.Run("offset limit copy", func(t *testing.T) {
		largeBuf := make([]byte, 1024*1024)
		f, err := generateFile(largeBuf)
		require.Nil(t, err)

		dest := path.Join(os.TempDir(), "dest_offset_limit.txt")
		err = Copy(f, dest, 35, 42)
		require.Nil(t, err)

		destF, err := os.Open(dest)
		require.Nil(t, err)
		contents, err := io.ReadAll(destF)
		require.Nil(t, err)
		err = destF.Close()
		require.Nil(t, err)
		require.Equal(t, largeBuf[35:35+42], contents)
	})

	t.Run("no file copy", func(t *testing.T) {
		dest := path.Join(os.TempDir(), "dest_dummy.txt")
		err := Copy("dummy", dest, 0, 0)
		require.Error(t, err)

		_, err = os.Open(dest)
		require.Error(t, err)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("dir copy", func(t *testing.T) {
		dest := path.Join(os.TempDir(), "dest_dir_1.txt")
		err := Copy(os.TempDir(), dest, 35, 42)
		require.Error(t, err)

		_, err = os.Open(dest)
		require.Error(t, err)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})

	t.Run("dev random copy", func(t *testing.T) {
		dest := path.Join(os.TempDir(), "dest_rand_1.txt")
		err := Copy("/dev/urandom", dest, 35, 42)
		require.Error(t, err)

		_, err = os.Open(dest)
		require.Error(t, err)
		require.True(t, errors.Is(err, os.ErrNotExist))
	})
}

func generateFile(largeBuf []byte) (string, error) {
	gen := rand.New(rand.NewSource(time.Now().Unix()))
	gen.Read(largeBuf)
	f, err := os.CreateTemp("", "copy_test_*.txt")
	if err != nil {
		return "", err
	}
	_, err = f.Write(largeBuf)
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}
