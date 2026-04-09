package storage

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vadim/filestorage/internal/model"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user already exists")

type UserStorage struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserStorage {
	return &UserStorage{pool: pool}
}

func (s *UserStorage) Create(ctx context.Context, username, password string) (model.User, error) {
	var user model.User
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users (username, password) VALUES ($1, $2)
		 RETURNING id, username, password, created_at`,
		username, password,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return model.User{}, ErrUserExists
	}
	return user, err
}

func (s *UserStorage) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := s.pool.QueryRow(ctx,
		`SELECT id, username, password, created_at FROM users WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.CreatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return model.User{}, ErrUserNotFound
	}
	return user, err
}

