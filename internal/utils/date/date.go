package date

import (
	"errors"
	"time"
)

func ParseDate(s string) (*time.Time, error) {
	dateTime, err := time.ParseInLocation("2006-01-02", s, time.Local)

	if err != nil {
		return nil, err
	}

	return &dateTime, nil
}

func ParseTime(s string) (*time.Time, error) {
	dateTime, err := time.ParseInLocation("15:04", s, time.Local)
	if err != nil {
		return nil, err
	}

	return &dateTime, nil
}

func ParseDateTime(s string) (*time.Time, error) {
	dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)

	if err != nil {
		return nil, err
	}

	return &dateTime, nil
}

func ToAtString(in *time.Time) string {
	return in.Format("15:04 02.01.2006")
}

func ToDateTimeString(in *time.Time) string {
	return in.Format("2006-01-02 15:04:05")
}

func ToTimeString(in *time.Time) string {
	return in.Format("15:04:05")
}

func GetBookingDateTime(leadDays int32, slotTime *time.Time, reservationDate *time.Time) (*time.Time, error) {
	currentTime := time.Now()
	bookingDateTime := time.Date(reservationDate.Year(), reservationDate.Month(), reservationDate.Day(), slotTime.Hour(), slotTime.Minute(), slotTime.Second(), 0, time.Local).AddDate(0, 0, -int(leadDays))

	if bookingDateTime.Before(currentTime) {
		return nil, errors.New("slots for this reservation date have already opened - resy-cli can't help you here 😢")
	}

	return &bookingDateTime, nil
}
