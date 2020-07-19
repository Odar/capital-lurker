package repositories

import "github.com/Odar/capital-lurker/pkg/app/models"

type AdminerRepo interface {
	GetUniversities() ([]models.University, error)
	//other issues
}
