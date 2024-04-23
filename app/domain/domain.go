package domain

import (
	"os"
	"path/filepath"
)

func GetEnv(environment string, default_value string) string {
	env := os.Getenv(environment)
	if len(env) > 0 {
		return env
	}
	return default_value
}

func GetFileName(environment string, default_name string) (string, error) {
	appPath, err := os.Executable()
	if err != nil {
		return "", err
	}
	fileName := filepath.Join(filepath.Dir(appPath), default_name)
	return GetEnv(environment, fileName), nil
}

func FileNotExists(filePath string) bool {

	_, err := os.Stat(filePath)
	return os.IsNotExist(err)

}

func CreateFile(filePath string) error {
	_, err := os.Create(filePath)
	return err
}
