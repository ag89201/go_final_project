package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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
	fmt.Printf("Creating file: %s\n", filePath)
	_, err := os.Create(filePath)
	return err
}

func GetNextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("repeat is empty string")
	} else if strings.Contains(repeat, "d ") {
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil {
			return "", err
		}
		if days > 400 {
			return "", errors.New("days is greater than 400")
		}

		pdate, err := time.Parse("20060102", date)
		if err != nil {
			return "", err
		}
		newDate := pdate.AddDate(0, 0, days)

		for newDate.Before(now) {
			newDate = newDate.AddDate(0, 0, days)
		}
		return newDate.Format("20060102"), nil
	} else if repeat == "y" {
		pdate, err := time.Parse("20060102", date)
		if err != nil {
			return "", err
		}
		newDate := pdate.AddDate(1, 0, 0)
		for newDate.Before(now) {
			newDate = newDate.AddDate(1, 0, 0)
		}
		return newDate.Format("20060102"), nil
	} else {
		return "", errors.New("repeat is not valid")
	}

}
