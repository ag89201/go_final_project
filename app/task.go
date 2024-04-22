package model

import (
	"errors"
	"time"

	"github.com/ag89201/go_final_project/app/domain"
)

const (
	DateFormat = "20060102"
	LimitTask  = 10
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type IdResponse struct {
	Id int `json:"id"`
}

type TaskResponse struct {
	Tasks []Task `json:"tasks"`
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (t *Task) CheckCorrectData() error {
	if t.Title == "" {
		return errors.New("title is required")
	}

	if len(t.Date) == 0 {
		t.Date = time.Now().Format(DateFormat)
	}

	now, err := time.Parse(DateFormat, t.Date)
	if err != nil {
		return errors.New("invalid date format")
	}

	if now.Before(time.Now()) {
		t.Date = time.Now().Format(DateFormat)
	}

	if len(t.Repeat) > 0 {
		_, err := domain.GetNextDate(now, t.Date, t.Repeat)
		if err != nil {
			// return errors.New("invalid repeat format")
			return err
		}

	}

	return nil
}
