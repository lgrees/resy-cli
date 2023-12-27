package book

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bcillie/resy-cli/internal/api"
	"github.com/bcillie/resy-cli/internal/utils/date"
	"github.com/rs/zerolog"
)

type BookingDetails struct {
	VenueId string
	// YYYY-MM-DD HH:MM:SS
	BookingDateTime string
	PartySize       string
	// YYYY-MM-DD
	ReservationDate string
	// HH:MM:SS
	ReservationTimes []string
	ReservationTypes []string
}

func (b BookingDetails) MarshalZerologObject(e *zerolog.Event) {
	e.Str("reservation_times", strings.Join(b.ReservationTimes, ",")).
		Str("reservation_types", strings.Join(b.ReservationTypes, ",")).
		Str("reservation_date", b.ReservationDate).
		Str("party_size", b.PartySize).
		Str("venue_id", b.VenueId).
		Str("booking_datetime", b.BookingDateTime)
}

type BookingConfig struct {
	ConfigId  string `json:"config_id"`
	Day       string `json:"day"`
	PartySize int64  `json:"party_size"`
}

func ToBookCmd(bookingDetails *BookingDetails, dryRun bool) string {
	resTypes := make([]string, 0)

	for _, resType := range bookingDetails.ReservationTypes {
		resTypes = append(resTypes, fmt.Sprintf("'%s'", resType))
	}

	types := strings.Join(resTypes, ",")
	times := strings.Join(bookingDetails.ReservationTimes, ",")
	resyExec, _ := os.Executable()

	return fmt.Sprintf("%s book --bookingDateTime='%s' --venueId=%s --partySize=%s --reservationDate=%s --reservationTimes=%s --reservationTypes=%s --dryRun=%t", resyExec, bookingDetails.BookingDateTime, bookingDetails.VenueId, bookingDetails.PartySize, bookingDetails.ReservationDate, times, types, dryRun)
}

func Book(bookingDetails *BookingDetails, dryRun bool, logger zerolog.Logger) error {
	venueId, err := strconv.Atoi(bookingDetails.VenueId)
	if err != nil {
		return err
	}
	partySize, err := strconv.Atoi(bookingDetails.PartySize)
	if err != nil {
		return err
	}
	resDate, err := date.NewResyDate(bookingDetails.ReservationDate, time.DateOnly)
	if err != nil {
		return err
	}
	findParams := api.FindParams{
		VenueId:         int32(venueId),
		PartySize:       int32(partySize),
		ReservationDate: *resDate,
	}
	slots, err := api.Find(&findParams)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch slots")
		return err
	}

	logger.Info().Array("available_slots", slots).Msg("found available slots")

	matchingSlots := findMatches(bookingDetails, slots)
	if len(matchingSlots) == 0 {
		err = errors.New("no matching slots")
		logger.Error().Err(err).Msg("no matching slots")
		return err

	}

	logger.Info().Array("matching_slots", matchingSlots).Msg("found matching slots")

	if dryRun {
		return nil
	}

	err = book(bookingDetails, matchingSlots, logger)
	if err != nil {
		logger.Error().Err(err).Msg("no booking occurred")
		return err
	}
	return nil
}

func WaitThenBook(bookingDetails *BookingDetails, dryRun bool, logger zerolog.Logger) error {
	bookTime, err := date.ParseDateTime(bookingDetails.BookingDateTime)
	if err != nil {
		logger.Error().Err(err).Msg("failed to parse booking datetime")
		return err
	}
	duration := time.Until(*bookTime)
	if duration.Minutes() > 5 {
		err = fmt.Errorf("cannot wait more than 5 minutes to book - it is currently %s before book time", duration.String())
		logger.Error().Err(err).Msg("stopped waiting to book")
		return err
	}

	if duration < 0 {
		logger.Warn().Msg("book time has already passed - this can occur when your computer is asleep/turned off during book time")
	}

	logger.Info().Msgf("waiting %d seconds until booking time: %s", duration/time.Second, bookingDetails.BookingDateTime)
	time.Sleep(duration + (time.Millisecond * 200))

	err = Book(bookingDetails, dryRun, logger)

	if err != nil {
		time.Sleep((time.Millisecond * 100))
		logger.Info().Msg("retrying book job")
		return Book(bookingDetails, dryRun, logger)
	}
	return nil
}

func findMatches(bookingDetails *BookingDetails, slots api.Slots) (matches api.Slots) {
	for _, slot := range slots {
		if isSlotMatch(bookingDetails, slot) {
			matches = append(matches, slot)
		}
	}
	return
}

func book(bookingDetails *BookingDetails, matchingSlots api.Slots, logger zerolog.Logger) error {
	for _, slot := range matchingSlots {
		logger.Info().Object("slot", slot).Msg("attempting to book slot")

		partySize, err := strconv.Atoi(bookingDetails.PartySize)
		if err != nil {
			return err
		}

		resDate, err := date.NewResyDate(bookingDetails.BookingDateTime, time.DateOnly)
		if err != nil {
			return err
		}
		dp := api.DetailsParams{
			ConfigId:  slot.Config.Token,
			Day:       *resDate,
			PartySize: int64(partySize),
		}
		dr, err := api.GetDetails(&dp)
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		err = api.Book(dr)
		if err == nil {
			return nil
		}

		logger.Warn().Err(err).Object("slot", slot).Msg("booking slot failed")
	}

	return errors.New("could not book any matching slots")
}

func isSlotMatch(bookingDetails *BookingDetails, slot api.Slot) bool {
	pieces := strings.Split(slot.Date.Start, " ")
	slotTime := pieces[1]
	slotType := strings.ToLower(slot.Config.Type)
	isTypeMatch := false
	if len(bookingDetails.ReservationTypes) == 0 {
		isTypeMatch = true
	}
	isTimeMatch := false

	for _, time := range bookingDetails.ReservationTimes {
		if time == slotTime {
			isTimeMatch = true
			break
		}
	}
	for _, resType := range bookingDetails.ReservationTypes {
		if resType == slotType {
			isTypeMatch = true
			break
		}
	}

	return isTimeMatch && isTypeMatch
}
