package speaker

import (
    "github.com/Masterminds/squirrel"
    "github.com/Odar/capital-lurker/pkg/api"
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
