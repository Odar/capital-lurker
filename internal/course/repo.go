package course

import (
	"strings"
	"time"

	"github.com/Odar/capital-lurker/pkg/app/models"

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

func (r *repo) GetCoursesForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) (
	[]api.CourseForAdmin, error) {
	sortBy = applySortByParameter(sortBy)
	coursesQuery := applyFilter("course", filter, r.builder.Select(
		"course.id, "+
			"course.name, "+
			"course.description, "+
			"course.added_at, "+
			"course.updated_at, "+
			"theme.id, "+
			"theme.name, "+
			"theme.slug, "+
			"theme.on_main_page, "+
			"theme.added_at, "+
			"theme.updated_at, "+
			"theme.position, "+
			"theme.img").
		From("course")).
		Limit(uint64(limit)).Offset(uint64((page - 1) * limit)) //.OrderBy("course." + sortBy)
	sql, args, err := coursesQuery.LeftJoin("theme ON course.theme_id = theme.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	courses := make([]api.CourseForAdmin, 0, limit)
	stmt, err := r.postgres.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	for stmt.Next() {
		var course api.CourseForAdmin
		var theme api.ThemeForAdmin
		err = stmt.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.AddedAt,
			&course.UpdatedAt,
			&theme.ID,
			&theme.Name,
			&theme.Slug,
			&theme.OnMainPage,
			&theme.AddedAt,
			&theme.UpdatedAt,
			&theme.Position,
			&theme.Img)
		if theme != (api.ThemeForAdmin{}) {
			course.Theme = &theme
		} else {
			course.Theme = nil
			err = nil
		}

		if err != nil {
			return nil, errors.Wrapf(err, "can not extract data from query `%s` with args %+v", sql, args)
		}

		courses = append(courses, course)
	}

	return courses, nil
}

func (r *repo) CountCoursesForAdmin(filter *api.Filter) (uint64, error) {
	coursesQuery := applyFilter("course",
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
	if request.Description != nil {
		updateRequest = updateRequest.Set("description", *request.Description)
	}
	updateRequest = updateRequest.Set("updated_at", time.Now().UTC())

	sql, args, err := updateRequest.Suffix(
		"RETURNING "+
			"id, "+
			"description, "+
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
		Columns("name", "description", "added_at", "updated_at").
		Values(request.Name, request.Description, time.Now().UTC(), time.Now().UTC()).
		Suffix("RETURNING id, name, description, added_at, updated_at").
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
