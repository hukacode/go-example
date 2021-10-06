package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Content     string    `json:"content"`
	IsCompleted bool      `json:"is-completed"`
}

type TaskModel struct {
	DB *sql.DB
}

func (tm TaskModel) Create(task *Task) error {
	query := `INSERT INTO task (content)
						VALUES ($1)
						RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return tm.DB.QueryRowContext(ctx, query, task.Content).Scan(&task.ID, &task.CreatedAt)
}

func (tm TaskModel) Read(id int64) (*Task, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, content
						FROM task
						WHERE id = $1`

	var task Task

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := tm.DB.QueryRowContext(ctx, query, id).Scan(&task.ID, &task.CreatedAt, &task.Content)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &task, nil
}

func (tm TaskModel) ReadAll(content string, isCompleted bool, filter Filter) ([]*Task, Metadata, error) {
	query := fmt.Sprintf(`SELECT count(*) OVER(),id, created_at, content, is_completed
						FROM task
						WHERE (to_tsvector('simple', content) @@ plainto_tsquery('simple', $1) OR $1 = '')
						AND (is_completed = $2 OR $2 = false)
						ORDER BY %s %s, id ASC
						LIMIT $3 OFFSET $4`, filter.sortColumn(), filter.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []interface{}{content, isCompleted, filter.limit(), filter.offset()}

	rows, err := tm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	tasks := []*Task{}

	for rows.Next() {
		var task Task
		err := rows.Scan(&totalRecords, &task.ID, &task.CreatedAt, &task.Content, &task.IsCompleted)
		if err != nil {
			return nil, Metadata{}, err
		}

		tasks = append(tasks, &task)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filter.Page, filter.PageSize)

	return tasks, metadata, nil
}

func (tm TaskModel) Update(task *Task) error {
	query := `UPDATE task
						SET content = $1, is_completed = $2
						WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := tm.DB.ExecContext(ctx, query, task.Content, task.IsCompleted, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (tm TaskModel) Delete(id int64) error {
	query := `DELETE FROM task
						WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := tm.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type MockTaskModel struct{}

func (t MockTaskModel) Create(task *Task) error {
	return nil
}

func (t MockTaskModel) Read(id int64) (*Task, error) {
	return nil, nil
}

func (t MockTaskModel) ReadAll(content string, IsCompleted bool, filter Filter) ([]*Task, Metadata, error) {
	return nil, Metadata{}, nil
}

func (t MockTaskModel) Update(task *Task) error {
	return nil
}

func (t MockTaskModel) Delete(id int64) error {
	return nil
}
