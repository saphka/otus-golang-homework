package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
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
		require.Equal(t, env["FOO"], EnvValue{Value: "bar"})
	})

	t.Run("add var + remove", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "env_load")
		require.NoError(t, err)

		err = generateFile(dir, "KEEP", []byte("bar"))
		require.NoError(t, err)
		err = generateFile(dir, "REM", []byte{})
		require.NoError(t, err)

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 2, len(env))
		require.Contains(t, env, "KEEP")
		require.Equal(t, env["KEEP"], EnvValue{Value: "bar"})
		require.Contains(t, env, "REM")
		require.Equal(t, env["REM"], EnvValue{NeedRemove: true})
	})

	t.Run("add var + trim", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "env_load")
		require.NoError(t, err)

		err = generateFile(dir, "TAIL", []byte("bar   "))
		require.NoError(t, err)

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 1, len(env))
		require.Contains(t, env, "TAIL")
		require.Equal(t, env["TAIL"], EnvValue{Value: "bar"})
	})

	t.Run("add var + newline", func(t *testing.T) {
		dir, err := os.MkdirTemp("", "env_load")
		require.NoError(t, err)

		err = generateFile(dir, "NEWLINE", append([]byte("bar"), 0x00, 0x00))
		require.NoError(t, err)

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, 1, len(env))
		require.Contains(t, env, "NEWLINE")
		require.Equal(t, env["NEWLINE"], EnvValue{Value: "bar\n\n"})
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
