package model

import (
	"database/sql"
	"errors"

	_ "modernc.org/sqlite"
)

type Db struct {
	db *sql.DB
}

func NewDB(db *sql.DB) Db {
	return Db{db: db}
}

var Database Db

func NewDataBase(filePath string) (Db, error) {
	db, err := sql.Open("sqlite", filePath)
	if err != nil {
		return Db{}, err
	}

	return NewDB(db), nil
}

func (d Db) Close() error {
	return d.db.Close()
}

func (s Db) CreateSchedulerTable() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT
		)`)
	if err != nil {

		return err
	}
	return nil
}

func (s Db) CreateIndex() error {
	_, err := s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date)`)
	return err

}

func (s Db) InsertTask(task Task) (int, error) {
	res, err := s.db.Exec(`INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s Db) GetTasks() ([]Task, error) {
	var tasks []Task
	rows, err := s.db.Query(`SELECT id, date, title, comment, repeat FROM scheduler LIMIT :limit`, sql.Named("limit", LimitTask))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)

	}

	return tasks, nil
}

func (s Db) GetTask(id int) (Task, error) {
	var task Task
	err := s.db.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id`, sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}
	return task, nil
}

func (s Db) UpdateTask(task Task) (int64, error) {

	res, err := s.db.Exec(`UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`,
		sql.Named("id", task.ID),
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func (s Db) DeleteTask(id int) error {

	res, err := s.db.Exec(`DELETE FROM scheduler WHERE id = :id`, sql.Named("id", id))
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {

		return errors.New("task not found")
	}

	return nil
}
