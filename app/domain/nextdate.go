package domain

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

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
