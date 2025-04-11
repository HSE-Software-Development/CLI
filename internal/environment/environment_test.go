package environment

import (
	"testing"
	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestConstructor(t *testing.T) {
	env := New()
	assert.Equal(t, len(env), 11)
	fmt.Println("------")
	fmt.Println(env["PWD"])
	fmt.Println("------")
}


func TestEnv(t *testing.T) {
	env := New()
	variables := map[string]string {
		"111": "xxx",
		"222": "yyy",
		"333": "zzz",
	}
	for k, v := range variables {
		env.Set(k, v)
	}
	for k, v := range variables {
		if val, err := env.Get(k); err != nil {
			assert.Error(t, err)
		} else {
			assert.Equal(t, v, val)
		}
	}
}