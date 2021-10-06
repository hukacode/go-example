package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Model struct {
	Task interface {
		Create(task *Task) error
		Read(id int64) (*Task, error)
		ReadAll(content string, IsCompleted bool, filter Filter) ([]*Task, Metadata, error)
		Update(task *Task) error
		Delete(id int64) error
	}
	User interface {
		Create(user *User) error
		ReadByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope, tokenPlaintext string) (*User, error)
	}
	Token interface {
		New(userID int64, ttl time.Duration, scope string) (*Token, error)
		Create(token *Token) error
		DeleteAllForUser(scope string, userID int64) error
	}
	Permission interface {
		ReadAllForUser(userID int64) (Permissions, error)
		AddForUser(userID int64, codes ...string) error
	}
}

func NewModel(db *sql.DB) Model {
	return Model{
		Task:       TaskModel{DB: db},
		User:       UserModel{DB: db},
		Token:      TokenModel{DB: db},
		Permission: PermissionModel{DB: db},
	}
}

func NewMockModel() Model {
	return Model{
		Task: MockTaskModel{},
	}
}
