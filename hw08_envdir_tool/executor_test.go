package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("no changes", func(t *testing.T) {
		env := make(Environment)
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		code := RunCmd([]string{"echo", "Hello", "world"}, env, in, out, err)
		require.Zero(t, code)
		require.Contains(t, "Hello world\n", out.String())
	})

	t.Run("remove var", func(t *testing.T) {
		_, ok := os.LookupEnv("USER")
		require.True(t, ok)
		env := make(Environment)
		env["USER"] = EnvValue{NeedRemove: true}
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		code := RunCmd([]string{"printenv"}, env, in, out, err)
		require.Zero(t, code)
		require.NotContains(t, out.String(), "\nUSER=")
	})

	t.Run("change var", func(t *testing.T) {
		realUser, ok := os.LookupEnv("USER")
		require.True(t, ok)
		env := make(Environment)
		env["USER"] = EnvValue{Value: realUser + "_dummy"}
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		code := RunCmd([]string{"printenv"}, env, in, out, err)
		require.Zero(t, code)
		require.Contains(t, out.String(), fmt.Sprintf("USER=%s_dummy\n", realUser))
	})

	t.Run("add var", func(t *testing.T) {
		_, ok := os.LookupEnv("FOO")
		require.False(t, ok)
		env := make(Environment)
		env["FOO"] = EnvValue{Value: "bar"}
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		code := RunCmd([]string{"printenv"}, env, in, out, err)
		require.Zero(t, code)
		require.Contains(t, out.String(), "FOO=bar")
	})

	t.Run("no command", func(t *testing.T) {
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}
		code := RunCmd([]string{"some_dummy_coomand"}, make(Environment), in, out, err)
		require.NotZero(t, code)
	})
}
