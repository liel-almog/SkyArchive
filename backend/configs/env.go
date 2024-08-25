package configs

import (
	"fmt"
	"sync"

	"github.com/joho/godotenv"
	"github.com/liel-almog/SkyArchive/backend/errors/apperrors"
)

var (
	initEnvOnce sync.Once

	envs map[string]*string
)

func InitEnv() {
	initEnvOnce.Do(func() {
		envs = make(map[string]*string)
		env, err := godotenv.Read()
		if err != nil {
			panic(err)
		}

		for key, value := range env {
			env := value
			envs[key] = &env
		}
	})
}

func GetAllEnvs() (map[string]string, error) {
	if envs == nil {
		return nil, fmt.Errorf("envs is nil")
	}

	envMap := make(map[string]string)
	for key, value := range envs {

		// Remeber that go reuse the value variable in the loop
		// https://medium.com/swlh/use-pointer-of-for-range-loop-variable-in-go-3d3481f7ffc9
		env := *value
		envMap[key] = env
	}

	return envMap, nil
}

func GetEnv(key string) (string, error) {
	if envs == nil {
		return "", fmt.Errorf("envs is nil")
	}

	value, ok := envs[key]
	if !ok {
		return "", fmt.Errorf("%w key %s not found", apperrors.ErrInvalidEnv, key)
	}

	return *value, nil
}
