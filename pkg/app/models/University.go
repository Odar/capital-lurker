package models

import "time"

type University struct {
	ID         uint64    `db:"id" json:"id"`
	Name       string    `db:"name" json:"name"`
	OnMainPage bool      `db:"on_main_page" json:"on_main_page"`
	InFilter   bool      `db:"in_filter" json:"in_filter"`
	AddedAt    time.Time `db:"added_at" json:"added_at"`
	UpdatedAt  time.Time `db:"updated_at" json:"updated_at"`
	Position   uint64    `db:"position" json:"position"`
	Img        string    `db:"img" json:"img"`
}
