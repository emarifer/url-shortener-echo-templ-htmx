package service

import (
	"context"

	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/model"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/repository"
)

// Service is the business logic of the application.
type Service interface {
	RegisterUser(ctx context.Context, email, username, password string) error
	LoginUser(ctx context.Context, email, password string) (*model.User, error)

	AddLink(ctx context.Context, link *entity.Link) error
	RecoverLinks(ctx context.Context, user_id string) ([]model.Link, error)
	RecoverLink(ctx context.Context, slug string) (*model.Link, error)
	UpdateLink(ctx context.Context, description string, id int) error
	RemoveLink(ctx context.Context, slug, user_id string) error
}

type serv struct {
	repo repository.Repository
}

func New(repo repository.Repository) Service {

	return &serv{
		repo: repo,
	}
}
