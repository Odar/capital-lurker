package theme

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/internal/general"
	"github.com/Odar/capital-lurker/pkg/api"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/gosimple/slug"
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

func (r *repo) GetThemesForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) ([]models.Theme, error) {
	base := r.builder.Select("*").From("theme")
	filtered := general.ApplyFilter("theme", filter, base)
	sorted := filtered.OrderBy(general.ApplySortByParameter(sortBy))
	paged := sorted.Limit(uint64(limit)).Offset(uint64((page - 1) * limit))

	sql, args, err := paged.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	themes := make([]models.Theme, 0, limit)
	err = r.postgres.Select(&themes, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return themes, nil
}

func (r *repo) CountThemesForAdmin(filter *api.Filter) (uint64, error) {
	base := r.builder.Select("count(*) as c").From("theme")
	filtered := general.ApplyFilter("theme", filter, base)
	sql, args, err := filtered.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "can not build sql")
	}

	var result uint64
	err = r.postgres.Get(&result, sql, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return result, nil
}

func (r *repo) DeleteTheme(ID uint64) (int64, error) {
	sql, args, err := r.builder.Delete("*").
		From("theme").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return -1, errors.Wrap(err, "can not build sql")
	}

	result, err := r.postgres.Exec(sql, args...)
	if err != nil {
		return -1, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return -1, errors.Wrapf(err, "can not estimate deleted rows")
	}

	return count, nil
}

func (r *repo) UpdateThemeForAdmin(ID uint64, request *api.UpdateThemeForAdminRequest) (*models.Theme, error) {
	updateRequest := r.builder.Update("theme")
	if request.Name != nil {
		updateRequest = updateRequest.Set("name", *request.Name).
			Set("slug", slug.Make(*request.Name))
	}
	if request.OnMainPage != nil {
		updateRequest = updateRequest.Set("on_main_page", *request.OnMainPage)
	}
	updateRequest = updateRequest.Set("updated_at", time.Now().UTC())
	if request.Position != nil {
		updateRequest = updateRequest.Set("position", *request.Position)
	}
	if request.Img != nil {
		updateRequest = updateRequest.Set("img", *request.Img)
	}

	sql, args, err := updateRequest.Suffix(
		"RETURNING "+
			"id, "+
			"name, "+
			"slug, "+
			"on_main_page, "+
			"added_at, "+
			"updated_at, "+
			"position, "+
			"img").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	res := &models.Theme{}
	err = r.postgres.Get(res, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return res, nil
}

func (r *repo) AddThemeForAdmin(request *api.AddThemeForAdminRequest) (*models.Theme, error) {
	sql, args, err := r.builder.Insert("theme").
		Columns("name", "slug", "on_main_page", "added_at", "updated_at", "position", "img").
		Values(request.Name, slug.Make(request.Name), request.OnMainPage, time.Now().UTC(), time.Now().UTC(),
			request.Position, request.Img).
		Suffix("RETURNING id, name, slug, on_main_page, added_at, updated_at, position, img").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	addedTheme := &models.Theme{}
	err = r.postgres.Get(addedTheme, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return addedTheme, nil
}
