package service

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/Mark-Grigorev/FinGo/internal/domain"
	"github.com/Mark-Grigorev/FinGo/internal/repository"
	"github.com/Mark-Grigorev/FinGo/pkg/token"
)

type AuthService struct {
	store     *repository.Store
	tokenMaker *token.Maker
	log       *slog.Logger
}

func NewAuth(store *repository.Store, maker *token.Maker, log *slog.Logger) *AuthService {
	return &AuthService{store: store, tokenMaker: maker, log: log}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*domain.User, string, *token.Payload, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, "", nil, domain.ErrInvalidInput
	}

	user, err := s.store.GetUserByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, "", nil, domain.ErrUnauthorized
	}
	if err != nil {
		return nil, "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", nil, domain.ErrUnauthorized
	}

	tokenStr, payload, err := s.tokenMaker.CreateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", nil, err
	}
	return user, tokenStr, payload, nil
}

func (s *AuthService) Register(ctx context.Context, email, name, password string) (*domain.User, string, *token.Payload, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	name = strings.TrimSpace(name)
	if email == "" || name == "" || len(password) < 6 {
		return nil, "", nil, domain.ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", nil, err
	}

	user, err := s.store.CreateUser(ctx, email, name, string(hash))
	if err != nil {
		return nil, "", nil, domain.ErrAlreadyExists
	}

	tokenStr, payload, err := s.tokenMaker.CreateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", nil, err
	}
	return user, tokenStr, payload, nil
}

func (s *AuthService) VerifyToken(tokenStr string) (*token.Payload, error) {
	return s.tokenMaker.VerifyToken(tokenStr)
}

func (s *AuthService) GetUser(ctx context.Context, userID int64) (*domain.User, error) {
	return s.store.GetUserByID(ctx, userID)
}
