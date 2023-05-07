package dao

import (
	"context"
	"database/sql"

	"github.com/tmp-friends/anirec-batch/functions/models"
	"github.com/tmp-friends/anirec-batch/functions/pkg/config"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type InsertEpisodesDao struct {
	DB *sql.DB
}

func NewInsertEpisodesDao() *InsertEpisodesDao {
	db := config.NewPostgresConnector()

	return &InsertEpisodesDao{
		DB: db.Conn,
	}
}

func (ied *InsertEpisodesDao) GetCurrentSeason() (*models.Season, error) {
	currentSeasonName, err := models.Seasons(
		qm.Select("id", "name"),
		qm.OrderBy("updated_at desc"),
	).One(context.Background(), ied.DB)
	if err != nil {
		return &models.Season{}, err
	}

	return currentSeasonName, nil
}

func (ied *InsertEpisodesDao) InsertWorks(works models.WorkSlice) error {
	for _, v := range works {
		err := v.Upsert(
			context.Background(),
			ied.DB,
			true,
			[]string{"id"},
			boil.Blacklist("created_at"),
			boil.Infer(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ied *InsertEpisodesDao) InsertEpisodes(episodes models.EpisodeSlice) error {
	for _, v := range episodes {
		err := v.Upsert(
			context.Background(),
			ied.DB,
			true,
			[]string{"id"},
			boil.Blacklist("created_at"),
			boil.Infer(),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
