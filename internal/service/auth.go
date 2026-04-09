package service

import (
	"context"
	"errors"

	"github.com/vadim/filestorage/internal/auth"
	"github.com/vadim/filestorage/internal/model"
	"github.com/vadim/filestorage/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var ErrWrongCredentials = errors.New("wrong username or password")

type AuthService struct {
	users     *storage.UserStorage
	jwtSecret string
}

func NewAuthService(users *storage.UserStorage, jwtSecret string) *AuthService {
	return &AuthService{users: users, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}
	return s.users.Create(ctx, req.Username, string(hashed))
}

func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (string, error) {
	user, err := s.users.GetByUsername(ctx, req.Username)
	if err != nil {
		if err == storage.ErrUserNotFound {
			return "", ErrWrongCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", ErrWrongCredentials
	}

	return auth.GenerateToken(user.ID, s.jwtSecret)
}
