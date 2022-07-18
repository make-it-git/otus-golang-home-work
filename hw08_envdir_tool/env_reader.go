package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	env := make(Environment)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.Contains(entry.Name(), "=") {
			continue
		}

		f, err := os.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, err
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			return nil, err
		}
		if info.Size() == 0 {
			env[entry.Name()] = EnvValue{"", true}
			continue
		}

		reader := bufio.NewReader(f)
		strBytes, _, err := reader.ReadLine()
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		strBytes = bytes.ReplaceAll(strBytes, []byte{0}, []byte{'\n'})

		str := string(strBytes)
		str = strings.TrimRight(str, " \t\r")

		env[entry.Name()] = EnvValue{str, false}
	}

	return env, nil
}

func prepareEnvironment(env Environment) []string {
	realEnv := make([]string, 0)

	oldEnv := make(map[string]string)
	for _, value := range os.Environ() {
		data := strings.Split(value, "=")
		oldEnv[data[0]] = data[1]
	}

	newEnv := make(map[string]EnvValue)
	for k, v := range env {
		newEnv[k] = v
	}

	for oldK, oldV := range oldEnv {
		found := false
		for newK := range newEnv {
			if oldK == newK {
				found = true
			}
		}

		if !found {
			realEnv = append(realEnv, fmt.Sprintf("%s=%s", oldK, oldV))
		}
	}

	for newK, newV := range newEnv {
		found := false
		for oldK := range oldEnv {
			if newK == oldK {
				found = true
				if !newV.NeedRemove {
					found = false
				}
			}
		}

		if !found {
			realEnv = append(realEnv, fmt.Sprintf("%s=%s", newK, newV.Value))
		}
	}

	return realEnv
}
