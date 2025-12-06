package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Infamous003/GoChat/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateUsername = errors.New("duplicate username")
	ErrRecordNotFound    = errors.New("record not found")
)

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  password  `json:"-"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

// Set takes a plain text password and generates a hash using it
func (p *password) Set(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plainPassword
	p.hash = hash

	return nil
}

// Matches compares plain text password with the hash using CompareHashAndPassword
// which generates the hash again and does the comparison.
func (p *password) Matches(plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainPassword))
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

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be atleast 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be longer than 72 bytes")
}

func ValidateUsername(v *validator.Validator, username string) {
	v.Check(username != "", "username", "must be provided")
	v.Check(len(username) >= 8, "username", "must be atleast 8 bytes long")
	v.Check(len(username) <= 32, "username", "must be not be longer than 32 bytes")
}

func ValidateUser(v *validator.Validator, user *User) {
	ValidateUsername(v, user.Username)

	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}
}

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, user.Username, user.Password.hash).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, password_hash, created_at, version FROM users
		WHERE username = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Version,
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
