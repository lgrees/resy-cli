package surveyHelpers

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/bcillie/resy-cli/internal/utils/date"
)

func CreateRegexValidator(s, e string) survey.Validator {
	return func(val interface{}) error {
		str, ok := val.(string)
		if !ok {
			return errors.New("input must be a string")
		}

		match, _ := regexp.MatchString(s, str)
		if !match {
			return errors.New(e)
		}
		return nil
	}
}

func DateValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return errors.New("input must be a string")
	}

	t, err := date.NewResyDate(str, time.DateOnly)
	if err != nil {
		return errors.New("input must be a valid date (YYYY-MM-DD)")
	}

	if !t.After(time.Now().Local()) {
		return errors.New("the date selected should be in the future")
	}

	return nil
}

func TimeValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return errors.New("input must be a string")
	}

	_, err := date.ParseTime(str)

	if err != nil {
		return errors.New("input must be a valid time (HH:MM)")
	}

	return nil
}

func TimesValidator(val interface{}) error {
	in, ok := val.(string)

	if !ok {
		return errors.New("input must be a string")
	}

	arr := strings.Split(in, "\n")

	for _, v := range arr {
		err := TimeValidator(v)
		if err != nil {
			return err
		}

	}

	return nil
}

func VenueValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return errors.New("input must be a string")
	}

	arr := strings.Split(str, " | ")
	if len(arr) < 5 {
		return errors.New("please tab to search and select a venue")
	}
	return nil
}

func TransformLowerCase(in interface{}) (out interface{}) {
	str, ok := in.(string)
	if !ok {
		return ""
	}

	return strings.ToLower(str)
}
