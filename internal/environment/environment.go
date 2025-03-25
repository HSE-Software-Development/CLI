package environment

import "fmt"

type Env map[string]string

func New() Env {
	return make(map[string]string)
}

func (env Env) Set(variable, value string) {
	env[variable] = value
}

func (env Env) Get(variable string) (string, error) {
	if v, ok := env[variable]; !ok {
		return v, nil
	}
	return "", fmt.Errorf("unknown command: %s", variable)
}