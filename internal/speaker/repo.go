package speaker

import (
	"strings"

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

func (r *repo) GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error) {
	if limit <= 0 {
		return nil, errors.New("bad request; parameter: limit should be positive")
	}
	sql, args, err := r.builder.Select("id, name, position, img").
		From("speaker").
		Where("on_main_page = true").
		OrderBy("position DESC").
		Limit(uint64(limit)).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	speakers := make([]api.SpeakerOnMain, 0, limit)
	err = r.postgres.Select(&speakers, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return speakers, nil
}

func (r *repo) GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) (
	[]models.Speaker, error,
) {
	sortBy = validateSortByParameter(sortBy)
	speakersQuery := validateFilterGetSpeakerForAdmin(filter, r.builder.Select("*").From("speaker"))
	sql, args, err := speakersQuery.Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).OrderBy(sortBy).ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	speakers := make([]models.Speaker, 0, limit)
	err = r.postgres.Select(&speakers, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return speakers, nil
}

func (r *repo) CountSpeakersForAdmin(filter *api.Filter) (uint64, error) {
	speakersQuery := validateFilterGetSpeakerForAdmin(filter, r.builder.Select("count(*)").From("speaker"))
	sql, args, err := speakersQuery.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "can not build sql")
	}

	var count uint64
	err = r.postgres.Get(&count, sql, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return count, nil
}

func validateFilterGetSpeakerForAdmin(filter *api.Filter, query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.ID != 0 {
		query = query.Where("id = ?", filter.ID)
	}
	if filter.Name != "" {
		query = query.Where("name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.OnMainPage != nil {
		query = query.Where("on_main_page = ?", *filter.OnMainPage)
	}
	if filter.InFilter != nil {
		query = query.Where("in_filter = ?", *filter.InFilter)
	}
	if !filter.AddedAtRange.From.IsZero() {
		query = query.Where("added_at >= ?", filter.AddedAtRange.From)
	}
	if !filter.AddedAtRange.To.IsZero() {
		query = query.Where("added_at < ?", filter.AddedAtRange.To)
	}
	if !filter.UpdatedAtRange.From.IsZero() {
		query = query.Where("updated_at >= ?", filter.UpdatedAtRange.From)
	}
	if !filter.UpdatedAtRange.To.IsZero() {
		query = query.Where("updated_at < ?", filter.UpdatedAtRange.To)
	}
	if filter.Position != 0 {
		query = query.Where("position = ?", filter.Position)
	}
	if filter.Img != "" {
		query = query.Where("img LIKE ?", "%"+filter.Img+"%")
	}
	return query
}

func validateSortByParameter(sortBy string) string {
	var columnNames = map[string]bool{
		"id":           true,
		"name":         true,
		"on_main_page": true,
		"in_filter":    true,
		"added_at":     true,
		"updated_at":   true,
		"position":     true,
		"img":          true,
	}
	var orderByKeywords = map[string]bool{
		"DESC": true,
		"ASC":  true,
	}
	words := strings.Split(sortBy, " ")
	_, foundColumnName := columnNames[words[0]]
	if len(words) == 1 {
		if !foundColumnName {
			sortBy = "id DESC"
		}
	} else {
		_, foundOrderByKeyword := orderByKeywords[words[1]]
		if !foundColumnName || sortBy == "" || len(words) > 2 || !foundOrderByKeyword {
			sortBy = "id DESC"
		}
	}

	return sortBy
}

func (r *repo) DeleteSpeakerForAdmin(ID uint64) (string, error) {
	sql, args, err := r.builder.Delete("*").
		From("speaker").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return "error", errors.Wrap(err, "can not build sql")
	}

	result, err := r.postgres.Exec(sql, args...)
	if err != nil {
		return "error", errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	count, err := result.RowsAffected()
	if err != nil {
		return "error", errors.Wrapf(err, "can not estimate deleted rows")
	}
	if count == 1 {
		return "deleted", nil
	}
	if count == 0 {
		return "nothing", nil
	}

	return "error", errors.Wrapf(err, "more than one rows were deleted")
}
