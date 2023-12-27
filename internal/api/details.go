package api

import (
	"encoding/json"
	"fmt"

	"github.com/bcillie/resy-cli/internal/utils/date"
	"github.com/bcillie/resy-cli/internal/utils/http"
)

type DetailsParams struct {
	ConfigId  string        `json:"config_id"`
	Day       date.ResyDate `json:"day"`
	PartySize int64         `json:"party_size"`
}

type DetailsResponse struct {
	BookToken struct {
		Value string `json:"value"`
	} `json:"book_token"`
	User struct {
		PaymentMethods []struct {
			Id int64 `json:"id"`
		} `json:"payment_methods"`
	} `json:"user"`
}

func GetDetails(detailsParams *DetailsParams) (*DetailsResponse, error) {
	body, err := json.Marshal(detailsParams)
	if err != nil {
		return nil, err
	}
	responseBody, statusCode, err := http.PostJSON("https://api.resy.com/3/details", &http.Req{Body: body})
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 || responseBody == nil {
		return nil, fmt.Errorf("failed to get booking details, status code: %d", statusCode)
	}

	var details DetailsResponse
	err = json.Unmarshal(responseBody, &details)
	if err != nil {
		return nil, err
	}

	return &details, nil
}
