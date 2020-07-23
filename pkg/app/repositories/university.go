package repositories

import (
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
)

type AdminerRepo interface {
	GetUniversities(filter *api.Filter, sortBy string) ([]models.University, error)
	AddUniversity(request api.PutRequest) (*models.University, error)
	DeleteUniversity(id uint64) (*string, error)
	UpdateUniversity(request api.PostIdRequest, id uint64) (*models.University, error)
	//other issues
}
