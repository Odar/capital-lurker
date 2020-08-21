package repositories

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

type AuthenticatorRepo interface {
	AddUser(email, password, firstName, lastName string, birthDate time.Time) (*models.User, error)
}
