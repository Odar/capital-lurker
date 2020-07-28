package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type GetUniversitiesRepoResponse struct {
	Universities []models.University
	Count        uint64
}

type AdminerRepo interface {
	GetUniversities(filter *api.Filter, sortBy string, limit, page int) (*GetUniversitiesRepoResponse, error)
}
