package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"hukaapp/todo/internal/validator"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreateAt  time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
}

type password struct {
	plaintext *string
	hash      []byte
}

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserModel struct {
	DB *sql.DB
}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (m UserModel) Create(user *User) error {
	query := `INSERT INTO app_user (name, email, password_hash, activated)
						VALUES ($1, $2, $3, $4)
						RETURNING id, created_at`

	args := []interface{}{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreateAt)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constaint "user_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) ReadByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated
						FROM app_user
						WHERE email = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.CreateAt, &user.Name, &user.Email, &user.Password.hash, &user.Activated)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (m UserModel) Update(user *User) error {
	query := `UPDATE app_user
						SET name = $1, email = $2, password_hash = $3, activated = $4
						WHERE id = $5`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constaint "user_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
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

func (m UserModel) GetForToken(tokenScope, tokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `SELECT app_user.id, app_user.created_at, app_user.name, app_user.email, app_user.password_hash, app_user.activated
	FROM app_user
	INNER JOIN token ON app_user.id = token.user_id
	WHERE token.hash = $1 AND token.scope = $2 AND token.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreateAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
