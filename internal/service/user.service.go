package service

import (
	"context"
	"errors"
	"strings"

	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/model"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  error = errors.New("user already exists")
	ErrInvalidCredentials error = errors.New("invalid credentials")
)

func (s *serv) RegisterUser(
	ctx context.Context, email, username, password string,
) error {
	u, _ := s.repo.GetUserByEmail(ctx, email)
	if u != nil {
		return ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		return err
	}

	newUser := &entity.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}

	return s.repo.SaveUser(ctx, newUser)
}

func (s *serv) LoginUser(
	ctx context.Context, email, password string,
) (*model.User, error) {
	u, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {

			return nil, ErrInvalidCredentials
		}

		return nil, err
	}

	// compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &model.User{
		UserID:   u.UserID,
		Email:    u.Email,
		Username: u.Username,
	}, nil
}
