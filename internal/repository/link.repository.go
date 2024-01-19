package repository

import (
	"context"

	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
)

const (
	qryInsertLink = `
		INSERT INTO links (url, slug, description, user_id)
		VALUES (:url, :slug, :description, :user_id);
	`

	qryGetLinkBySlug = `
		SELECT * FROM links
		WHERE slug = $1;
	`

	qryGetAllLinks = `
		SELECT * FROM links
		WHERE user_id = $1
		ORDER BY created_at DESC;
	`

	qryUpdateLink = `
		UPDATE links
		SET description = $1
		WHERE id = $2;
	`

	qryDeleteLink = `
		DELETE FROM links
		WHERE slug = $1 AND user_id = $2;
	`
)

func (r *repo) SaveLink(ctx context.Context, link *entity.Link) error {
	_, err := r.db.NamedExecContext(ctx, qryInsertLink, link)

	return err
}

func (r *repo) GetLink(
	ctx context.Context, slug string,
) (*entity.Link, error) {
	l := &entity.Link{}

	err := r.db.GetContext(ctx, l, qryGetLinkBySlug, slug)
	if err != nil {
		return nil, err
	}

	return l, nil
}

func (r *repo) GetLinks(
	ctx context.Context, user_id string,
) ([]entity.Link, error) {
	ll := []entity.Link{}

	err := r.db.SelectContext(ctx, &ll, qryGetAllLinks, user_id)
	if err != nil {
		return nil, err
	}

	return ll, nil
}

func (r *repo) UpdateLink(
	ctx context.Context, description string, id int,
) error {

	_, err := r.db.ExecContext(ctx, qryUpdateLink, description, id)

	return err
}

func (r *repo) DeleteLink(
	ctx context.Context, slug, user_id string,
) (int64, error) {
	var row int64

	result, err := r.db.ExecContext(ctx, qryDeleteLink, slug, user_id)
	if err != nil {
		return 0, err
	}

	if row, err = result.RowsAffected(); err != nil {
		return 0, err
	}

	return row, nil
}
