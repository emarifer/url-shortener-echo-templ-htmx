package model

import "time"

type Link struct {
	ID          int       `json:"id"`
	Url         string    `json:"url"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}
