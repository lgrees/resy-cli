package api

import (
	"encoding/json"
	"fmt"

	"github.com/bcillie/resy-cli/internal/utils/http"
)

type Response struct {
	Message string
}

func Ping() {
	body, statusCode, err := http.Get("https://api.resy.com/2/user", &http.Req{})

	var jsonObj Response
	json.Unmarshal(body, &jsonObj)

	if err != nil {
		fmt.Printf("Error: could not ping the auth server: %s\n", err)
	} else if statusCode >= 400 {
		fmt.Println("Error: Could not authenticate with resy.")

		fmt.Printf("Status Code: %d\n", statusCode)
		if jsonObj.Message != "" {
			fmt.Printf("Message: %s\n", jsonObj.Message)
		}

		fmt.Println("Run `resy setup` to reset your authentication information.")
	} else {
		fmt.Println("Success! You're all set to begin booking.")
	}
}
