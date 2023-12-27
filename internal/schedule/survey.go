package schedule

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bcillie/resy-cli/internal/api"
	"github.com/bcillie/resy-cli/internal/utils/surveyHelpers"
)

type surveyVenue struct {
	Name     string
	Location string
	Rating   string
	Cuisine  string
	Id       int
}

type surveyInputs struct {
	DryRun           bool
	Venue            surveyVenue
	SlotTime         string
	PartySize        string
	ReservationDate  string
	ReservationTimes string
	ReservationTypes string
}

func (venue *surveyVenue) WriteAnswer(name string, value interface{}) error {
	s := value.(string)
	arr := strings.Split(s, " | ")
	if len(arr) < 5 {
		return nil
	}
	venue.Name = arr[0]
	fmt.Sscanf(arr[4], "ID: %d", &venue.Id)
	return nil
}

func suggestVenues(toComplete string) []string {
	venues, err := api.SearchVenues(toComplete)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	ret := make([]string, 0)
	for _, v := range *venues {
		ret = append(ret, v.String())
	}

	return ret
}

var questions = []*survey.Question{
	{
		Name: "venue",
		Prompt: &survey.Input{
			Message: "Venue:",
			Suggest: suggestVenues,
		},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.VenueValidator),
	},
	{
		Name:     "partySize",
		Prompt:   &survey.Input{Message: "Party Size:"},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.CreateRegexValidator("[0-9]+", "Party Size must be a number.")),
	},
	{
		Name:     "reservationDate",
		Prompt:   &survey.Input{Message: "Reservation Date (YYYY-MM-DD):"},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.DateValidator),
	},
	{
		Name:     "reservationTimes",
		Prompt:   &survey.Multiline{Message: "Reservation Times (HH:MM):"},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.TimesValidator),
	},
	{
		Name: "reservationTypes",
		Prompt: &survey.Multiline{
			Message: "Reservation Types (ex. 'Indoor dining') - optional:",
			Help:    "Generally, this corresponds directly to the tag that you see under the reservation (though not always). Leave this empty to book any type of reservation.",
		},
		Transform: surveyHelpers.TransformLowerCase,
	},
	{
		Name: "slotTime",
		Prompt: &survey.Input{
			Message: "What time do slots open? (HH:MM)",
		},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.TimeValidator),
	},
	{
		Name: "dryRun",
		Prompt: &survey.Confirm{
			Message: "Is this a dry run?",
			Default: false,
			Help:    "Dry runs will not actually attempt to book your reservation."},
		Validate: survey.Required,
	},
}

func surveyDetails() (*surveyInputs, error) {
	answers := surveyInputs{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}

	bookingDateTime, err := getBookingDateTime(&answers)
	if err != nil {
		return nil, err
	}

	confirm := false
	survey.AskOne(&survey.Confirm{Message: "Schedule to book with the above information?"}, &confirm)
	if confirm {
		fmt.Printf(`
		Great, resy-cli will attempt to book your reservation for a party of %s at %s.
		The booking will be attempted at %s.
		Make sure that your credentials are up to date before then by running 'resy ping'.
		Additionally, make sure that your computer is awake at this time.

		Happy dining! 😋
		`, answers.PartySize, answers.Venue.Name, bookingDateTime)
		fmt.Println("")
		return &answers, nil
	} else {
		fmt.Println("Okay, I won't try to book anything.")
		return nil, nil
	}
}
