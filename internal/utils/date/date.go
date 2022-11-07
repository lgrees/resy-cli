package date

import "time"

func ParseDate(s string) (*time.Time, error) {
	dateTime, err := time.ParseInLocation("2006-01-02", s, time.Local)

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
