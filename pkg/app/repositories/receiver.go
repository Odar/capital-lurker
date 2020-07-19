package repositories

import "github.com/Odar/capital-lurker/pkg/app/models"

type ReceiverRepo interface {
	GetCache(word string) (*models.ReceiverCache, error)
	SetCache(model models.ReceiverCache) error
	AddUsing(model *models.ReceiverCache) error
}
