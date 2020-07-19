package models

import "time"

type University struct {
	ID         uint64    `db:"id"`
	Name       string    `db:"name"`
	OnMainPage bool      `db:"on_main_page"`
	InFilter   bool      `db:"in_filter"`
	AddedAt    time.Time `db:"added_at"`
	UpdatedAt  time.Time `db:"updated_at"`
	Position   uint64    `db:"position"`
	Img        string    `db:"img"`
}
