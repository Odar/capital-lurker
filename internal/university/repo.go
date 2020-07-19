package university

import (
	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/jmoiron/sqlx"
)

func NewRepo(postgres *sqlx.DB) *repo {
	return &repo{
		postgres: postgres,
		builder:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

type repo struct {
	postgres *sqlx.DB
	builder  squirrel.StatementBuilderType
}

func (r *repo) GetUniversities(filter api.Filter) ([]models.University, error) {

}

func (r *repo) GetNumUniversities() (uint64, error) {

}
