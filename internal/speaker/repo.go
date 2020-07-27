package speaker

import (
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

func (r *repo) GetSpeakerForAdminFromDB(decodedParams *api.GetSpeakerForAdminRequest) ([]models.Speaker, uint64, error) {
    if decodedParams.Limit <= 0 {
        decodedParams.Limit = 10
    }
    if decodedParams.Page <= 0 {
        decodedParams.Page = 1
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
    _, found := columnNames[decodedParams.SortBy]
    if !found || decodedParams.SortBy == "" {
        decodedParams.SortBy = "id DESC"
    }

    filterId := decodedParams.Filter.ID
    filterName := decodedParams.Filter.Name
    filterOnMainPage := decodedParams.Filter.OnMainPage
    filterInFilter := decodedParams.Filter.InFilter
    filterAddedAtFrom := decodedParams.Filter.AddedAtRange.From
    filterAddedAtTo := decodedParams.Filter.AddedAtRange.To
    filterUpdatedAtFrom := decodedParams.Filter.UpdatedAtRange.From
    filterUpdateAtTo := decodedParams.Filter.UpdatedAtRange.To
    filterPosition := decodedParams.Filter.Position
    filterImg := decodedParams.Filter.Img

    speakersQuery := r.builder.Select("*").From("speaker")
    if filterId != 0 {
        speakersQuery = speakersQuery.Where("id = ?", filterId)
    }
    if filterName != "" {
        speakersQuery = speakersQuery.Where("name LIKE ?", "%"+filterName+"%")
    }
    if filterOnMainPage != nil {
        speakersQuery = speakersQuery.Where("on_main_page = ?", *filterOnMainPage)
    }
    if filterInFilter != nil {
        speakersQuery = speakersQuery.Where("in_filter = ?", *filterInFilter)
    }
    if !filterAddedAtFrom.IsZero() {
        speakersQuery = speakersQuery.Where("added_at >= ?", filterAddedAtFrom)
    }
    if !filterAddedAtTo.IsZero() {
        speakersQuery = speakersQuery.Where("added_at < ?", filterAddedAtTo)
    }
    if !filterUpdatedAtFrom.IsZero() {
        speakersQuery = speakersQuery.Where("updated_at >= ?", filterUpdatedAtFrom)
    }
    if !filterUpdateAtTo.IsZero() {
        speakersQuery = speakersQuery.Where("updated_at < ?", filterUpdateAtTo)
    }
    if filterPosition != 0 {
        speakersQuery = speakersQuery.Where("position = ?", filterPosition)
    }
    if filterImg != "" {
        speakersQuery = speakersQuery.Where("img LIKE ?", "%"+filterImg+"%")
    }

    sql, args, err := speakersQuery.ToSql()

    if err != nil {
        return nil, 0, errors.Wrap(err, "can not build sql")
    }

    speakers := make([]models.Speaker, 0, decodedParams.Limit)
    err = r.postgres.Select(&speakers, sql, args...)
    if err != nil {
        return nil, 0, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
    }

    if len(speakers) > 0 {
        if int64(len(speakers)) > decodedParams.Limit {
            return speakers[:decodedParams.Limit], uint64(decodedParams.Limit), nil
        } else {
            return speakers, uint64(len(speakers)), nil
        }
    }

    return nil, 0, nil
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
