package date

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	switch p := any.(type) {
	case string:
		t, err := time.Parse(format, p)
		if err != nil {
			return nil, err
		}
		return NewResyDate(t, format)
	case time.Time:
		return &ResyDate{Time: p, FormatStr: format}, nil
	default:
		return nil, fmt.Errorf("unsupported type: %s", reflect.TypeOf(any))
	}
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
