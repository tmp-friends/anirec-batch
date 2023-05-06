package logic

import (
	"log"

	"github.com/tmp-friends/anirec-batch/functions/pkg/insert_episodes/service"
)

type InsertEpisodesLogic struct {
	service *service.InsertEpisodesService
}

func NewInsertEpisodesLogic() *InsertEpisodesLogic {
	ias := service.NewInsertEpisodesService()

	return &InsertEpisodesLogic{
		service: ias,
	}
}

func (ial *InsertEpisodesLogic) DoExecute() {
	// season_nameは手動で登録する
	// 今期のアニメ作品情報をAnnictAPIから取得する
	works, err := ial.service.GetWorks()
	if err != nil {
		log.Fatal(err)
	}

	workIds := make([]int, len(works))
	for i, work := range works {
		workIds[i] = work.ID
	}

	// 今期のアニメエピソード情報をAnnictAPIから取得する
	episodes, err := ial.service.GetEpisodes(workIds)
	if err != nil {
		log.Fatal(err)
	}

	// DBに今期のアニメ作品情報を登録する
	err = ial.service.InsertWorks(works)
	if err != nil {
		log.Fatal(err)
	}

	err = ial.service.InsertEpisodes(episodes)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("InsertEpisodes Done!\n対象作品数: %d\n対象エピソード数: %d\n", len(works), len(episodes))
}
