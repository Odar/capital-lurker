package models

import "time"

type User struct {
	ID             uint64    `db:"id" json:"id"`
	Email          string    `db:"email" json:"email"`
	Password       string    `db:"password" json:"password"`
	FirstName      string    `db:"first_name" json:"first_name"`
	LastName       string    `db:"last_name" json:"last_name"`
	BirthDate      time.Time `db:"birth_date" json:"birth_date"`
	SignedUpAt     time.Time `db:"signed_up_at" json:"signed_up_at"`
	LastSignedInAt time.Time `db:"last_signed_in_at" json:"last_signed_in_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
	Img            string    `db:"img" json:"img"`
}
