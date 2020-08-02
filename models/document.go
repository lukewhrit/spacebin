package models

// Document is the structure of a document in the database
type Document struct {
	ID        string `db:"id"`
	Content   string `db:"content"`
	Extension string `db:"extension"`
	CreatedAt int    `db:"created_at"`
	UpdatedAt int    `db:"updated_at"`
}
