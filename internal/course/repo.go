package course

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/Odar/capital-lurker/internal/general"
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

func (r *repo) GetCoursesForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) (
	[]api.CourseForAdmin, error) {
	sortBy = general.ApplySortByParameter(sortBy)
	coursesQuery := general.ApplyFilter("course", filter, r.builder.Select(
		"course.id AS course_id, "+
			"course.name AS course_name, "+
			"course.description AS course_description, "+
			"course.position AS course_position, "+
			"course.added_at AS course_added_at, "+
			"course.updated_at AS course_updated_at, "+
			"theme.id AS theme_id, "+
			"theme.name AS theme_name, "+
			"theme.slug AS theme_slug, "+
			"theme.on_main_page AS theme_on_main_page, "+
			"theme.added_at AS theme_added_at, "+
			"theme.updated_at AS theme_updated_at, "+
			"theme.position AS theme_position, "+
			"theme.img AS theme_img").
		From("course")).
		Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).OrderBy("course." + sortBy)
	sql, args, err := coursesQuery.LeftJoin("theme ON course.theme_id = theme.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	unparsedCourses := make([]api.UnparsedCourseForAdmin, 0, limit)
	err = r.postgres.Select(&unparsedCourses, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	courses := make([]api.CourseForAdmin, len(unparsedCourses), len(unparsedCourses))
	for i := range unparsedCourses {
		courses[i].ID = unparsedCourses[i].ID
		courses[i].Name = unparsedCourses[i].Name
		courses[i].Description = unparsedCourses[i].Description
		courses[i].Position = unparsedCourses[i].Position
		courses[i].AddedAt = unparsedCourses[i].AddedAt
		courses[i].UpdatedAt = unparsedCourses[i].UpdatedAt
		if unparsedCourses[i].ThemeID == nil {
			courses[i].Theme = nil
		} else {
			courses[i].Theme = &models.Theme{
				ID:         *unparsedCourses[i].ThemeID,
				Name:       *unparsedCourses[i].ThemeName,
				Slug:       *unparsedCourses[i].ThemeSlug,
				OnMainPage: *unparsedCourses[i].ThemeOnMainPage,
				AddedAt:    *unparsedCourses[i].ThemeAddedAt,
				UpdatedAt:  *unparsedCourses[i].ThemeUpdatedAt,
				Position:   *unparsedCourses[i].ThemePosition,
				Img:        *unparsedCourses[i].ThemeImg,
			}
		}
	}

	return courses, nil
}

func (r *repo) CountCoursesForAdmin(filter *api.Filter) (uint64, error) {
	coursesQuery := general.ApplyFilter("course",
		filter, r.builder.Select("count(*)").From("course"))
	sql, args, err := coursesQuery.ToSql()
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

func (r *repo) DeleteCourse(ID uint64) (int64, error) {
	sql, args, err := r.builder.Delete("*").
		From("course").
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

func (r *repo) UpdateCourseForAdmin(ID uint64, request *api.UpdateCourseForAdminRequest) (*models.Course, error) {
	updateRequest := r.builder.Update("course")
	if request.Name != nil {
		updateRequest = updateRequest.Set("name", *request.Name)
	}
	if request.ThemeID != nil {
		updateRequest = updateRequest.Set("theme_id", *request.ThemeID)
	}
	if request.Description != nil {
		updateRequest = updateRequest.Set("description", *request.Description)
	}
	if request.Position != nil {
		updateRequest = updateRequest.Set("position", *request.Position)
	}
	updateRequest = updateRequest.Set("updated_at", time.Now().UTC())

	sql, args, err := updateRequest.Suffix(
		"RETURNING "+
			"id, "+
			"name, "+
			"theme_id, "+
			"description, "+
			"position, "+
			"added_at, "+
			"updated_at").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	res := &models.Course{}
	err = r.postgres.Get(res, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return res, nil
}

func (r *repo) AddCourseForAdmin(request *api.AddCourseForAdminRequest) (*models.Course, error) {
	sql, args, err := r.builder.Insert("course").
		Columns("name", "theme_id", "description", "position", "added_at", "updated_at").
		Values(request.Name, request.ThemeID, request.Description, request.Position, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING id, name, theme_id, description, position, added_at, updated_at").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	addedCourse := &models.Course{}
	err = r.postgres.Get(addedCourse, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return addedCourse, nil
}
