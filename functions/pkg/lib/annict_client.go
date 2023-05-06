package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

func (ac *AnnictClient) GetWorks(seasonName string) (models.WorkSlice, error) {
	req, err := http.NewRequest("GET", URL+"/works", nil)
	if err != nil {
		return models.WorkSlice{}, err
	}

	// QueryParameters
	q := req.URL.Query()
	q.Add("access_token", os.Getenv("ACCESS_TOKEN"))
	q.Add("fields", "id,title")
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

func (ac *AnnictClient) GetEpisodes(workId int) (models.EpisodeSlice, error) {
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
