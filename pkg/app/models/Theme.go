package models

import "time"

type Theme struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	Slug       string    `db:"slug"`
	OnMainPage bool      `db:"on_main_page"`
	AddedAt    time.Time `db:"added_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Position   uint64    `db:"position"`
	Img        string    `db:"img"`
}
