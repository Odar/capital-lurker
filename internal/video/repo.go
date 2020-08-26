package video

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

	unparsedVideos := make([]models.UnparsedVideoForAdmin, 0, limit)
	err = r.postgres.Select(&unparsedVideos, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	videos := make([]api.VideoForAdmin, 0, len(unparsedVideos))
	for i := range unparsedVideos {
		video := api.VideoForAdmin{
			ID:               unparsedVideos[i].ID,
			Name:             unparsedVideos[i].Name,
			Img:              unparsedVideos[i].Img,
			Video:            unparsedVideos[i].Video,
			YouTubeVideo:     unparsedVideos[i].YouTubeVideo,
			PositionInCourse: unparsedVideos[i].PositionInCourse,
			AddedAt:          unparsedVideos[i].AddedAt,
			UpdatedAt:        unparsedVideos[i].UpdatedAt,
			UploadedAt:       unparsedVideos[i].UpdatedAt,
			YouTubedAt:       unparsedVideos[i].YouTubedAt,
		}
		if unparsedVideos[i].CourseID == nil {
			video.Course = nil
		} else {
			video.Course = &api.CourseForAdmin{
				ID:          *unparsedVideos[i].CourseID,
				Name:        *unparsedVideos[i].CourseName,
				Description: *unparsedVideos[i].CourseDescription,
				Position:    *unparsedVideos[i].CoursePosition,
				AddedAt:     *unparsedVideos[i].CourseAddedAt,
				UpdatedAt:   *unparsedVideos[i].CourseUpdatedAt,
			}
			if unparsedVideos[i].ThemeID == nil {
				video.Course.Theme = nil
			} else {
				video.Course.Theme = &models.Theme{
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

		videos = append(videos, video)
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

func (r *repo) DeleteVideo(ID uint64) (int64, error) {
	sql, args, err := r.builder.Delete("*").
		From("video").
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

func (r *repo) UpdateVideoForAdmin(ID uint64, request *api.UpdateVideoForAdminRequest) (*models.Video, error) {
	updateRequest := r.builder.Update("video")
	if request.Name != nil {
		updateRequest = updateRequest.Set("name", *request.Name)
	}
	if request.Img != nil {
		updateRequest = updateRequest.Set("img", *request.Img)
	}
	if request.Video != nil {
		updateRequest = updateRequest.Set("video", *request.Video)
	}
	if request.YouTubeVideo != nil {
		updateRequest = updateRequest.Set("youtube_video", *request.YouTubeVideo)
	}

	if request.CourseID != nil {
		updateRequest = updateRequest.Set("course_id", *request.CourseID)
	}
	if request.PositionInCourse != nil {
		updateRequest = updateRequest.Set("position_in_course", *request.PositionInCourse)
	}
	if request.UploadedAt != nil {
		updateRequest = updateRequest.Set("uploaded_at", *request.UploadedAt)
	}
	if request.YouTubedAt != nil {
		updateRequest = updateRequest.Set("youtubed_at", *request.YouTubedAt)
	}
	updateRequest = updateRequest.Set("updated_at", time.Now().UTC())

	sql, args, err := updateRequest.Suffix(
		"RETURNING "+
			"id, "+
			"name, "+
			"img, "+
			"video, "+
			"youtube_video, "+
			"course_id, "+
			"position_in_course, "+
			"added_at, "+
			"updated_at, "+
			"uploaded_at, "+
			"youtubed_at").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	res := &models.Video{}
	err = r.postgres.Get(res, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return res, nil
}

func (r *repo) AddVideoForAdmin(request *api.AddVideoForAdminRequest) (*models.Video, error) {
	sql, args, err := r.builder.Insert("video").
		Columns("name", "img", "video", "youtube_video", "course_id",
			"position_in_course", "added_at", "updated_at", "uploaded_at", "youtubed_at").
		Values(request.Name, request.Img, request.Video, request.YouTubeVideo, request.CourseID,
			request.PositionInCourse, time.Now().UTC(), time.Now().UTC(), request.UploadedAt, request.YouTubedAt).
		Suffix("RETURNING id, name, img, video, youtube_video, course_id, " +
			"position_in_course, added_at, updated_at, uploaded_at, youtubed_at").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	addedVideo := &models.Video{}
	err = r.postgres.Get(addedVideo, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return addedVideo, nil
}
