package schedule

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fanniva/resy-cli/internal/book"
	"github.com/fanniva/resy-cli/internal/utils/date"
)

func Add(s string) error {
	inputs, err := surveyDetails()
	if err != nil || inputs == nil {
		return err
	}

	return schedule(inputs)
}

func getBookingDateTime(inputs *surveyInputs) (*time.Time, error) {
	res, err := book.FetchVenueDetails(inputs.Venue.Id)
	if err != nil {
		return nil, err
	}

	slotTime, err := date.ParseTime(inputs.SlotTime)
	if err != nil {
		return nil, err
	}

	reservationDate, err := date.ParseDate(inputs.ReservationDate)
	if err != nil {
		return nil, err
	}

	return date.GetBookingDateTime(res.LeadTimeInDays, slotTime, reservationDate)
}

func schedule(inputs *surveyInputs) error {
	types := strings.Split(inputs.ReservationTypes, "\n")
	_times := strings.Split(inputs.ReservationTimes, "\n")
	times := make([]string, len(_times))
	for i, time := range _times {
		foo, _ := date.ParseTime(time)
		times[i] = date.ToTimeString(foo)
	}

	bookingDateTime, err := getBookingDateTime(inputs)
	if err != nil {
		return err
	}

	bookCmd := book.ToBookCmd(&book.BookingDetails{
		ReservationDate:  inputs.ReservationDate,
		ReservationTimes: times,
		ReservationTypes: types,
		BookingDateTime:  date.ToDateTimeString(bookingDateTime),
		PartySize:        inputs.PartySize,
		VenueId:          inputs.Venue.Id,
	}, inputs.DryRun)

	jobStartTime := bookingDateTime.Add(-time.Minute) // Job should execute slightly before desired book time (at is not terribly reliable here)

	atCmd := fmt.Sprintf("at %s", date.ToAtString(&jobStartTime))
	cmd := fmt.Sprintf("echo \"%s\" | %s", bookCmd, atCmd)
	_, err = exec.Command("sh", "-c", cmd).Output()

	return err
}
