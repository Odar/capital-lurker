package repositories

import (
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"
)

const (
	TokenStatusLoggedOut    = "loggedOut"
	TokenStatusOK           = "OK"
	TokenStatusInvalidToken = "invalid token"
	TokenStatusError        = "error"
)

type DbResponseIdPass struct {
	ID   uint64 `db:"id"`
	Pass string `db:"password"`
}

type AuthenticatorRepo interface {
	AddUser(email, password, firstName, lastName string, birthDate time.Time, vkID uint64) (*models.User, error)
	CheckAuth(email, password string) (uint64, bool, error)
	CheckRegistration(vkID uint64) (bool, error)
	GetIDForVk(vkID uint64) (uint64, error)
	UpdateTokenWithID(id uint64, token string) error
	InvalidateToken(id uint64) error
	CheckToken(id uint64, token string) (string, error)
}
