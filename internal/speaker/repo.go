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
	sql, args, err := r.builder.Select("speaker.id, speaker.name, speaker.position, speaker.img, " +
		"university.id, university.name, university.img").
		From("speaker").
		Where("speaker.on_main_page = true").
		OrderBy("speaker.position DESC").
		Limit(uint64(limit)).
		LeftJoin("university ON speaker.university_id = university.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	speakers := make([]api.SpeakerOnMain, 0, limit)
	stmt, err := r.postgres.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	for stmt.Next() {
		var speaker api.SpeakerOnMain
		var university api.UniversityOnMain
		err = stmt.Scan(
			&speaker.ID,
			&speaker.Name,
			&speaker.Position,
			&speaker.Img,
			&university.ID,
			&university.Name,
			&university.Img)
		if university != (api.UniversityOnMain{}) {
			speaker.University = &university
		} else {
			speaker.University = nil
			err = nil
		}

		if err != nil {
			return nil, errors.Wrapf(err, "can not extract data from query `%s` with args %+v", sql, args)
		}

		speakers = append(speakers, speaker)
	}

	return speakers, nil
}

func (r *repo) GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) (
	[]api.SpeakerForAdmin, error) {
	sortBy = applySortByParameter(sortBy)
	speakersQuery := applyFilter("speaker", filter, r.builder.Select(
		"speaker.id, "+
			"speaker.name, "+
			"speaker.on_main_page, "+
			"speaker.in_filter, "+
			"speaker.added_at, "+
			"speaker.updated_at, "+
			"speaker.position, "+
			"speaker.img, "+
			"university.id, "+
			"university.name, "+
			"university.on_main_page, "+
			"university.in_filter, "+
			"university.added_at, "+
			"university.updated_at, "+
			"university.position, "+
			"university.img").
		From("speaker")).
		Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).OrderBy("speaker." + sortBy)
	sql, args, err := speakersQuery.LeftJoin("university ON speaker.university_id = university.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	speakers := make([]api.SpeakerForAdmin, 0, limit)
	stmt, err := r.postgres.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	for stmt.Next() {
		var speaker api.SpeakerForAdmin
		var university api.UniversityForAdmin
		err = stmt.Scan(
			&speaker.ID,
			&speaker.Name,
			&speaker.OnMainPage,
			&speaker.InFilter,
			&speaker.AddedAt,
			&speaker.UpdatedAt,
			&speaker.Position,
			&speaker.Img,
			&university.ID,
			&university.Name,
			&university.OnMainPage,
			&university.InFilter,
			&university.AddedAt,
			&university.UpdatedAt,
			&university.Position,
			&university.Img)
		if university != (api.UniversityForAdmin{}) {
			speaker.University = &university
		} else {
			speaker.University = nil
			err = nil
		}

		if err != nil {
			return nil, errors.Wrapf(err, "can not extract data from query `%s` with args %+v", sql, args)
		}

		speakers = append(speakers, speaker)
	}

	return speakers, nil
}

func (r *repo) CountSpeakersForAdmin(filter *api.Filter) (uint64, error) {
	speakersQuery := applyFilter("speaker",
		filter, r.builder.Select("count(*)").From("speaker"))
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

func (r *repo) DeleteSpeaker(ID uint64) (int64, error) {
	sql, args, err := r.builder.Delete("*").
		From("speaker").
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

func (r *repo) UpdateSpeakerForAdmin(ID uint64, request *api.UpdateSpeakerForAdminRequest) (
	*models.Speaker, error) {
	updateRequest := r.builder.Update("speaker")
	if request.Name != nil {
		updateRequest = updateRequest.Set("name", *request.Name)
	}
	if request.OnMainPage != nil {
		updateRequest = updateRequest.Set("on_main_page", *request.OnMainPage)
	}
	if request.InFilter != nil {
		updateRequest = updateRequest.Set("in_filter", *request.InFilter)
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
			"on_main_page, "+
			"in_filter, "+
			"added_at, "+
			"updated_at, "+
			"position, "+
			"img").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	res := &models.Speaker{}
	err = r.postgres.Get(res, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return res, nil
}

func (r *repo) AddSpeakerForAdmin(request *api.AddSpeakerForAdminRequest) (*models.Speaker, error) {
	sql, args, err := r.builder.Insert("speaker").
		Columns("name", "on_main_page", "in_filter", "added_at", "updated_at", "position", "img").
		Values(request.Name, request.OnMainPage, request.InFilter, time.Now().UTC(), time.Now().UTC(),
			request.Position, request.Img).
		Suffix("RETURNING id, name, on_main_page, in_filter, added_at, updated_at, position, img").
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

func applyFilter(tableName string, filter *api.Filter, query squirrel.SelectBuilder) squirrel.SelectBuilder {
	if filter != nil {
		if filter.ID != nil { //add fix to query: id starts from 1
			query = query.Where(tableName+".id = ?", *filter.ID)
		}
		if filter.Name != nil {
			query = query.Where(tableName+".name LIKE ?", "%"+*filter.Name+"%")
		}
		if filter.OnMainPage != nil { //how to parse blanks?
			query = query.Where(tableName+".on_main_page = ?", *filter.OnMainPage)
		}
		if filter.InFilter != nil {
			query = query.Where(tableName+".in_filter = ?", *filter.InFilter)
		}
		if filter.AddedAtRange != nil {
			query = query.Where(tableName+".added_at >= ? AND "+
				tableName+".added_at < ?", filter.AddedAtRange.From, filter.AddedAtRange.To)
		}
		if filter.UpdatedAtRange != nil {
			query = query.Where(tableName+".updated_at >= ? AND "+
				tableName+".updated_at < ?", filter.UpdatedAtRange.From, filter.UpdatedAtRange.To)
		}
		if filter.Position != nil {
			query = query.Where(tableName+".position = ?", *filter.Position)
		}
		if filter.Img != nil {
			query = query.Where(tableName+".img LIKE ?", "%"+*filter.Img+"%")
		}
	}
	return query
}

func applySortByParameter(sortBy string) string {
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
