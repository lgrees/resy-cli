package schedule

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/fanniva/resy-cli/internal/utils/http"
)

type searchRequest struct {
	Query   string   `json:"query"`
	PerPage int64    `json:"per_page"`
	Types   []string `json:"types"`
}

type searchResponse struct {
	Search struct {
		Hits []struct {
			Locality string `json:"locality"`
			Rating   struct {
				Average float32 `json:"average"`
				Count   int64   `json:"count"`
			} `json:"rating"`

			Id struct {
				Resy int64 `json:"resy"`
			}
			Name         string   `json:"name"`
			Neighborhood string   `json:"neighborhood"`
			Cuisine      []string `json:"cuisine"`
		} `json:"hits"`
	} `json:"search"`
}

func searchVenues(query string) (*[]surveyVenue, error) {
	searchRequest := searchRequest{
		Query:   query,
		PerPage: 70,
		Types:   []string{"venue"},
	}
	body, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := http.PostJSON("https://api.resy.com/3/venuesearch/search", &http.Req{Body: body})

	if err != nil {
		return nil, err
	}
	if statusCode >= 400 || responseBody == nil {
		return nil, fmt.Errorf("failed to get search results, status code: %d", statusCode)
	}

	var res searchResponse
	_ = json.Unmarshal(responseBody, &res)

	ret := make([]surveyVenue, 0)

	for _, s := range res.Search.Hits {
		v := surveyVenue{
			Id:       strconv.Itoa(int(s.Id.Resy)),
			Name:     s.Name,
			Rating:   fmt.Sprintf("%s (%d reviews)", strconv.FormatFloat(float64(s.Rating.Average), 'f', 2, 64), s.Rating.Count),
			Location: fmt.Sprintf("%s, %s", s.Neighborhood, s.Locality),
			Cuisine:  s.Cuisine[0],
		}
		ret = append(ret, v)
	}

	return &ret, nil
}
