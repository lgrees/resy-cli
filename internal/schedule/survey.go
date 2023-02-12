package schedule

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lgrees/resy-cli/internal/utils/surveyHelpers"
)

type surveyVenue struct {
	Name     string
	Location string
	Rating   string
	Cuisine  string
	Id       string
}

type surveyInputs struct {
	DryRun           bool
	Venue            surveyVenue
	BookingDateTime  string
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
	venue.Id = arr[4]
	return nil
}

func (venue *surveyVenue) ToString() string {
	return strings.Join([]string{
		venue.Name,
		venue.Cuisine,
		venue.Location,
		venue.Rating,
		venue.Id,
	}, " | ")
}

func suggestVenues(toComplete string) []string {
	venues, err := searchVenues(toComplete)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	ret := make([]string, 0)
	for _, v := range *venues {
		ret = append(ret, v.ToString())
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
		Prompt:   &survey.Multiline{Message: "Reservation Times (HH:MM:SS):"},
		Validate: survey.Required,
	},
	{
		Name: "reservationTypes",
		Prompt: &survey.Multiline{
			Message: "Reservation Types (ex. 'Indoor dining'):",
			Help:    "Generally, this corresponds directly to the tag that you see under the reservation (though not always). Leave this empty to book any type of reservation.",
		},
		Transform: surveyHelpers.TransformLowerCase,
	},
	{
		Name: "bookingDateTime",
		Prompt: &survey.Input{
			Message: "What date/time should resy-cli attempt to book this reservation? (YYYY-MM-DD HH:MM:SS)",
			Help:    "Generally, this should be when the restaurant opens slots for the date you are trying to book."},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.DateTimeValidator),
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

	confirm := false
	survey.AskOne(&survey.Confirm{Message: "Schedule to book with the above information?"}, &confirm)
	if confirm {
		fmt.Printf("\nGreat, I'll attempt to book your reservation for %s at %s on %s.", answers.PartySize, answers.Venue.Name, answers.BookingDateTime)
		fmt.Println("\nMake sure that your credentials are up to date before then by running `resy ping`!")
		return &answers, err
	} else {
		fmt.Println("Okay, I won't try to book anything.")
		return nil, nil
	}
}
