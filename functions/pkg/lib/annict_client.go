package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/tmp-friends/anirec-batch/functions/dto"
	"github.com/tmp-friends/anirec-batch/functions/models"
)

const (
	URL      = "https://api.annict.com/v1"
	PER_PAGE = 50
)

type AnnictClient struct {
	client *http.Client
}

func NewAnnictClient() *AnnictClient {
	return &AnnictClient{
		client: &http.Client{},
	}
}

type WorksResponse struct {
	Works      models.WorkSlice `json:"works"`
	TotalCount int              `json:"total_count"`
	NextPage   int              `json:"next_page"`
	PrevPage   int              `json:"prev_page"`
}

func (ac *AnnictClient) FetchWorks(seasonName string) (models.WorkSlice, error) {
	req, err := http.NewRequest("GET", URL+"/works", nil)
	if err != nil {
		return models.WorkSlice{}, err
	}

	// QueryParameters
	q := req.URL.Query()
	q.Add("access_token", os.Getenv("ACCESS_TOKEN"))
	q.Add("fields", "id,title,media")
	q.Add("filter_season", seasonName)
	q.Add("per_page", fmt.Sprintf("%d", PER_PAGE))

	var works models.WorkSlice

	page := 1
	for {
		var worksResponse WorksResponse

		q.Set("page", fmt.Sprintf("%d", page))
		req.URL.RawQuery = q.Encode()

		// Request
		res, err := ac.client.Do(req)
		if err != nil {
			return models.WorkSlice{}, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&worksResponse)
		if err != nil {
			return models.WorkSlice{}, err
		}

		works = append(works, worksResponse.Works...)

		if page*PER_PAGE >= worksResponse.TotalCount {
			break
		}
		page++
	}

	return works, nil
}

type EpisodesResponse struct {
	Episodes   models.EpisodeSlice `json:"episodes"`
	TotalCount int                 `json:"total_count"`
	NextPage   int                 `json:"next_page"`
	PrevPage   int                 `json:"prev_page"`
}

func (ac *AnnictClient) FetchEpisodes(workId int) (models.EpisodeSlice, error) {
	req, err := http.NewRequest("GET", URL+"/episodes", nil)
	if err != nil {
		return models.EpisodeSlice{}, err
	}

	// QueryParameters
	q := req.URL.Query()
	q.Add("access_token", os.Getenv("ACCESS_TOKEN"))
	q.Add("fields", "id,number,sort_number,title")
	q.Add("per_page", fmt.Sprintf("%d", PER_PAGE))
	q.Add("filter_work_id", fmt.Sprintf("%d", workId))

	var episodes models.EpisodeSlice
	page := 1
	for {
		var episodesResponse EpisodesResponse

		q.Set("page", fmt.Sprintf("%d", page))
		req.URL.RawQuery = q.Encode()

		// Request
		res, err := ac.client.Do(req)
		if err != nil {
			return models.EpisodeSlice{}, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&episodesResponse)
		if err != nil {
			return models.EpisodeSlice{}, err
		}

		episodes = append(episodes, episodesResponse.Episodes...)

		if page*PER_PAGE >= episodesResponse.TotalCount {
			break
		}
		page++
	}

	return episodes, nil
}

type RecordsResponse struct {
	Records    []dto.Record `json:"records"`
	TotalCount int          `json:"total_count"`
	NextPage   int          `json:"next_page"`
	PrevPage   int          `json:"prev_page"`
}

func (ac *AnnictClient) FetchRecords(
	episodeId int,
	startTime time.Time,
) ([]dto.Record, error) {
	req, err := http.NewRequest("GET", URL+"/records", nil)
	if err != nil {
		return []dto.Record{}, err
	}

	// QueryParameters
	q := req.URL.Query()
	q.Add("access_token", os.Getenv("ACCESS_TOKEN"))
	q.Add("fields", "id,comment,rating_state,created_at,user.id,work.id")
	q.Add("sort_id", "desc")
	q.Add("per_page", fmt.Sprintf("%d", PER_PAGE))
	q.Add("filter_episode_id", fmt.Sprintf("%d", episodeId))

	var records []dto.Record
	page := 1
	for {
		var recordsResponse RecordsResponse

		q.Set("page", fmt.Sprintf("%d", page))
		req.URL.RawQuery = q.Encode()

		// Request
		res, err := ac.client.Do(req)
		if err != nil {
			return []dto.Record{}, err
		}
		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&recordsResponse)
		if err != nil {
			return []dto.Record{}, err
		}

		records = append(records, recordsResponse.Records...)

		// 終了条件1: 指定した期間より前のレコードが取得された場合
		// (FirestoreにはUpsertを行うので、この条件はAPIリクエスト回数を押さえる目的)
		if len(records) > 0 && records[len(records)-1].CreatedAt.Before(startTime) {
			break
		}
		// 終了条件2: Pageが最後まで到達した場合
		if page*PER_PAGE >= recordsResponse.TotalCount {
			break
		}
		page++
	}

	return records, nil
}
