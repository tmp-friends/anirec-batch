package service

import (
	"time"

	"github.com/tmp-friends/anirec-batch/functions/models"
	"github.com/tmp-friends/anirec-batch/functions/pkg/insert_episodes/dao"
	"github.com/tmp-friends/anirec-batch/functions/pkg/lib"
)

type InsertEpisodesService struct {
	dao *dao.InsertEpisodesDao
	lib *lib.AnnictClient
}

func NewInsertEpisodesService() *InsertEpisodesService {
	ied := dao.NewInsertEpisodesDao()
	ac := lib.NewAnnictClient()

	return &InsertEpisodesService{
		dao: ied,
		lib: ac,
	}
}

func (ies *InsertEpisodesService) GetWorks() (models.WorkSlice, error) {
	// DBから今期のseason_nameを取得する
	currentSeason, err := ies.dao.GetCurrentSeason()
	if err != nil {
		return models.WorkSlice{}, err
	}

	// 今期のアニメ作品情報をAnnictAPIから取得する
	works, err := ies.lib.FetchWorks(currentSeason.Name)
	if err != nil {
		return models.WorkSlice{}, err
	}

	for _, v := range works {
		v.SeasonID = currentSeason.ID
	}

	return works, nil
}

func (ies *InsertEpisodesService) GetEpisodes(workIds []int) (models.EpisodeSlice, error) {
	// 今期のアニメエピソード情報をAnnictAPIから取得する
	var episodes models.EpisodeSlice
	for _, workId := range workIds {
		res, err := ies.lib.FetchEpisodes(workId)
		if err != nil {
			return models.EpisodeSlice{}, err
		}

		for _, v := range res {
			v.WorkID = workId
		}

		episodes = append(episodes, res...)

		// Requestを送る間隔をあける
		time.Sleep(1 * time.Second)
	}

	return episodes, nil
}

func (ies *InsertEpisodesService) InsertWorks(works models.WorkSlice) error {
	err := ies.dao.InsertWorks(works)

	return err
}

func (ies *InsertEpisodesService) InsertEpisodes(episodes models.EpisodeSlice) error {
	err := ies.dao.InsertEpisodes(episodes)

	return err
}
