package main

import (
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
		f := generateFile(largeBuf)

		dest := path.Join(os.TempDir(), "dest.txt")
		err := Copy(f, dest, 0, 0)
		require.Nil(t, err)

		destF, err := os.Open(dest)
		require.Nil(t, err)
		contents, _ := io.ReadAll(destF)
		_ = destF.Close()
		require.Equal(t, largeBuf, contents)
	})

	t.Run("offset copy", func(t *testing.T) {
		largeBuf := make([]byte, 1024*1024)
		f := generateFile(largeBuf)

		dest := path.Join(os.TempDir(), "dest.txt")
		err := Copy(f, dest, 35, 0)
		require.Nil(t, err)

		destF, err := os.Open(dest)
		require.Nil(t, err)
		contents, _ := io.ReadAll(destF)
		_ = destF.Close()
		require.Equal(t, largeBuf[35:], contents)
	})

	t.Run("offset limit copy", func(t *testing.T) {
		largeBuf := make([]byte, 1024*1024)
		f := generateFile(largeBuf)

		dest := path.Join(os.TempDir(), "dest.txt")
		err := Copy(f, dest, 35, 42)
		require.Nil(t, err)

		destF, err := os.Open(dest)
		require.Nil(t, err)
		contents, _ := io.ReadAll(destF)
		_ = destF.Close()
		require.Equal(t, largeBuf[35:35+42], contents)
	})
}

func generateFile(largeBuf []byte) string {
	gen := rand.New(rand.NewSource(time.Now().Unix()))
	gen.Read(largeBuf)
	f, _ := os.CreateTemp("", "copy_test_*.txt")
	_, _ = f.Write(largeBuf)
	_ = f.Close()
	return f.Name()
}
