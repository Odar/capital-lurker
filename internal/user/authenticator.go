package user

import "github.com/Odar/capital-lurker/pkg/app/repositories"

func New(repo repositories.AuthenticatorRepo) *authenticator {
	return &authenticator{
		repo: repo,
	}
}

type authenticator struct {
	repo repositories.AuthenticatorRepo
}
