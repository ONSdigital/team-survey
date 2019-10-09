package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/ONSdigital/team-survey/internal/survey"
)

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func checkErrorMessage(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	var m map[string]string
	json.Unmarshal(r.Body.Bytes(), &m)
	assertString(t, expected, m["error"])
}

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

func mockValidCFEnvironment() {
	os.Setenv("VCAP_APPLICATION", `{"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`)
	os.Setenv("HOME", "/home/vcap/app")
	os.Setenv("MEMORY_LIMIT", "512m")
	os.Setenv("PWD", "/home/vcap")
	os.Setenv("TMPDIR", "/home/vcap/tmp")
	os.Setenv("USER", "vcap")
	os.Setenv("VCAP_SERVICES", `{"elephantsql":[{"name":"elephantsql-dev-c6c60","label":"elephantsql-dev","tags":["New Product","relational","Data Store","postgresql"],"plan":"turtle","credentials":{"uri":"postgres://seilbmbd:PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn@foo.elephantsql.com:5432/seilbmbd"}}],"sendgrid":[{"name":"mysendgrid","label":"sendgrid","tags":["smtp","Email"],"plan":"free","credentials":{"hostname":"smtp.sendgrid.net","username":"QvsXMbJ3rK","password":"HCHMOYluTv"}}],"nfs":[{"credentials":{},"label":"nfs","name":"nfs","plan":"Existing","tags":["nfs"],"volume_mounts":[{"container_dir":"/testpath","device_type":"shared","mode":"rw"}]}]}`)
}

func mockValidCFEnvironmentNoElephantSQL() {
	os.Setenv("VCAP_APPLICATION", `{"instance_id":"451f045fd16427bb99c895a2649b7b2a","application_id":"abcabc123123defdef456456","cf_api": "https://api.system_domain.com","instance_index":0,"host":"0.0.0.0","port":61857,"started_at":"2013-08-12 00:05:29 +0000","started_at_timestamp":1376265929,"start":"2013-08-12 00:05:29 +0000","state_timestamp":1376265929,"limits":{"mem":512,"disk":1024,"fds":16384},"application_version":"c1063c1c-40b9-434e-a797-db240b587d32","application_name":"styx-james","application_uris":["styx-james.a1-app.cf-app.com"],"version":"c1063c1c-40b9-434e-a797-db240b587d32","name":"styx-james","space_id":"3e0c28c5-6d9c-436b-b9ee-1f4326e54d05","space_name":"jdk","uris":["styx-james.a1-app.cf-app.com"],"users":null}`)
	os.Setenv("HOME", "/home/vcap/app")
	os.Setenv("MEMORY_LIMIT", "512m")
	os.Setenv("PWD", "/home/vcap")
	os.Setenv("TMPDIR", "/home/vcap/tmp")
	os.Setenv("USER", "vcap")
	os.Setenv("VCAP_SERVICES", `{}`)
}

func teardownCFEnvironment() {
	os.Unsetenv("VCAP_APPLICATION")
	os.Unsetenv("HOME")
	os.Unsetenv("MEMORY_LIMIT")
	os.Unsetenv("PWD")
	os.Unsetenv("TMPDIR")
	os.Unsetenv("USER")
	os.Unsetenv("VCAP_SERVICES")
}

func deleteTestDB() {
	if _, err := os.Stat("test.db"); err == nil {
		os.Remove("test.db")
	}
}

func init() {
	deleteTestDB()
}

func TestMain(m *testing.M) {
	app.IsTest = true
	dbct, _ := app.GetDBConnectionType()
	config = AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	app.Initialize(dbct, config)
	code := m.Run()
	os.Exit(code)
}

func TestInitializeError(t *testing.T) {
	a := App{}
	dbct := DBConnectionCredentials{
		Dialect:          "invalid",
		ConnectionString: "whoknows",
	}
	config = AuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	}
	err := a.Initialize(dbct, config)

	assertString(t, `sql: unknown driver "invalid" (forgotten import?)`, err.Error())
}

func TestGetDBConnectionTypeDefault(t *testing.T) {
	d, err := app.GetDBConnectionType()

	assertErrorNotNil(t, err)
	assertString(t, `sqlite3`, d.Dialect)
	assertString(t, `test.db`, d.ConnectionString)
}

func TestGetDBConnectionTypeCFInvalidVcap(t *testing.T) {
	defer teardownCFEnvironment()

	os.Setenv("VCAP_APPLICATION", "{}")

	d, err := app.GetDBConnectionType()

	assertString(t, `unexpected end of JSON input`, err.Error())
	assertString(t, "", d.Dialect)
	assertString(t, "", d.ConnectionString)
}

func TestGetDBConnectionTypeCFValidVcap(t *testing.T) {
	defer teardownCFEnvironment()

	mockValidCFEnvironment()
	d, err := app.GetDBConnectionType()

	assertErrorNotNil(t, err)
	assertString(t, `postgres`, d.Dialect)
	assertString(t, "host=foo.elephantsql.com port=5432 user=seilbmbd dbname=seilbmbd password=PHxTPJSbkcDakfK4cYwXHiIX9Q8p5Bxn", d.ConnectionString)

}

func TestGetDBConnectionTypeCFValidVcapNoElephantSQL(t *testing.T) {
	defer teardownCFEnvironment()

	mockValidCFEnvironmentNoElephantSQL()
	d, err := app.GetDBConnectionType()

	assertErrorNotNil(t, err)
	assertString(t, `sqlite3`, d.Dialect)
	assertString(t, `test.db`, d.ConnectionString)
}

func TestHomePageHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetSurveyQuestionnaireHandlerInvalid(t *testing.T) {
	req, _ := http.NewRequest("GET", "/survey/sausage-and-mash/", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetSurveyQuestionnaireHandlerValid(t *testing.T) {
	var s survey.Survey
	s.Name = "this-is-my-new-survey"
	s.Team.Name = "the-testers"
	s.Save(app.DB)
	req, _ := http.NewRequest("GET", "/survey/this-is-my-new-survey/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetAdminDashboardHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/dashboard/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTeamAdminDashboardHandlerInvalid(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/dashboard/team/bad-team/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetTeamAdminDashboardHandlerValid(t *testing.T) {
	var s survey.Survey
	req, _ := http.NewRequest("GET", "/admin/dashboard/team/the-testers-rule/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	s.Name = "this-is-my-new-survey-again"
	s.Team.Name = "the-testers-rule"
	s.Save(app.DB)

	req, _ = http.NewRequest("GET", "/admin/dashboard/team/the-testers-rule/", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTeamAdminDashboardHandlerAll(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/dashboard/team/All/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetSurveyAdminDashboardHandlerInvalid(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/dashboard/survey/bad-survey/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetSurveyAdminDashboardHandlerValid(t *testing.T) {
	var s survey.Survey
	req, _ := http.NewRequest("GET", "/admin/dashboard/survey/spartacus/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	s.Name = "spartacus"
	s.Team.Name = "the-testers-rule"
	s.Save(app.DB)

	req, _ = http.NewRequest("GET", "/admin/dashboard/survey/spartacus/", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetSurveyAdminDeleteConfirmHandlerValid(t *testing.T) {
	var s survey.Survey
	req, _ := http.NewRequest("GET", "/admin/dashboard/survey/pluto-isnt-working/delete", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	s.Name = "pluto-isnt-working"
	s.Team.Name = "the-testers-rule"
	s.Save(app.DB)

	req, _ = http.NewRequest("GET", "/admin/dashboard/survey/pluto-isnt-working/delete", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetSurveyAdminDeleteHandlerValid(t *testing.T) {
	var s survey.Survey
	req, _ := http.NewRequest("GET", "/admin/dashboard/survey/pluto-isnt-working-again/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	s.Name = "pluto-isnt-working-again"
	s.Team.Name = "the-testers-rule"
	s.ShareCode = "sc123"
	s.Save(app.DB)

	req, _ = http.NewRequest("GET", "/admin/dashboard/survey/pluto-isnt-working-again/", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/admin/dashboard/survey/pluto-isnt-working-again/delete?confirm=%s", "333"), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/admin/dashboard/survey/pluto-isnt-working-again/", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/admin/dashboard/survey/pluto-isnt-working-again/delete?confirm=%s", s.ShareCode), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/admin/dashboard/survey/pluto-isnt-working-again/", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestGetCreateNewSurveyHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/survey/new/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetPublicCreateNewSurveyHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/public/survey/new/?access_token=%s", os.Getenv("CREATE_SURVEY_ACCESS_TOKEN")), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetAllSurveysHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/admin/survey/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTeamsHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/teams/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetTeamsHandlerError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/teams/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestPostSurveyHandlerEmptyBody(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/survey/", strings.NewReader(""))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkErrorMessage(t, response, "empty request body")
}

func TestPostSurveyHandlerEmptyJSON(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/survey/", strings.NewReader("{}"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkErrorMessage(t, response, "survey has no name")
}

func TestPostSurveyHandlerInvalidTeam(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/survey/", strings.NewReader(`{"name":"fry-bagel"}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkErrorMessage(t, response, "team not defined")
}

func TestPostSurveyHandler(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/v1/survey/", strings.NewReader(`{"name":"fry-bagel","team":{"name":"onion-soup"}}`))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetStatsHandlerInvalid(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/survey/chocolate-icecream/stats/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
	checkErrorMessage(t, response, "record not found")
}

func TestGetStatsHandlerNoResults(t *testing.T) {
	var s survey.Survey
	s.Name = "who-loves-nachos"
	s.Team.Name = "we-all-do"
	s.Save(app.DB)

	req, _ := http.NewRequest("GET", "/api/v1/survey/who-loves-nachos/stats/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusFound, response.Code)
	checkErrorMessage(t, response, "no results to calculate")
}

func TestGetStatsHandler(t *testing.T) {
	var s survey.Survey
	s.Name = "who-loves-trifle"
	s.Team.Name = "only-some-people"
	s.Results = append(s.Results, survey.Result{})
	s.Save(app.DB)

	req, _ := http.NewRequest("GET", "/api/v1/survey/who-loves-trifle/stats/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetAllStatsHandler(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/stats/all/", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetAllStatsHandlerDateRangeInvalidRange(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/stats/all/?dateStart=2019-02-02&dateEnd=2019-02-01", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkErrorMessage(t, response, "invalid date range")
}

func TestGetAllStatsHandlerDateRangeInvalidDateFormat(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/stats/all/?dateStart=aa-bb-cc&dateEnd=2019-02-01", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
	checkErrorMessage(t, response, `cannot parse startdate: parsing time "aa-bb-cc" as "2006-01-02": cannot parse "aa-bb-cc" as "2006"`)
}

func TestGetAllStatsHandlerDateRange(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/v1/stats/all/?dateStart=2019-02-01&dateEnd=2019-02-28", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestCleanup(t *testing.T) {
	deleteTestDB()
}
