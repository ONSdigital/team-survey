package utils

import (
	"fmt"
	"testing"
)

const VCapServices string = `postgres://seilbmbd:examplepass@babar.elephantsql.com:5432/seilbmbd`
const ValidPostgresString string = `host=babar.elephantsql.com port=5432 user=seilbmbd dbname=seilbmbd password=examplepass`

func assertString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}

func assertErrorNotNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}
}

func TestConvertCFPostgresConnectionString(t *testing.T) {
	cString := ConvertCFPostgresConnectionString(VCapServices)

	if cString != ValidPostgresString {
		t.Errorf("Expected '%s', got '%s'", ValidPostgresString, cString)
	}
}

func TestShortDateFromString(t *testing.T) {
	d, err := ShortDateFromString("2019-04-01")

	assertErrorNotNil(t, err)
	assertString(t, "2019-04-01 00:00:00 +0000 UTC", fmt.Sprintf("%s", d))
}

func TestShortDateFromStringInvalid(t *testing.T) {
	_, err := ShortDateFromString("scpc")
	assertString(t, `parsing time "scpc" as "2006-01-02": cannot parse "scpc" as "2006"`, err.Error())
}

func TestCheckDataBoundariesStr(t *testing.T) {
	_, err := CheckDataBoundariesStr("2019-01-01", "2019-01-30")
	assertErrorNotNil(t, err)
}

func TestCheckDataBoundariesEndBeforeStart(t *testing.T) {
	_, err := CheckDataBoundariesStr("2019-02-01", "2019-01-30")
	assertString(t, `invalid date range`, err.Error())
}

func TestCheckDataBoundariesBadStart(t *testing.T) {
	_, err := CheckDataBoundariesStr("scpc", "2019-01-30")
	assertString(t, `cannot parse startdate: parsing time "scpc" as "2006-01-02": cannot parse "scpc" as "2006"`, err.Error())
}

func TestCheckDataBoundariesBadEnd(t *testing.T) {
	_, err := CheckDataBoundariesStr("2019-02-01", "scpc")
	assertString(t, `cannot parse enddate: parsing time "scpc" as "2006-01-02": cannot parse "scpc" as "2006"`, err.Error())
}
