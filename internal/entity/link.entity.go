package entity

import "time"

type Link struct {
	ID          int       `db:"id"`
	Url         string    `db:"url"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UserID      string    `db:"user_id"`
}
