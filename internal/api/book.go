package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/bcillie/resy-cli/internal/utils/http"
)

func Book(params *DetailsResponse) error {

	token := fmt.Sprintf("book_token=%s", url.PathEscape(params.BookToken.Value))
	var paymentDetails string
	if params.User.PaymentMethods != nil {
		if len(params.User.PaymentMethods) != 0 {
			body, _ := json.Marshal(struct {
				Id int64 `json:"id"`
			}{Id: params.User.PaymentMethods[0].Id})
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
