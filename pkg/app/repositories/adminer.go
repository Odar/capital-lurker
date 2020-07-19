package repositories

import "github.com/Odar/capital-lurker/pkg/app/models"

type AdminerRepo interface {
	PostUniversity() ([]models.University, error)
	//other issues
}
