package repository

import (
	"context"

	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
	"github.com/jmoiron/sqlx"
)

// Repository is the interface that wraps the basic CRUD operations.
type Repository interface {
	SaveUser(ctx context.Context, user *entity.User) error
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	SaveLink(ctx context.Context, link *entity.Link) error
	GetLink(ctx context.Context, slug string) (*entity.Link, error)
	GetLinks(ctx context.Context, user_id string) ([]entity.Link, error)
	UpdateLink(ctx context.Context, description string, id int) error
	DeleteLink(ctx context.Context, slug, user_id string) (int64, error)
}

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {

	return &repo{
		db: db,
	}
}
