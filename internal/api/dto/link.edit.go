package dto

type LinkEdit struct {
	Slug        string `query:"slug"`
	Description string `query:"description"`
	ID          int    `query:"id" header:"X-Id"`
}
