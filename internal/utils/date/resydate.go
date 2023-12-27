package date

import (
	"encoding/json"
	"time"
)

const (
	AtFmt = "15:04 02.01.2006"
)

type ResyDate struct {
	time.Time
	FormatStr string
}

func NewResyDate(any interface{}, format string) (*ResyDate, error) {
	var t time.Time
	switch p := any.(type) {
	case time.Time:
		t = p
	case string:
		strT, err := time.Parse(format, p)
		if err != nil {
			return nil, err
		}
		t = strT
	}
	return &ResyDate{Time: t, FormatStr: format}, nil
}

func (d *ResyDate) UnmarshalJSON(b []byte) error {
	date, err := time.Parse(d.FormatStr, string(b))
	if err != nil {
		return err
	}
	d.Time = date
	return nil
}

func (d *ResyDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *ResyDate) String() string {
	return d.Time.Format(d.FormatStr)
}
