package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/bcillie/resy-cli/internal/utils/http"
)

type FindResponse struct {
	Results struct {
		Venues []struct {
			Slots Slots
		}
	}
}

type VenueResponse struct {
	Venue struct {
		Name string `json:"name"`
	} `json:"venue"`
	LeadTimeInDays int32 `json:"lead_time_in_days"`
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

func FetchVenueDetails(venueId string) (*VenueDetails, error) {
	params := make(map[string]string)
	params["venue_id"] = venueId

	body, statusCode, err := http.Get("https://api.resy.com/2/config", &http.Req{QueryParams: params})
	if err != nil {
		return nil, err
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("failed to fetch venue details for venue id %s", venueId)
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

func fetchSlots(bookingDetails *BookingDetails) (Slots, error) {
	params := make(map[string]string)
	params["party_size"] = bookingDetails.PartySize
	params["venue_id"] = bookingDetails.VenueId
	params["day"] = bookingDetails.ReservationDate
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

func bookSlot(bookingDetails *BookingDetails, slot Slot) error {
	// Get booking token
	partySize, _ := strconv.Atoi(bookingDetails.PartySize)
	bookingConfig := BookingConfig{
		ConfigId:  slot.Config.Token,
		Day:       bookingDetails.ReservationDate,
		PartySize: int64(partySize),
	}
	body, err := json.Marshal(bookingConfig)
	if err != nil {
		return err
	}
	responseBody, statusCode, err := http.PostJSON("https://api.resy.com/3/details", &http.Req{Body: body})
	if err != nil {
		return err
	}
	if statusCode >= 400 || responseBody == nil {
		return fmt.Errorf("failed to get booking details, status code: %d", statusCode)
	}

	var details DetailsResponse
	_ = json.Unmarshal(responseBody, &details)

	// Actually book with token
	token := fmt.Sprintf("book_token=%s", url.PathEscape(details.BookToken.Value))
	var paymentDetails string
	if details.User.PaymentMethods != nil {
		if len(details.User.PaymentMethods) != 0 {
			body, _ := json.Marshal(struct {
				Id int64 `json:"id"`
			}{Id: details.User.PaymentMethods[0].Id})
			paymentDetails = fmt.Sprintf("struct_payment_method=%s", url.PathEscape(string(body)))
		}
	}

	var form string
	if paymentDetails != "" {
		form = strings.Join([]string{token, paymentDetails}, "&")
	} else {
		form = token
	}
	_, statusCode, err = http.PostForm("https://api.resy.com/3/book", &http.Req{Body: []byte(form)})
	if err != nil {
		return err
	}
	if statusCode >= 400 {
		return fmt.Errorf("failed to book reservation, status code: %d", statusCode)
	}

	return nil
}
