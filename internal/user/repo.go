package user

import (
	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/pkg/app/models"
	"github.com/Odar/capital-lurker/pkg/app/repositories"
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

func (r *repo) AddUser(email, password, firstName, lastName string, birthDate time.Time, vkID uint64) (*models.User, error) {
	res := models.User{}
	sql, args, err := r.builder.Insert("users").Columns("vk_id, email, password, first_name, last_name, birth_date, signed_up_at, last_signed_in_at, updated_at, img, token").
		Values(vkID, email, password, firstName, lastName, birthDate, time.Now().UTC(), time.Unix(0, 0).UTC(), time.Now().UTC(), "default.png", "").
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

func (r *repo) CheckAuth(email, password string) (uint64, bool, error) {
	res := repositories.DbResponseIdPass{}
	sql, args, err := r.builder.Select("id, password").From("users").Where("email = ?", email).ToSql()

	if err != nil {
		return 0, false, errors.Wrap(err, "can not build sql")
	}

	err = r.postgres.Get(&res, sql, args...)
	if err != nil {
		return 0, false, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	if res.Pass != password {
		return 0, false, errors.New("Wrong email/password")
	}
	return res.ID, true, nil
}

func (r *repo) CheckRegistration(vkID uint64) (bool, error) {
	sql, args, err := r.builder.Select("COUNT (*) as c").From("users").Where("vk_id = ?", vkID).ToSql()
	if err != nil {
		return false, errors.Wrap(err, "can not build sql")
	}

	var result uint64
	err = r.postgres.Get(&result, sql, args...)
	if err != nil {
		return false, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	if result > 0 {
		return true, nil
	}
	return false, nil
}

func (r *repo) GetIDForVk(vkID uint64) (uint64, error) {
	sql, args, err := r.builder.Select("id").From("users").Where("vk_id = ?", vkID).ToSql()
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

func (r *repo) UpdateTokenWithID(id uint64, token string) error {
	sql, args, err := r.builder.Update("users").Set("token", token).Where("id = ?", id).ToSql()
	if err != nil {
		return errors.Wrap(err, "can not build sql")
	}

	res, err := r.postgres.Exec(sql, args...)

	if numRowsInserted, _ := res.RowsAffected(); numRowsInserted < 1 {
		return errors.Wrap(err, "can not update token")
	}

	return nil
}

func (r *repo) InvalidateToken(id uint64) error {
	return r.UpdateTokenWithID(id, "")
}

func (r *repo) CheckToken(id uint64, token string) (string, error) {
	sql, args, err := r.builder.Select("token").From("users").Where("id = ?", id).ToSql()
	if err != nil {
		return repositories.TokenStatusError, errors.Wrap(err, "can not build sql")
	}

	var result string
	err = r.postgres.Get(&result, sql, args...)
	if err != nil {
		return repositories.TokenStatusError, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	if result == "" {
		return repositories.TokenStatusLoggedOut, nil
	}
	if result != token {
		return repositories.TokenStatusInvalidToken, nil
	}
	return repositories.TokenStatusOK, nil

}
