package api

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bcillie/resy-cli/internal/utils/http"
)

const SearchUrl = "https://api.resy.com/3/venuesearch/search"

type SearchRequest struct {
	Query   string   `json:"query"`
	PerPage int64    `json:"per_page"`
	Types   []string `json:"types"`
}

type SearchHit struct {
	Locality string `json:"locality"`
	Rating   struct {
		Average float32 `json:"average"`
		Count   int64   `json:"count"`
	} `json:"rating"`

	Id struct {
		Resy int32 `json:"resy"`
	}
	Name         string   `json:"name"`
	Neighborhood string   `json:"neighborhood"`
	Cuisine      []string `json:"cuisine"`
}

type SearchResponse struct {
	Search struct {
		Hits []SearchHit `json:"hits"`
	} `json:"search"`
}

type SearchResult struct {
	Name     string
	Location string
	Rating   float32
	Cuisine  []string
	Id       int32
}

func NewSearchResult(hit SearchHit) *SearchResult {
	return &SearchResult{
		Name:     hit.Name,
		Location: hit.Locality,
		Rating:   hit.Rating.Average,
		Cuisine:  hit.Cuisine,
		Id:       hit.Id.Resy,
	}
}

func (result *SearchResult) String() string {
	return strings.Join([]string{
		result.Name,
		strings.Join(result.Cuisine, ", "),
		result.Location,
		fmt.Sprintf("%.2f stars", result.Rating),
		fmt.Sprintf("ID: %d", result.Id),
	}, " | ")
}

func SearchVenues(query string) (*[]SearchResult, error) {
	searchRequest := SearchRequest{
		Query:   query,
		PerPage: 5,
		Types:   []string{"venue"},
	}
	body, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := http.PostJSON(SearchUrl, &http.Req{Body: body})

	if err != nil {
		return nil, err
	}
	if statusCode >= 400 || responseBody == nil {
		return nil, fmt.Errorf("failed to get search results, status code: %d", statusCode)
	}

	var res SearchResponse
	_ = json.Unmarshal(responseBody, &res)

	ret := make([]SearchResult, 0)

	for _, hit := range res.Search.Hits {
		ret = append(ret, *NewSearchResult(hit))
	}

	return &ret, nil
}
