package service

import (
	"time"

	"github.com/tmp-friends/anirec-batch/functions/dto"
	"github.com/tmp-friends/anirec-batch/functions/models"
	"github.com/tmp-friends/anirec-batch/functions/pkg/insert_records/dao"
	"github.com/tmp-friends/anirec-batch/functions/pkg/lib"
)

type InsertRecordsService struct {
	dao *dao.InsertRecordsDao
	lib *lib.AnnictClient
}

func NewInsertRecordsService() *InsertRecordsService {
	ird := dao.NewInsertRecordsDao()
	ac := lib.NewAnnictClient()

	return &InsertRecordsService{
		dao: ird,
		lib: ac,
	}
}

func (irs *InsertRecordsService) GetEpisodes() (models.EpisodeSlice, error) {
	currentSeason, err := irs.dao.GetCurrentSeason()
	if err != nil {
		return models.EpisodeSlice{}, err
	}

	episodes, err := irs.dao.GetEpisodes(currentSeason.ID)
	if err != nil {
		return models.EpisodeSlice{}, err
	}

	return episodes, nil
}

func (irs *InsertRecordsService) GetRecords(
	episodes models.EpisodeSlice,
	startTime time.Time,
) ([]dto.Record, error) {
	var records []dto.Record
	for _, episode := range episodes {
		res, err := irs.lib.FetchRecords(episode.ID, startTime)
		if err != nil {
			return []dto.Record{}, err
		}

		records = append(records, res...)

		// Requestを送る間隔をあける
		time.Sleep(1 * time.Second)
	}

	return records, nil
}

func (irs *InsertRecordsService) InsertRecords(records []dto.Record) error {
	err := irs.dao.InsertRecords(records)
	if err != nil {
		return err
	}

	return nil
}
