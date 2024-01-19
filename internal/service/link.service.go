package service

import (
	"context"
	"errors"
	"strings"

	"github.com/emarifer/url-shortener-echo-templ-htmx/encryption"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/entity"
	"github.com/emarifer/url-shortener-echo-templ-htmx/internal/model"
)

var (
	ErrDuplicateUrl     error = errors.New("duplicate url")
	ErrResourceNotFound error = errors.New("resource not found")
)

func (s *serv) AddLink(ctx context.Context, link *entity.Link) error {
	var (
		newSlug string
		err     error
	)

	if len(link.Slug) != 6 {
		newSlug, err = encryption.CreateSlug(6)
		if err != nil {
			return err
		}

		link.Slug = newSlug
	}

	err = s.repo.SaveLink(ctx, link)
	if err != nil {
		if strings.Contains(err.Error(), "violates unique constraint") {
			return ErrDuplicateUrl
		}

		return err
	}

	return nil
}

func (s *serv) RecoverLinks(
	ctx context.Context, user_id string,
) ([]model.Link, error) {
	ll := []model.Link{}

	entityLinks, err := s.repo.GetLinks(ctx, user_id)
	if err != nil {
		return nil, err
	}

	for _, el := range entityLinks {
		l := model.Link{
			ID:          el.ID,
			Url:         el.Url,
			Slug:        el.Slug,
			Description: el.Description,
			CreatedAt:   el.CreatedAt,
		}

		ll = append(ll, l)
	}

	return ll, nil
}

func (s *serv) RecoverLink(
	ctx context.Context, slug string,
) (*model.Link, error) {
	entityLink, err := s.repo.GetLink(ctx, slug)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, ErrResourceNotFound
		}

		return nil, err
	}

	ml := &model.Link{
		ID:          entityLink.ID,
		Url:         entityLink.Url,
		Slug:        entityLink.Slug,
		Description: entityLink.Description,
		CreatedAt:   entityLink.CreatedAt,
	}

	return ml, nil
}

func (s *serv) UpdateLink(
	ctx context.Context, description string, id int,
) error {
	err := s.repo.UpdateLink(ctx, description, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return ErrResourceNotFound
		}

		return err
	}

	return nil
}

func (s *serv) RemoveLink(ctx context.Context, slug, user_id string) error {
	row, err := s.repo.DeleteLink(ctx, slug, user_id)
	if err != nil {
		return err
	}

	if row != 1 {
		return ErrResourceNotFound
	}

	return nil
}
