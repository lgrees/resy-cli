package surveyHelpers

import (
	"errors"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/lgrees/resy-cli/internal/utils/date"
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

	_, err := date.ParseDate(str)
	if err != nil {
		return errors.New("input must be a valid date (YYYY-MM-DD)")
	}

	return nil
}

func DateTimeValidator(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return errors.New("Input must be a string.")
	}

	_, err := date.ParseDateTime(str)
	if err != nil {
		return errors.New("Input must be a valid date time (YYYY-MM-DD HH:MM:SS)")
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
