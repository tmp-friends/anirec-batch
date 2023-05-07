package function

import (
	"log"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/joho/godotenv"
	insertEpisodesLogic "github.com/tmp-friends/anirec-batch/functions/pkg/insert_episodes/logic"
	insertRecordsLogic "github.com/tmp-friends/anirec-batch/functions/pkg/insert_records/logic"
)

// Targetの定義
func init() {
	loadEnv()

	functions.HTTP("InsertEpisodes", insertEpisodes)
	functions.HTTP("InsertRecords", insertRecords)
}

func loadEnv() {
	// SOURCE_DIRは本番のみGCPConsole上で設定
	// 開発: "", 本番: "serverless_function_source_code/"
	err := godotenv.Load(os.Getenv("SOURCE_DIR") + ".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func insertEpisodes(w http.ResponseWriter, r *http.Request) {
	iel := insertEpisodesLogic.NewInsertEpisodesLogic()
	iel.DoExecute()
}

func insertRecords(w http.ResponseWriter, r *http.Request) {
	irl := insertRecordsLogic.NewInsertRecordsLogic()
	irl.DoExecute()
}
