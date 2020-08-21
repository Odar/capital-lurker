package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
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

func (r *repo) AddUser(email, password, firstName, lastName string, birthDate time.Time) (*models.User, error) {
	res := models.User{}
	sql, args, err := r.builder.Insert("users").Columns("email, password, first_name, last_name, birth_date, signed_up_at, last_signed_in_at, updated_at, img").
		Values(email, password, firstName, lastName, birthDate, time.Now().UTC(), time.Unix(0, 0).UTC(), time.Now().UTC(), "default.png").
		Suffix("RETURNING id, email, password, first_name, last_name, birth_date, signed_up_at, last_signed_in_at, updated_at, img").
		ToSql()

	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	err = r.postgres.Get(&res, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return &res, nil
}
