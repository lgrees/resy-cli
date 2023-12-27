package book

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/bcillie/resy-cli/internal/api"
	"github.com/bcillie/resy-cli/internal/utils/date"
	"github.com/bcillie/resy-cli/internal/utils/http"
)

func bookSlot(bookingDetails *BookingDetails, slot api.Slot) error {
	// Get booking token
	partySize, err := strconv.Atoi(bookingDetails.PartySize)
	if err != nil {
		return err
	}

	resDate, err := date.NewResyDate(bookingDetails.ReservationDate, time.DateOnly)
	if err != nil {
		return err
	}

	detailsParams := api.DetailsParams{
		ConfigId:  slot.Config.Token,
		Day:       *resDate,
		PartySize: int64(partySize),
	}

	details, err := api.GetDetails(&detailsParams)
	if err != nil {
		return err
	}

	// book with token
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
	_, statusCode, err := http.PostForm("https://api.resy.com/3/book", &http.Req{Body: []byte(form)})
	if err != nil {
		return err
	}
	if statusCode >= 400 {
		return fmt.Errorf("failed to book reservation, status code: %d", statusCode)
	}

	return nil
}
