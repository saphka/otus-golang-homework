package main

import (
	"bufio"
	"bytes"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	result := make(Environment, len(entries))
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		stat, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if stat.Size() == 0 {
			result[entry.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.Open(entry.Name())
		if err != nil {
			return nil, err
		}
		value, err := readFile(file)
		if err != nil {
			return nil, err
		}
		result[entry.Name()] = value
	}
	return result, nil
}

func readFile(file *os.File) (EnvValue, error) {
	defer func() {
		_ = file.Close()
	}()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return EnvValue{}, err
	}

	line = bytes.Replace(line, []byte{0x00}, []byte("\n"), -1)
	result := strings.TrimRight(string(line), " \t")
	return EnvValue{Value: result}, nil
}
