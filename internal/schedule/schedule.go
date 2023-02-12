package schedule

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/lgrees/resy-cli/internal/book"
	"github.com/lgrees/resy-cli/internal/utils/date"
)

func Add(s string) error {
	inputs, err := surveyDetails()
	if err != nil || inputs == nil {
		return err
	}

	return schedule(inputs)
}

func schedule(inputs *surveyInputs) error {
	types := strings.Split(inputs.ReservationTypes, "\n")
	times := strings.Split(inputs.ReservationTimes, "\n")

	bookCmd := book.ToBookCmd(&book.BookingDetails{
		ReservationDate:  inputs.ReservationDate,
		ReservationTimes: times,
		ReservationTypes: types,
		BookingDateTime:  inputs.BookingDateTime,
		PartySize:        inputs.PartySize,
		VenueId:          inputs.Venue.Id,
	}, inputs.DryRun)

	t, _ := date.ParseDateTime(inputs.BookingDateTime)
	bookTime := t.Add(-time.Minute) // Job should execute slightly before desired book time (at is not terribly reliable here)

	atCmd := fmt.Sprintf("at %s", date.ToAtString(&bookTime))
	cmd := fmt.Sprintf("echo \"%s\" | %s", bookCmd, atCmd)
	_, err := exec.Command("sh", "-c", cmd).Output()

	return err
}
