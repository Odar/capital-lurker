package university

import (
	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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

func (r *repo) GetUniversities(filter *api.Filter, sortBy string, limit, page int) (*repositories.GetUniversitiesRepoResponse, error) {
	base := r.builder.Select("*").From("university")
	filtered := applyFilter(base, filter)

	sorted := filtered
	switch sortBy {
	case "":
		sorted = sorted.OrderBy("id DESC")
	default:
		sorted = sorted.OrderBy(sortBy)
	}

	paged := sorted
	paged = paged.Limit(uint64(limit))
	paged = paged.Offset(uint64((page - 1) * limit))

	sql, args, err := paged.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}
	content := make([]models.University, 0)
	err = r.postgres.Select(&content, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	getCount := applyFilter(r.builder.Select("count(*)").From("university"), filter)
	sql, args, err = getCount.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}
	var count []uint64
	err = r.postgres.Select(&count, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	if len(content) > 0 {
		return &repositories.GetUniversitiesRepoResponse{
			Universities: content,
			Count:        count[0],
		}, nil
	}
	return nil, nil
}

func applyFilter(base squirrel.SelectBuilder, filter *api.Filter) squirrel.SelectBuilder {
	filtered := base
	if filter != nil {
		if filter.Id != 0 { //add fix to query: id starts from 1
			filtered = filtered.Where("id = ?", filter.Id)
		}
		if filter.Name != "" {
			filtered = filtered.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.OnMainPage != nil { //how to parse blanks?
			filtered = filtered.Where("on_main_page = ?", *filter.OnMainPage)
		}
		if filter.InFilter != nil {
			filtered = filtered.Where("in_filter = ?", *filter.InFilter)
		}
		if filter.AddedAtRange != nil {
			filtered = filtered.Where("added_at >= ? AND added_at < ?", filter.AddedAtRange.From, filter.AddedAtRange.To)
		}
		if filter.UpdatedAtRange != nil {
			filtered = filtered.Where("updated_at >= ? AND updated_at < ?", filter.UpdatedAtRange.From, filter.UpdatedAtRange.To)
		}
		if filter.Position != 0 {
			filtered = filtered.Where("position = ?", filter.Position)
		}
		if filter.Img != "" {
			filtered = filtered.Where("img LIKE ?", "%"+filter.Img+"%")
		}
	}
	return filtered
}
