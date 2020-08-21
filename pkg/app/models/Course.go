package models

import "time"

type Course struct {
	ID          uint64    `db:"id"`
	Name        string    `db:"name"`
	ThemeID     uint64    `db:"theme_id"`
	Description string    `db:"description"`
	Position    uint64    `db:"position"`
	AddedAt     time.Time `db:"added_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
