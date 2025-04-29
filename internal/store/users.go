package store

import (
	"context"
	"database/sql"
	"errors"
)

type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}
type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username, email,password) 
	VALUES ($1,$2,$3) RETURNING id,created_at
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserStore) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	query := `
		SELECT id, email, username, created_at
		FROM users
		WHERE id = $1
	`
	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.Id,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
