package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/bcillie/resy-cli/internal/utils"
	"github.com/bcillie/resy-cli/internal/utils/date"
	"github.com/bcillie/resy-cli/internal/utils/http"
	"github.com/rs/zerolog"
)

type FindParams struct {
	VenueId   int32 `query:"venue_id"`
	PartySize int32 `query:"party_size"`
	// YYYY-MM-DD
	ReservationDate date.ResyDate `query:"day"`
}

type Slot struct {
	Date struct {
		Start string
	}

	Config struct {
		Type  string
		Token string
	}
}

type Slots []Slot

func (s Slot) MarshalZerologObject(e *zerolog.Event) {
	e.Str("reservation_time", s.Date.Start).
		Str("reservation_type", s.Config.Type)
}

func (s Slots) MarshalZerologArray(a *zerolog.Array) {
	for _, s := range s {
		a.Object(s)
	}
}

type FindResponse struct {
	Results struct {
		Venues []struct {
			Slots Slots
		}
	}
}

func Find(findParams *FindParams) (Slots, error) {

	params := utils.GetQueryParams(*findParams)
	// Seemingly deprecated but still required by the resy API
	params["lat"] = "0"
	params["long"] = "0"

	body, statusCode, err := http.Get("https://api.resy.com/4/find", &http.Req{QueryParams: params})
	if err != nil {
		return nil, err
	}
	// body := make([]byte, 0)
	// statusCode := 400
	// var err error = nil
	if statusCode != 200 {
		return nil, fmt.Errorf("failed to fetch slots for date, status code: %d", statusCode)
	}

	var res FindResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	if len(res.Results.Venues) == 0 || len(res.Results.Venues[0].Slots) == 0 {
		return nil, errors.New("no slots for date")
	}

	return res.Results.Venues[0].Slots, nil
}
