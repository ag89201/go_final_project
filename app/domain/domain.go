package domain

import (
	"os"
)

func GetEnv(environment string, default_value string) string {
	env := os.Getenv(environment)
	if len(env) > 0 {
		return env
	}
	return default_value
}

func FileNotExists(filePath string) bool {

	_, err := os.Stat(filePath)
	return os.IsNotExist(err)

}

func CreateFile(filePath string) error {
	_, err := os.Create(filePath)
	return err
}
