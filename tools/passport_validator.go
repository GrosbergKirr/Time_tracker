package tools

import (
	"regexp"
)

func ValidatePassport(passport string) (bool, error) {
	layout := `^\d{4} \d{6}$`
	re, err := regexp.Compile(layout)
	if err != nil {
		return false, err
	}
	return re.MatchString(passport), nil

}
