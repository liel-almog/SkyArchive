package configs

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
)

var (
	initEnvOnce sync.Once

	envs map[string]*string
)

func InitEnv() {
	envs = make(map[string]*string)
	initEnvOnce.Do(func() {
		env, err := godotenv.Read()
		if err != nil {
			panic(err)
		}

		for key, value := range env {
			envs[key] = &value
		}
	})
}

func GetAllEnvs() (map[string]string, error) {
	if envs == nil {
		return nil, fmt.Errorf("envs is nil")
	}

	envMap := make(map[string]string)
	for key, value := range envs {
		envMap[key] = *value
	}

	return envMap, nil
}

func GetEnv(key string) (string, error) {
	if envs == nil {
		return "", fmt.Errorf("envs is nil")
	}

	value, ok := envs[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}

	return *value, nil
}
