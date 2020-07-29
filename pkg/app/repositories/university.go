package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type AdminerRepo interface {
	GetUniversitiesList(filter *api.Filter, sortBy string, limit, page int) ([]models.University, error)
	CountUniversities(filter *api.Filter) (uint64, error)
}
