package utils

import "os"

type EnvManager struct {
	defaultValue map[string]string
}

func NewEnvManager() (*EnvManager, func()) {
	m := &EnvManager{
		defaultValue: make(map[string]string),
	}

	destructor := func() {
		for k, v := range m.defaultValue {
			_ = os.Setenv(k, v)
		}
	}

	return m, destructor
}

func (t *EnvManager) Set(key, value string) {
	oldV := os.Getenv(key)
	t.defaultValue[key] = oldV
	_ = os.Setenv(key, value)
}
