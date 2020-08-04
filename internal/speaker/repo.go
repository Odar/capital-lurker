package speaker

import (
	"strings"
	"time"

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
	nestedSelectUniversity := r.builder.Select("id, name, img").
		From("university")
	sql, args, err := r.builder.Select("speaker.id, speaker.name, speaker.position, speaker.img").
		From("speaker").
		Where("on_main_page = true").
		OrderBy("position DESC").
		Limit(uint64(limit)).
		JoinClause(nestedSelectUniversity.Prefix("LEFT JOIN (").Suffix(") u ON speaker.university_id = u.id")).
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

func (r *repo) GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.SpeakerForAdminFilter) (
	[]api.SpeakerForAdmin, error) {
	sortBy = validateSortByParameter(sortBy)
	speakersQuery := validateFilterGetSpeakerForAdmin(filter, r.builder.Select(
		"speaker.id, "+
			"speaker.name, "+
			"speaker.on_main_page, "+
			"speaker.in_filter, "+
			"speaker.added_at, "+
			"speaker.updated_at, "+
			"speaker.position, "+
			"speaker.img").From("speaker")).
		Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).OrderBy("speaker." + sortBy)
	sql, args, err := speakersQuery.LeftJoin("university ON speaker.university_id = university.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	speakers := make([]api.SpeakerForAdmin, 0, limit)
	err = r.postgres.Select(&speakers, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return speakers, nil
}

func (r *repo) CountSpeakersForAdmin(filter *api.SpeakerForAdminFilter) (uint64, error) {
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

func validateFilterGetSpeakerForAdmin(filter *api.SpeakerForAdminFilter, query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter.ID != 0 {
		query = query.Where("speaker.id = ?", filter.ID)
	}
	if filter.Name != "" {
		query = query.Where("speaker.name LIKE ?", "%"+filter.Name+"%")
	}
	if filter.OnMainPage != nil {
		query = query.Where("speaker.on_main_page = ?", *filter.OnMainPage)
	}
	if filter.InFilter != nil {
		query = query.Where("speaker.in_filter = ?", *filter.InFilter)
	}
	if !filter.AddedAtRange.From.IsZero() {
		query = query.Where("speaker.added_at >= ?", filter.AddedAtRange.From)
	}
	if !filter.AddedAtRange.To.IsZero() {
		query = query.Where("speaker.added_at < ?", filter.AddedAtRange.To)
	}
	if !filter.UpdatedAtRange.From.IsZero() {
		query = query.Where("speaker.updated_at >= ?", filter.UpdatedAtRange.From)
	}
	if !filter.UpdatedAtRange.To.IsZero() {
		query = query.Where("speaker.updated_at < ?", filter.UpdatedAtRange.To)
	}
	if filter.Position != 0 {
		query = query.Where("speaker.position = ?", filter.Position)
	}
	if filter.Img != "" {
		query = query.Where("speaker.img LIKE ?", "%"+filter.Img+"%")
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
func (r *repo) DeleteSpeakerForAdminFromDB(ID uint64) (string, error) {
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

func (r *repo) UpdateSpeakerForAdminInDB(ID uint64, request *api.UpdateSpeakerForAdminRequest) (
	*models.Speaker, error) {
	sql, args, err := r.builder.Update("speaker").
		Set("name", request.Name).
		Set("on_main_page", request.OnMainPage).
		Set("in_filter", request.InFilter).
		Set("updated_at", time.Now().UTC()).
		Set("position", request.Position).
		Set("img", request.Img).
		Suffix("RETURNING *").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	updateSpeaker := &models.Speaker{}
	err = r.postgres.Get(updateSpeaker, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return updateSpeaker, nil
}

func (r *repo) AddSpeakerForAdminInDB(request *api.AddSpeakerForAdminRequest) (*models.Speaker, error) {
	sql, args, err := r.builder.Insert("speaker").
		Columns("name", "on_main_page", "in_filter", "added_at", "updated_at", "position", "img").
		Values(request.Name, request.OnMainPage, request.InFilter, time.Now().UTC(), time.Now().UTC(),
			request.Position, request.Img).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	addedSpeaker := &models.Speaker{}
	err = r.postgres.Get(addedSpeaker, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return addedSpeaker, nil
}
