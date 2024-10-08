package apinext

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Структура для получения задач из фронтэнда и оправки
type Task struct {
	Id      string `json:"id,omitempty"`
	Date    string `json:"date,omitempty"`
	Title   string `json:"title,omitempty"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

const (
	DateFormat = "20060102"
)

func nextDaily(now, date time.Time, repeat string) (string, error) {
	args := strings.Split(repeat, " ")
	if len(args) != 2 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	v, err := strconv.Atoi(args[1])
	if err != nil {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	if v < 1 || v > 400 {
		return "", fmt.Errorf("uncorrect repeat format")
	}

	date = date.AddDate(0, 0, v)
	for now.After(date) {
		date = date.AddDate(0, 0, v)
	}

	return date.Format(DateFormat), nil
}

func dataAnswer(now, date time.Time) (string, error) {
	date = date.AddDate(1, 0, 0)
	for now.After(date) {
		date = date.AddDate(1, 0, 0)
	}

	return date.Format(DateFormat), nil
}

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("repeat is empty")
	}

	d, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("cannot parse date")
	}

	switch {
	case strings.HasPrefix(repeat, "d"):
		return nextDaily(now, d, repeat)
	case repeat == "y":
		return dataAnswer(now, d)
	}

	return "", fmt.Errorf("unexpected type")
}

func (task *Task) Valid() (bool, error) {
	if task.Id == "" {
		return false, fmt.Errorf("не указан id")
	}
	if id, err := strconv.Atoi(task.Id); err != nil || id < 0 {
		return false, fmt.Errorf("id должен быть положительным числом")
	}

	if task.Title == "" {
		return false, fmt.Errorf("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(DateFormat)
	} else {
		_, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			return false, err
		}
	}

	if task.Repeat != "" {
		_, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
