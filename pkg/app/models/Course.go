package models

import "time"

type Course struct {
	ID          uint64    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	AddedAt     time.Time `db:"added_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
