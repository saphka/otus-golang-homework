package main

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("no changes", func(t *testing.T) {
		env := make(Environment)
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		RunCmd([]string{"echo", "Hello", "world"}, env, in, out, err)
		require.Contains(t, "Hello world\n", out.String())
	})

	t.Run("remove var", func(t *testing.T) {
		env := make(Environment)
		env["USER"] = EnvValue{NeedRemove: true}
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		RunCmd([]string{"printenv"}, env, in, out, err)
		require.NotContainsf(t, out.String(), "USER=", "USER variable must be removed")
	})

	t.Run("change var", func(t *testing.T) {
		env := make(Environment)
		realUser, ok := os.LookupEnv("USER")
		require.True(t, ok)
		env["USER"] = EnvValue{Value: realUser + "_dummy"}
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		RunCmd([]string{"printenv"}, env, in, out, err)
		require.Containsf(t, out.String(), fmt.Sprintf("USER=%s_dummy\n", realUser), "Env must contain new user name")
	})

	t.Run("change var", func(t *testing.T) {
		env := make(Environment)
		_, ok := os.LookupEnv("FOO")
		require.False(t, ok)
		env["FOO"] = EnvValue{Value: "bar"}
		in := bytes.NewReader([]byte{})
		out := &strings.Builder{}
		err := &strings.Builder{}

		RunCmd([]string{"printenv"}, env, in, out, err)
		require.Containsf(t, out.String(), "FOO=bar", "Env must contain FOO")
	})
}
