package speaker

import (
    "github.com/Masterminds/squirrel"
    "github.com/Odar/capital-lurker/pkg/api"
    "github.com/Odar/capital-lurker/pkg/app/models"
    "github.com/jmoiron/sqlx"
    "github.com/pkg/errors"
    "strings"
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

func (r *repo) GetSpeakerOnMainFromDB(limit int64) ([]api.SpeakerOnMain, error) {
    if limit <= 0 {
        return nil, errors.New("bad request; parameter: limit should be positive")
    }
    sql, args, err := r.builder.Select("id, name, position, img").
        From("speaker").
        Where("on_main_page = true").
        OrderBy("position DESC").
        ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "can not build sql")
    }

    speakers := make([]api.SpeakerOnMain, 0, limit)
    err = r.postgres.Select(&speakers, sql, args...)
    if err != nil {
        return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
    }

    if len(speakers) > 0 {
        return speakers, nil
    }

    return nil, nil
}

func (r *repo) GetSpeakersForAdminFromDB(limit int64, page int64, sortBy string, filter *api.SpeakerForAdminFilter) (
    []models.Speaker, uint64, error) {
    if limit <= 0 {
        limit = 10
    }
    if page <= 0 {
        page = 1
    }
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

    speakersQuery := validateFilterGetSpeakerForAdmin(filter, r.builder.Select("*").From("speaker"))

    sql, args, err := speakersQuery.Limit(uint64(limit)).OrderBy(sortBy).ToSql()

    if err != nil {
        return nil, 0, errors.Wrap(err, "can not build sql")
    }

    speakers := make([]models.Speaker, 0, limit)
    err = r.postgres.Select(&speakers, sql, args...)
    if err != nil {
        return nil, 0, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
    }

    if len(speakers) > 0 {
        return speakers, uint64(len(speakers)), nil
    }

    return nil, 0, nil
}

func validateFilterGetSpeakerForAdmin(filter *api.SpeakerForAdminFilter, query squirrel.SelectBuilder) squirrel.SelectBuilder {
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
