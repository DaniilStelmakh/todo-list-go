package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	apinext "github.com/DaniilStelmakh/go_final_project_main/apinext"

	_ "github.com/mattn/go-sqlite3"
)

const (
	selectLimit = 50
)

type Storage struct {
	db *sql.DB
}

// Ф-ия дл ясоздания таблицы
func CreateTable(storagePath string) (*Storage, error) {

	log.Printf("Storage %s\n", storagePath)
	_, err := os.Stat(storagePath)

	var install bool
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}

	if install {

		stmt, err := db.Prepare(`
			CREATE TABLE IF NOT EXISTS scheduler(
			id INTEGER PRIMARY KEY,
			date CHAR(8) NOT NULL,
			title TEXT NOT NULL,
			comment TEXT,
			repeat VARCHAR(128));
			CREATE INDEX IF NOT EXISTS idx_sched_date ON scheduler(date);
			`)

		if err != nil {
			return nil, err
		}

		_, err = stmt.Exec()
		if err != nil {
			return nil, err
		}
	}

	return &Storage{db: db}, nil
}

// Закрываем подключение к БД
func (s *Storage) Close() error {
	return s.db.Close()
}

// Ф-ия для изменения таблицы
func (s *Storage) CreateTask(task apinext.Task) (int, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)"

	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Ф-ия для получения задачи по id
func (s *Storage) GetTaskById(id string) (apinext.Task, error) {
	query := "SELECT * FROM scheduler WHERE id = ?"
	row := s.db.QueryRow(query, id)
	task := apinext.Task{}

	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return apinext.Task{}, fmt.Errorf("failed scan from database: %w", err)
	}
	return task, nil
}

// Ф-ия для получиния задачи
func (s *Storage) GetTasks() ([]apinext.Task, error) {
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"

	tasks := []apinext.Task{}

	rows, err := s.db.Query(query, selectLimit)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return []apinext.Task{}, nil
	}

	for rows.Next() {
		task := apinext.Task{}

		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Ф-ия для получения ближайщих задач
func (s *Storage) GetTasksBySearch(search string) ([]apinext.Task, error) {
	query := "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?"

	tasks := []apinext.Task{}

	rows, err := s.db.Query(query, search, search, selectLimit)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return []apinext.Task{}, nil
	}

	for rows.Next() {
		task := apinext.Task{}

		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Ф-ия для получения задач по дате
func (s *Storage) GetTasksByDate(date string) ([]apinext.Task, error) {
	query := "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT ?"

	tasks := []apinext.Task{}

	rows, err := s.db.Query(query, date, selectLimit)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
		return []apinext.Task{}, nil
	}

	for rows.Next() {
		task := apinext.Task{}

		if err := rows.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// Ф-ия для изменения задач
func (s *Storage) UpdateTask(task *apinext.Task) (int64, error) {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := s.db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.Id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Ф-ия для удаления задач
func (s *Storage) DeleteTaskById(id string) (int64, error) {
	query := `DELETE FROM scheduler WHERE id=?`
	res, err := s.db.Exec(query, id)
	if err != nil {
		return 0, err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}
