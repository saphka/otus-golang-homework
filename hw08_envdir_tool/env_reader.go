package main

import (
	"bufio"
	"bytes"
	"os"
	"path"
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

		value, err := readFile(path.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		result[entry.Name()] = value
	}
	return result, nil
}

func readFile(fileName string) (EnvValue, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return EnvValue{}, err
	}
	defer func() {
		_ = file.Close()
	}()

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return EnvValue{}, err
	}

	line = bytes.ReplaceAll(line, []byte{0x00}, []byte("\n"))
	result := strings.TrimRight(string(line), " \t")
	return EnvValue{Value: result}, nil
}
