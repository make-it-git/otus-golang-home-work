package main

import (
	"os"
	"strings"
	"testing"
)

var fixtures = "testdata/env"

func TestReadDir(t *testing.T) {
	env, err := ReadDir(fixtures)
	if err != nil {
		t.Error(err)
	}

	expected := make(map[string]EnvValue, 5)
	expected["BAR"] = EnvValue{"bar", false}
	expected["EMPTY"] = EnvValue{"", false}
	expected["FOO"] = EnvValue{"   foo\nwith new line", false}
	expected["HELLO"] = EnvValue{"\"hello\"", false}
	expected["UNSET"] = EnvValue{"", true}

	for k, v := range env {
		envValue, ok := expected[k]
		if !ok {
			t.Errorf("Expected %s exists", k)
			t.FailNow()
		}
		if envValue.Value != v.Value {
			t.Errorf("Expected key='%s', value='%s' to be equal to '%s'", k, v.Value, envValue.Value)
			t.FailNow()
		}
		if envValue.NeedRemove != v.NeedRemove {
			t.Errorf("Expected key='%s', NeedRemove='%v' to be equal to '%v'", k, v.NeedRemove, envValue.NeedRemove)
			t.FailNow()
		}
		delete(expected, k)
	}

	if len(expected) != 0 {
		t.Errorf("Unexpected values found %v", expected)
	}
}

func TestEnvironmentPersisted(t *testing.T) {
	err := os.Setenv("SHOULD_EXIST", "SOME VALUE")
	if err != nil {
		t.Error(err)
	}

	defer os.Clearenv()

	env, err := ReadDir(fixtures)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	realEnv := prepareEnvironment(env)

	envMap := sliceToMap(realEnv)

	v, ok := envMap["SHOULD_EXIST"]
	if !ok {
		t.Error("Key 'SHOULD_EXIST' not found")
		t.FailNow()
	}
	if v != "SOME VALUE" {
		t.Errorf("Key 'SHOULD_EXIST' contains value '%s', expectedCode 'SOME VALUE'", v)
		t.FailNow()
	}
}

func TestEnvironmentRemoved(t *testing.T) {
	err := os.Setenv("UNSET", "VALUE EXISTS")
	if err != nil {
		t.Error(err)
	}

	defer os.Clearenv()

	env, err := ReadDir(fixtures)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	realEnv := prepareEnvironment(env)

	envMap := sliceToMap(realEnv)

	if _, ok := envMap["UNSET"]; ok {
		t.Error("Key 'UNSET' should be remove")
		t.FailNow()
	}
}

func sliceToMap(s []string) map[string]string {
	res := make(map[string]string)

	for _, v := range s {
		value := strings.Split(v, "=")
		if len(value) == 1 {
			res[value[0]] = ""
			continue
		}
		res[value[0]] = value[1]
	}

	return res
}
