package tasks

import (
	"fmt"
	"regexp"
	"time"

	apinext "github.com/DaniilStelmakh/go_final_project_main/apinext"
)

type Service struct {
	store StoreJobs
}

// Интерфейс для реализации БД
type StoreJobs interface {
	CreateTask(task *apinext.Task) (int, error)
	GetTaskById(id string) (*apinext.Task, error)
	GetTasks() ([]apinext.Task, error)
	GetTasksBySearch(search string) ([]apinext.Task, error)
	GetTasksByDate(date string) ([]apinext.Task, error)
	UpdateTask(*apinext.Task) (int64, error)
	DeleteTaskById(id string) (int64, error)
}

func New(store StoreJobs) *Service {
	return &Service{store: store}
}

func (s *Service) Add(date, title, comment, repeat string) (int, error) {
	_, err := time.Parse(apinext.DateFormat, date)
	if err != nil {
		return 0, err
	}
	n := time.Now().Format(apinext.DateFormat)

	if date < n {
		if repeat != "" {
			date, err = apinext.NextDate(time.Now(), date, repeat)
			if err != nil {
				return 0, err
			}
		} else {
			date = time.Now().Format(apinext.DateFormat)
		}
	}

	task := &apinext.Task{
		Title:   title,
		Date:    date,
		Comment: comment,
		Repeat:  repeat,
	}

	return s.store.CreateTask(task)
}

func (s *Service) GetAll(search string) ([]apinext.Task, error) {

	if search != "" {
		_, err := regexp.Compile(`[0-3][0-9]\.[0-1][0-9]\.20[0-9][0-9]`)
		if err != nil {
			return nil, err
		}
		matched, _ := regexp.MatchString(`[0-3][0-9]\.[0-1][0-9]\.20[0-9][0-9]`, search)
		if matched {
			date, err := time.Parse("02.01.2006", search)
			if err != nil {
				return nil, err
			}
			search = date.Format(apinext.DateFormat)
			return s.store.GetTasksByDate(search)
		}

		search = "%" + search + "%"
		return s.store.GetTasksBySearch(search)
	}

	return s.store.GetTasks()
}

func (s *Service) Get(id string) (*apinext.Task, error) {
	task, err := s.store.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("not found")
	}

	return task, nil
}

func (s *Service) Update(task *apinext.Task) error {
	count, err := s.store.UpdateTask(task)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows with such id")
	}
	return nil
}

func (s *Service) Delete(id string) error {
	count, err := s.store.DeleteTaskById(id)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("nothing to delete")
	}

	return nil
}

func (s *Service) Done(id string) error {
	task, err := s.Get(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err = s.Delete(task.Id)
		if err != nil {
			return err
		}
		return nil
	}

	date, err := apinext.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = date
	err = s.Update(task)
	if err != nil {
		return err
	}

	return nil
}
