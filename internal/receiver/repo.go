package receiver

import (
	"github.com/Masterminds/squirrel"
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

func (r *repo) GetCache(word string) (*models.ReceiverCache, error) {
	sql, args, err := r.builder.Select("id, word, answer, uses").
		From("receiver_cache").
		Where("word = ?", word).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	cache := make([]models.ReceiverCache, 0)
	err = r.postgres.Select(&cache, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	if len(cache) > 0 {
		return &cache[0], nil
	}

	return nil, nil
}

func (r *repo) SetCache(model models.ReceiverCache) error {
	sql, args, err := r.builder.Insert("receiver_cache").
		Columns("word, answer, uses").
		Values(model.Word, model.Answer, model.Uses).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "can not build sql")
	}

	_, err = r.postgres.Exec(sql, args...)
	if err != nil {
		return errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return nil
}

func (r *repo) AddUsing(model *models.ReceiverCache) error {
	sql, args, err := r.builder.Update("receiver_cache").
		Set("uses", model.Uses+1).
		Where("id = ?", model.ID).
		ToSql()
	if err != nil {
		return errors.Wrap(err, "can not build sql")
	}

	_, err = r.postgres.Exec(sql, args...)
	if err != nil {
		return errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	model.Uses++

	return nil
}
