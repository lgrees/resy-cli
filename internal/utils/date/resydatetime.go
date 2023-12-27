package date

import (
	"time"
)

type ResyDate struct {
	time.Time
	Format string
}

func NewResyDate(any interface{}, format string) *ResyDate {
	var t time.Time
	switch p := any.(type) {
	case time.Time:
		t = p
	case string:
		t, _ = time.Parse(time.DateOnly, p)
	}
	return &ResyDate{Time: t, Format: format}
}

func (d *ResyDate) UnmarshalJSON(b []byte) error {
	date, err := time.Parse(d.Format, string(b))
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}

func (d *ResyDate) MarshalJSON() ([]byte, error) {
	return []byte(d.String()), nil
}

func (d *ResyDate) String() string {
	return d.Time.Format(d.Format)
}
