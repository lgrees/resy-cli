package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bcillie/resy-cli/internal/utils"
	"github.com/bcillie/resy-cli/internal/utils/http"
)

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

type FindResponse struct {
	Results struct {
		Venues []struct {
			Slots Slots
		}
	}
}

type FindParams struct {
	VenueId   int32 `query:"venue_id"`
	PartySize int32 `query:"party_size"`
	// YYYY-MM-DD
	ReservationDate time.Time `query:"day" fmt:"2006-01-06"`
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
