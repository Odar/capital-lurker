package speaker

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

func (r *repo) GetSpeakersOnMain(limit int64) ([]api.SpeakerOnMain, error) {
	if limit <= 0 {
		return nil, errors.New("bad request; parameter: limit should be positive")
	}
	sql, args, err := r.builder.Select(
		"speaker.id AS speaker_id, " +
			"speaker.name AS speaker_name, " +
			"speaker.position AS speaker_position, " +
			"speaker.img AS speaker_img, " +
			"university.id AS university_id, " +
			"university.name AS university_name, " +
			"university.img AS university_img").
		From("speaker").
		Where("speaker.on_main_page = true").
		OrderBy("speaker.position DESC").
		Limit(uint64(limit)).
		LeftJoin("university ON speaker.university_id = university.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	unparsedSpeakers := make([]models.UnparsedSpeakerOnMain, 0, limit)
	err = r.postgres.Select(&unparsedSpeakers, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	speakers := make([]api.SpeakerOnMain, 0, len(unparsedSpeakers))
	for i := range unparsedSpeakers {
		speaker := api.SpeakerOnMain{
			ID:       unparsedSpeakers[i].ID,
			Name:     unparsedSpeakers[i].Name,
			Position: unparsedSpeakers[i].Position,
			Img:      unparsedSpeakers[i].Img,
		}
		if unparsedSpeakers[i].UniversityID == nil {
			speaker.University = nil
		} else {
			speaker.University = &api.UniversityOnMain{
				ID:   *unparsedSpeakers[i].UniversityID,
				Name: *unparsedSpeakers[i].UniversityName,
				Img:  *unparsedSpeakers[i].UniversityImg,
			}
		}

		speakers = append(speakers, speaker)
	}

	return speakers, nil
}

func (r *repo) GetSpeakersForAdmin(limit int64, page int64, sortBy string, filter *api.Filter) (
	[]api.SpeakerForAdmin, error) {
	sortBy = general.ApplySortByParameter(sortBy)
	speakersQuery := general.ApplyFilter("speaker", filter, r.builder.Select(
		"speaker.id AS speaker_id, "+
			"speaker.name AS speaker_name, "+
			"speaker.on_main_page AS speaker_on_main_page, "+
			"speaker.in_filter AS speaker_in_filter, "+
			"speaker.added_at AS speaker_added_at, "+
			"speaker.updated_at AS speaker_updated_at, "+
			"speaker.position AS speaker_position, "+
			"speaker.img AS speaker_img, "+
			"university.id AS university_id, "+
			"university.name AS university_name, "+
			"university.on_main_page AS university_on_main_page, "+
			"university.in_filter AS university_in_filter, "+
			"university.added_at AS university_added_at, "+
			"university.updated_at AS university_updated_at, "+
			"university.position AS university_position, "+
			"university.img AS university_img").
		From("speaker")).
		Limit(uint64(limit)).Offset(uint64((page - 1) * limit)).OrderBy("speaker." + sortBy)
	sql, args, err := speakersQuery.LeftJoin("university ON speaker.university_id = university.id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	unparsedSpeakers := make([]models.UnparsedSpeakerForAdmin, 0, limit)
	err = r.postgres.Select(&unparsedSpeakers, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	speakers := make([]api.SpeakerForAdmin, 0, len(unparsedSpeakers))
	for i := range unparsedSpeakers {
		speaker := api.SpeakerForAdmin{
			ID:         unparsedSpeakers[i].ID,
			Name:       unparsedSpeakers[i].Name,
			OnMainPage: unparsedSpeakers[i].OnMainPage,
			InFilter:   unparsedSpeakers[i].InFilter,
			AddedAt:    unparsedSpeakers[i].AddedAt,
			UpdatedAt:  unparsedSpeakers[i].UpdatedAt,
			Position:   unparsedSpeakers[i].Position,
			Img:        unparsedSpeakers[i].Img,
		}
		if unparsedSpeakers[i].UniversityID == nil {
			speaker.University = nil
		} else {
			speaker.University = &models.University{
				ID:         *unparsedSpeakers[i].UniversityID,
				Name:       *unparsedSpeakers[i].UniversityName,
				OnMainPage: *unparsedSpeakers[i].UniversityOnMainPage,
				InFilter:   *unparsedSpeakers[i].UniversityInFilter,
				AddedAt:    *unparsedSpeakers[i].UniversityAddedAt,
				UpdatedAt:  *unparsedSpeakers[i].UniversityUpdatedAt,
				Position:   *unparsedSpeakers[i].UniversityPosition,
				Img:        *unparsedSpeakers[i].UniversityImg,
			}
		}

		speakers = append(speakers, speaker)
	}

	return speakers, nil
}

func (r *repo) CountSpeakersForAdmin(filter *api.Filter) (uint64, error) {
	speakersQuery := general.ApplyFilter("speaker",
		filter, r.builder.Select("count(*)").From("speaker"))
	sql, args, err := speakersQuery.ToSql()
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

func (r *repo) DeleteSpeaker(ID uint64) (int64, error) {
	sql, args, err := r.builder.Delete("*").
		From("speaker").
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

func (r *repo) UpdateSpeakerForAdmin(ID uint64, request *api.UpdateSpeakerForAdminRequest) (
	*models.Speaker, error) {
	updateRequest := r.builder.Update("speaker")
	if request.Name != nil {
		updateRequest = updateRequest.Set("name", *request.Name)
	}
	if request.OnMainPage != nil {
		updateRequest = updateRequest.Set("on_main_page", *request.OnMainPage)
	}
	if request.InFilter != nil {
		updateRequest = updateRequest.Set("in_filter", *request.InFilter)
	}
	updateRequest = updateRequest.Set("updated_at", time.Now().UTC())
	if request.Position != nil {
		updateRequest = updateRequest.Set("position", *request.Position)
	}
	if request.Img != nil {
		updateRequest = updateRequest.Set("img", *request.Img)
	}
	if request.UniversityID != nil {
		updateRequest = updateRequest.Set("university_id", *request.UniversityID)
	}

	sql, args, err := updateRequest.Suffix(
		"RETURNING "+
			"id, "+
			"name, "+
			"on_main_page, "+
			"in_filter, "+
			"added_at, "+
			"updated_at, "+
			"position, "+
			"img"+
			"university_id").
		Where("id = ?", ID).
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	res := &models.Speaker{}
	err = r.postgres.Get(res, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return res, nil
}

func (r *repo) AddSpeakerForAdmin(request *api.AddSpeakerForAdminRequest) (*models.Speaker, error) {
	sql, args, err := r.builder.Insert("speaker").
		Columns("name", "on_main_page", "in_filter", "added_at", "updated_at",
			"position", "img", "university_id").
		Values(request.Name, request.OnMainPage, request.InFilter, time.Now().UTC(), time.Now().UTC(),
			request.Position, request.Img, request.UniversityID).
		Suffix("RETURNING id, name, on_main_page, in_filter, added_at, updated_at," +
			"position, img, university_id").
		ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "can not build sql")
	}

	addedSpeaker := &models.Speaker{}
	err = r.postgres.Get(addedSpeaker, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "can not exec query `%s` with args %+v", sql, args)
	}

	return addedSpeaker, nil
}
