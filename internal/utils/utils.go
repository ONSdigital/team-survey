package utils

import (
	"fmt"
	"regexp"
	"time"
)

// ConvertCFPostgresConnectionString takes a connection string from VCAP_SERVICES
func ConvertCFPostgresConnectionString(s string) string {
	var re = regexp.MustCompile(`postgres://(.*?):(.*)@(.*):([0-9]+)/(.*)`)
	result := re.FindAllStringSubmatch(s, -1)[0]
	match := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", result[3], result[4], result[1], result[5], result[2])
	return match
}

//ShortDateFromString parse shot date from string
func ShortDateFromString(ds string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", ds)
	if err != nil {
		return t, err
	}
	return t, nil
}

//CheckDataBoundariesStr checks is startdate <= enddate
func CheckDataBoundariesStr(startdate, enddate string) (bool, error) {

	tstart, err := ShortDateFromString(startdate)
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}
	tend, err := ShortDateFromString(enddate)
	if err != nil {
		return false, fmt.Errorf("cannot parse enddate: %v", err)
	}

	if tstart.After(tend) {
		return false, fmt.Errorf("invalid date range")
	}
	return true, err
}
