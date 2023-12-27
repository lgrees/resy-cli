package api

import (
	"encoding/json"
	"fmt"

	"github.com/bcillie/resy-cli/internal/utils/http"
)

type VenueResponse struct {
	Venue struct {
		Name string `json:"name"`
	} `json:"venue"`
	LeadTimeInDays int32 `json:"lead_time_in_days"`
}

type VenueDetails struct {
	Name           string
	LeadTimeInDays int32
}

func GetConfig(venueId int32) (*VenueDetails, error) {
	params := make(map[string]string)
	params["venue_id"] = fmt.Sprint(venueId)
	body, statusCode, err := http.Get("https://api.resy.com/2/config", &http.Req{QueryParams: params})
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("failed to fetch venue details for venue id %d", venueId)
	}

	var res VenueResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	return &VenueDetails{
		Name:           res.Venue.Name,
		LeadTimeInDays: res.LeadTimeInDays,
	}, nil
}
