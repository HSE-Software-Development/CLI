package environment

import (
	"fmt"
	"os"
)

type Env map[string]string

func New() Env {
	env := Env{}
	env_var := []string{
		"PWD", "SHELL", "TERM", "USER", "OLDPWD", "LS_COLORS", "MAIL", "PATH", "LANG", "HOME", "_*",
	}
	for _, v := range env_var {
		cmd := os.Getenv(v)
		env[v] = string(cmd)
	}
	
	return env
}

func (env Env) Set(variable, value string) {
	env[variable] = value
}

func (env Env) Get(variable string) (string, error) {
	if v, ok := env[variable]; ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown command: %s", variable)
}