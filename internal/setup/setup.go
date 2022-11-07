package setup

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/viper"
)

var questions = []*survey.Question{
	{
		Name:     "apiKey",
		Prompt:   &survey.Input{Message: "Api Key:"},
		Validate: survey.Required,
	},
	{
		Name:     "authToken",
		Prompt:   &survey.Input{Message: "Auth Token:"},
		Validate: survey.Required,
	},
}

func SurveyConfig() error {
	answers := struct {
		ApiKey    string
		AuthToken string
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		return err
	}

	confirm := false
	survey.AskOne(&survey.Confirm{Message: "Does this look correct?"}, &confirm)

	if confirm {
		viper.Set("resy_api_key", answers.ApiKey)
		viper.Set("resy_auth_token", answers.AuthToken)
		viper.WriteConfig()
		fmt.Println("Your user info has been saved! You're all set to start booking.")

	} else {
		fmt.Println("Your user info was not saved.")
	}

	return nil
}
