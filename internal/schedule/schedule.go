package schedule

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/lgrees/resy-cli/internal/book"
	"github.com/lgrees/resy-cli/internal/utils/date"
	"github.com/lgrees/resy-cli/internal/utils/surveyHelpers"

	"github.com/AlecAivazis/survey/v2"
)

type Inputs struct {
	DryRun           bool
	VenueId          string
	BookingDateTime  string
	PartySize        string
	ReservationDate  string
	ReservationTimes string
	ReservationTypes string
}

func Add(s string) error {
	inputs, err := surveyDetails()
	if err != nil || inputs == nil {
		return err
	}

	return schedule(inputs)
}

func schedule(inputs *Inputs) error {
	types := strings.Split(inputs.ReservationTypes, "\n")
	times := strings.Split(inputs.ReservationTimes, "\n")

	bookCmd := book.ToBookCmd(&book.BookingDetails{
		ReservationDate:  inputs.ReservationDate,
		ReservationTimes: times,
		ReservationTypes: types,
		BookingDateTime:  inputs.BookingDateTime,
		PartySize:        inputs.PartySize,
		VenueId:          inputs.VenueId,
	}, inputs.DryRun)

	t, _ := date.ParseDateTime(inputs.BookingDateTime)
	bookTime := t.Add(-time.Minute) // Job should execute slightly before desired book time (at is not terribly reliable here)

	atCmd := fmt.Sprintf("at %s", date.ToAtString(&bookTime))
	cmd := fmt.Sprintf("echo \"%s\" | %s", bookCmd, atCmd)
	_, err := exec.Command("sh", "-c", cmd).Output()

	return err
}

var questions = []*survey.Question{
	{
		Name:     "venueId",
		Prompt:   &survey.Input{Message: "Venue Id:"},
		Validate: survey.ComposeValidators(survey.Required, surveyHelpers.CreateRegexValidator("[0-9]+", "Venue Id must be a number.")),
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

func surveyDetails() (*Inputs, error) {
	answers := Inputs{}
	err := survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}

	confirm := false
	survey.AskOne(&survey.Confirm{Message: "Schedule to book with the above information?"}, &confirm)
	if confirm {
		fmt.Printf("\nGreat, I'll attempt to book your reservation at %s.", answers.BookingDateTime)
		fmt.Println("\nMake sure that your credentials are up to date before then by running resy ping!")
		return &answers, err
	} else {
		fmt.Println("Okay, I won't try to book anything.")
		return nil, nil
	}
}
