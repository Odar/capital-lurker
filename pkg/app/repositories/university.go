package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type AdminerRepo interface {
	GetUniversities(filter api.Filter) ([]models.University, error)
	GetNumUniversities() (uint64, error)
	//other issues
}
