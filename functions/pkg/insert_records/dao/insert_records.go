package dao

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"google.golang.org/api/option"

	"github.com/tmp-friends/anirec-batch/functions/dto"
	"github.com/tmp-friends/anirec-batch/functions/models"
	"github.com/tmp-friends/anirec-batch/functions/pkg/config"
)

type InsertRecordsDao struct {
	DB     *sql.DB
	client *firestore.Client
}

func NewInsertRecordsDao() *InsertRecordsDao {
	db := config.NewPostgresConnector()
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	firestoreClient, err := firestore.NewClient(context.Background(), "anirec-385905", opt)
	if err != nil {
		log.Fatalf("Failed to create firestore client: %v", err)
	}

	return &InsertRecordsDao{
		DB:     db.Conn,
		client: firestoreClient,
	}
}

func (ird *InsertRecordsDao) GetCurrentSeason() (*models.Season, error) {
	currentSeasonName, err := models.Seasons(
		qm.Select("id", "name"),
		qm.OrderBy("updated_at desc"),
	).One(context.Background(), ird.DB)
	if err != nil {
		return &models.Season{}, err
	}

	return currentSeasonName, nil
}

func (ird *InsertRecordsDao) GetEpisodes(seasonId int) (models.EpisodeSlice, error) {
	episodes, err := models.Episodes(
		qm.Select(
			"episodes.id as id",
			"episodes.number as number",
			"episodes.sort_number as sort_number",
			"episodes.title as title",
			"episodes.work_id as work_id",
		),
		qm.InnerJoin("works on episodes.work_id = works.id"),
		qm.Where("works.season_id = ?", seasonId),
		qm.Where("works.media = ?", "tv"),
	).All(context.Background(), ird.DB)

	if err != nil {
		return models.EpisodeSlice{}, err
	}

	return episodes, nil
}

func (ird *InsertRecordsDao) InsertRecords(records []dto.Record) error {
	// firestoreへUpsert
	// TODO: 500件ずつに分けてBatchWriteする
	for _, v := range records {
		data := map[string]interface{}{
			"comment":      v.Comment,
			"rating_state": v.RatingState,
			"user_id":      v.User.ID,
			"work_id":      v.Work.ID,
			"created_at":   v.CreatedAt,
		}

		docRef := ird.client.Collection("records").Doc(strconv.Itoa(v.ID))
		_, err := docRef.Set(context.Background(), data, firestore.MergeAll)
		if err != nil {
			return err
		}
	}

	return nil
}
