package logic

import (
	"log"
	"time"

	"github.com/tmp-friends/anirec-batch/functions/pkg/insert_records/service"
)

const LAYOUT = "2006-01-02 15:04:05"

type InsertRecordsLogic struct {
	service *service.InsertRecordsService
}

func NewInsertRecordsLogic() *InsertRecordsLogic {
	irs := service.NewInsertRecordsService()

	return &InsertRecordsLogic{
		service: irs,
	}
}

func (irl *InsertRecordsLogic) DoExecute() {
	startTime := irl.setStartTime()
	log.Printf("Target within time: %s ~ \n", startTime)

	episodes, err := irl.service.GetEpisodes()
	if err != nil {
		log.Fatal(err)
	}

	records, err := irl.service.GetRecords(episodes, startTime)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(records)

	err = irl.service.InsertRecords(records)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("InsertRecords Done!\n対象レコード数: %d\n", len(records))
}

// (日本標準時間で)開始時刻を取得
func (irl *InsertRecordsLogic) setStartTime() time.Time {
	// @see: https://tutuz-tech.hatenablog.com/entry/2021/01/30/192956
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}

	// 現在時刻の取得
	nowJST := time.Now().In(jst)

	startTime := time.Date(
		nowJST.Year(),
		nowJST.Month(),
		nowJST.Day()-1,
		4,
		0,
		0,
		0,
		jst,
	)

	return startTime
}
