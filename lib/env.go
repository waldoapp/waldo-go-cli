package lib

import (
	"os"
	"strings"
)

type Environment map[string]string

//-----------------------------------------------------------------------------

func CurrentEnvironment() Environment {
	env := make(map[string]string)

	for _, item := range os.Environ() {
		pair := strings.SplitN(item, "=", 2)

		if len(pair) != 2 {
			continue
		}

		key := pair[0]
		value := pair[1]

		if len(key) > 0 && len(value) > 0 {
			env[key] = value
		}
	}

	return env
}

//-----------------------------------------------------------------------------

func (env Environment) Flatten() []string {
	var flatEnv []string

	for key, value := range env {
		flatEnv = append(flatEnv, key+"="+value)
	}

	return flatEnv
}
