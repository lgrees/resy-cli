package schedule

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/bcillie/resy-cli/internal/api"
	"github.com/bcillie/resy-cli/internal/book"
	"github.com/bcillie/resy-cli/internal/utils/date"
)

func Add(s string) error {
	inputs, err := surveyDetails()
	if err != nil || inputs == nil {
		return err
	}

	return schedule(inputs)
}

func getBookingDateTime(inputs *surveyInputs) (*date.ResyDate, error) {
	res, err := api.GetConfig(int32(inputs.Venue.Id))
	if err != nil {
		return nil, err
	}

	slotTime, err := date.NewResyDate(inputs.SlotTime, "15:04")
	if err != nil {
		return nil, err
	}

	reservationDate, err := date.NewResyDate(inputs.ReservationDate, time.DateOnly)
	if err != nil {
		return nil, err
	}

	return date.GetBookingDateTime(res.LeadTimeInDays, slotTime, reservationDate)
}

func schedule(inputs *surveyInputs) error {
	types := strings.Split(inputs.ReservationTypes, "\n")
	_times := strings.Split(inputs.ReservationTimes, "\n")
	times := make([]string, len(_times))
	for i, t := range _times {
		foo, _ := date.ParseTime(t)
		times[i] = foo.Format(time.TimeOnly)
	}

	bookingDateTime, err := getBookingDateTime(inputs)
	if err != nil {
		return err
	}

	bookCmd := book.ToBookCmd(&book.BookingDetails{
		ReservationDate:  inputs.ReservationDate,
		ReservationTimes: times,
		ReservationTypes: types,
		BookingDateTime:  bookingDateTime.String(),
		PartySize:        inputs.PartySize,
		VenueId:          fmt.Sprintf("%d", inputs.Venue.Id),
	}, inputs.DryRun)

	jobStartTime := bookingDateTime.Add(-time.Minute) // Job should execute slightly before desired book time (at is not terribly reliable here)

	atCmd := fmt.Sprintf("at %s", jobStartTime.Format(date.AtFmt))
	cmd := fmt.Sprintf("echo \"%s\" | %s", bookCmd, atCmd)
	_, err = exec.Command("sh", "-c", cmd).Output()

	return err
}
