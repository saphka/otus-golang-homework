package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("add simple var", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "env_load")
		require.NoError(t, err)

		err = generateFile(dir, "FOO", []byte("bar"))
		require.NoError(t, err)

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 1, len(env))
		require.Contains(t, env, "FOO")

	})
}

func generateFile(dir, name string, content []byte) error {
	file, err := os.Create(path.Join(dir, name))
	if err != nil {
		return err
	}
	_, err = file.Write(content)
	if err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
