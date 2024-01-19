package repository

import (
	"context"

	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
)

const (
	qryInsertUser = `
		INSERT INTO users (email, username, password)
		VALUES (:email, :username, :password);
	`

	qryGetUserByEmail = `
		SELECT
			user_id,
			email,
			username,
			password
		FROM users
		WHERE email = $1;	
	`
)

func (r *repo) SaveUser(
	ctx context.Context, user *entity.User,
) error {
	_, err := r.db.NamedExecContext(ctx, qryInsertUser, user)

	return err
}

func (r *repo) GetUserByEmail(
	ctx context.Context, email string,
) (*entity.User, error) {
	u := &entity.User{}

	err := r.db.GetContext(ctx, u, qryGetUserByEmail, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}
