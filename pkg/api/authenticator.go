package api

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

type SignUpRequest struct {
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	BirthDate time.Time `json:"birth_date"`
}

type SignUpResponse struct {
	Result models.User `json:"result"`
}
