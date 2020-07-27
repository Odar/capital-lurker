package university

import (
	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
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

func MakeFilters(base squirrel.SelectBuilder, filter *api.Filter) squirrel.SelectBuilder {
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

func (r *repo) GetUniversities(request api.PostRequest) ([]models.University, error) {
	filter := request.Filter
	base := r.builder.Select("*").From("university")
	filtered := MakeFilters(base, filter)

	sorted := filtered
	switch request.SortBy {
	case "":
		sorted = sorted.OrderBy("id DESC")
	default:
		sorted = sorted.OrderBy(request.SortBy)
	}

	paged := sorted
	paged = paged.Limit(uint64(request.Limit))
	paged = paged.Offset(uint64((request.Page - 1) * request.Limit))

	sql, args, err := paged.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}
	size := 0 // size of content
	content := make([]models.University, size)
	err = r.postgres.Select(&content, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}
	if len(content) > 0 {
		return content, nil
	}
	return nil, nil
}
