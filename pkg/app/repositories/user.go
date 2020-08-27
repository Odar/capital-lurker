package repositories

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

type AuthenticatorRepo interface {
	AddUser(email, password, firstName, lastName string, birthDate time.Time, vkID uint64) (*models.User, error)
	CheckAuth(email, password string) (uint64, bool, error)
	CheckRegistration(vkID uint64) (bool, error)
	GetIDForVk(vkID uint64) (uint64, error)
}
