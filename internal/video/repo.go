package video

import (
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

func (r *repo) GetVideosForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) (
	[]api.VideoForAdmin, error) {
	sortBy = general.ApplySortByParameter(sortBy)
	videosQuery := general.ApplyFilter("video", filter, r.builder.Select(
		"video.id AS video_id, "+
			"video.name AS video_name, "+
			"video.img AS video_img, "+
			"video.video AS video_video, "+
			"video.youtube_video AS video_youtube_video, "+
			"video.position_in_course AS video_position_in_course, "+
			"video.added_at AS video_added_at, "+
			"video.updated_at AS video_updated_at, "+
			"video.uploaded_at AS video_uploaded_at, "+
			"video.youtubed_at AS video_youtubed_at, "+
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
		From("video")).
		Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).OrderBy("video." + sortBy)
	sql, args, err := videosQuery.LeftJoin("course ON video.course_id = course.id").
		LeftJoin("theme ON course.theme_id = theme.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	unparsedVideos := make([]api.UnparsedVideoForAdmin, 0, limit)
	err = r.postgres.Select(&unparsedVideos, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	videos := make([]api.VideoForAdmin, len(unparsedVideos), len(unparsedVideos))
	for i := range unparsedVideos {
		videos[i].ID = unparsedVideos[i].ID
		videos[i].Name = unparsedVideos[i].Name
		videos[i].Img = unparsedVideos[i].Img
		videos[i].Video = unparsedVideos[i].Video
		videos[i].YouTubeVideo = unparsedVideos[i].YouTubeVideo
		videos[i].PositionInCourse = unparsedVideos[i].PositionInCourse
		videos[i].AddedAt = unparsedVideos[i].AddedAt
		videos[i].UpdatedAt = unparsedVideos[i].UpdatedAt
		videos[i].UploadedAt = unparsedVideos[i].UploadedAt
		videos[i].YouTubedAt = unparsedVideos[i].YouTubedAt
		if unparsedVideos[i].CourseID == nil {
			videos[i].Course = nil
		} else {
			videos[i].Course = &api.CourseForAdmin{
				ID:          *unparsedVideos[i].CourseID,
				Name:        *unparsedVideos[i].CourseName,
				Description: *unparsedVideos[i].CourseDescription,
				Position:    *unparsedVideos[i].CoursePosition,
				AddedAt:     *unparsedVideos[i].CourseAddedAt,
				UpdatedAt:   *unparsedVideos[i].CourseUpdatedAt,
			}
			if unparsedVideos[i].ThemeID == nil {
				videos[i].Course.Theme = nil
			} else {
				videos[i].Course.Theme = &models.Theme{
					ID:         *unparsedVideos[i].ThemeID,
					Name:       *unparsedVideos[i].ThemeName,
					Slug:       *unparsedVideos[i].ThemeSlug,
					OnMainPage: *unparsedVideos[i].ThemeOnMainPage,
					AddedAt:    *unparsedVideos[i].ThemeAddedAt,
					UpdatedAt:  *unparsedVideos[i].ThemeUpdatedAt,
					Position:   *unparsedVideos[i].ThemePosition,
					Img:        *unparsedVideos[i].ThemeImg,
				}
			}
		}
	}

	return videos, nil
}

func (r *repo) CountVideosForAdmin(filter *api.Filter) (uint64, error) {
	videosQuery := general.ApplyFilter("video",
		filter, r.builder.Select("count(*)").From("video"))
	sql, args, err := videosQuery.ToSql()
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
