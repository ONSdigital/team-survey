package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ONSdigital/team-survey/internal/stats"
	"github.com/ONSdigital/team-survey/internal/survey"
	"github.com/ONSdigital/team-survey/internal/utils"

	"github.com/Pallinder/go-randomdata"
	"github.com/gorilla/mux"
	"github.com/rs/xid"
)

// GetSurveyQuestionnaireHandler returns a survey questionnaire
func (a App) GetSurveyQuestionnaireHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	s, err := survey.Survey{}.GetByName(a.DB, vars["name"])

	if err != nil {
		respondWithErrorCheck(w, err, "record not found", http.StatusNotFound)
		return
	}

	builder := survey.YamlFileToSurveyJSBuilder{}
	builder.Survey = "assets/survey-components/survey.yml"
	builder.Pages = []string{
		"assets/survey-components/about/ons.yml",
		"assets/survey-components/culture/westrum.yml",
		"assets/survey-components/cohesion/lencioni.yml",
		"assets/survey-components/metrics/accelerate.yml",
	}
	sjs, err := survey.BuildSurveyFromYMLFiles(builder)

	if err != nil {
		log.Printf("%s", err)
	}

	sp := SitePage{
		Template: "web/templates/survey.html",
		Survey:   s,
		SurveyJS: template.JS(sjs),
	}
	renderTemplate(w, sp)
}

//GetAdminDashboardHandler is a route handler
func (a App) GetAdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	sp := SitePage{
		Template: "web/templates/admin/dashboard-team.html",
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}
	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}

//GetSurveyAdminDeleteHandler is a route handler
func (a App) GetSurveyAdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
	s, err := a.API.GetSurveyByName(getMuxVar(r, "survey"))

	if err != nil {
		respondWithErrorCheck(w, err, "record not found", http.StatusNotFound)
		return
	}

	sp := SitePage{
		Template: "web/templates/admin/survey-delete-confirm.html",
		Survey:   s,
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}

	if r.URL.Query().Get("confirm") == s.ShareCode {

		s, err = a.API.GetSurveyByNameAndShareCode(getMuxVar(r, "survey"), r.URL.Query().Get("confirm"))
		if err != nil || s.ID == 0 {
			respondWithErrorCheck(w, err, "record not found", http.StatusNotFound)
			return
		}

		a.DB.Delete(&s)

		sp = SitePage{
			Template: "web/templates/admin/survey-deleted.html",
			Survey:   s,
			AppUser:  authProvider.GetUserFromContext(r.Context()),
		}

	}

	renderTemplate(w, sp)
}

//GetTeamAdminDashboardHandler is a route handler
func (a App) GetTeamAdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var t survey.Team
	t.Name = getMuxVar(r, "team") // Do this to check for "All"

	if t.Name != "All" {
		t, err = survey.Team{}.GetByName(a.DB, t.Name)
		if err != nil {
			respondWithErrorCheck(w, err, "record not found", http.StatusNotFound)
			return
		}
	}

	sp := SitePage{
		Template: "web/templates/admin/dashboard-team.html",
		Team:     t,
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}
	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}

//GetSurveyAdminDashboardHandler is a route handler
func (a App) GetSurveyAdminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	var s survey.Survey
	var err error

	if r.URL.Query().Get("share_code") != "" {
		s, err = a.API.GetSurveyByNameAndShareCode(getMuxVar(r, "survey"), r.URL.Query().Get("share_code"))
	} else {
		s, err = a.API.GetSurveyByName(getMuxVar(r, "survey"))
	}

	if err != nil {
		respondWithErrorCheck(w, err, "record not found", http.StatusNotFound)
		return
	}
	sp := SitePage{
		Template: "web/templates/admin/dashboard-survey.html",
		Survey:   s,
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}
	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}

//GetCreateNewSurveyHandler is a route handler
func (a App) GetCreateNewSurveyHandler(w http.ResponseWriter, r *http.Request) {
	s := survey.Survey{
		Name: strings.ToLower(fmt.Sprintf("%s-%s-%d", randomdata.Noun(), randomdata.Noun(), randomdata.Number(1000, 9999))),
	}
	sp := SitePage{
		Template: "web/templates/admin/create-survey.html",
		Survey:   s,
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}

	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}

//GetPublicCreateNewSurveyHandler is a route handler
func (a App) GetPublicCreateNewSurveyHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("access_token") != os.Getenv("CREATE_SURVEY_ACCESS_TOKEN") {
		respondWithError(w, http.StatusBadRequest, "Invalid access token")
		return
	}

	s := survey.Survey{
		Name: strings.ToLower(fmt.Sprintf("%s-%s-%d", randomdata.Noun(), randomdata.Noun(), randomdata.Number(1000, 9999))),
	}
	sp := SitePage{
		Template: "web/templates/admin/public-create-survey.html",
		Survey:   s,
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}

	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}

//GetAllSurveysHandler is a route handler
func (a App) GetAllSurveysHandler(w http.ResponseWriter, r *http.Request) {
	s, _ := survey.Survey{}.GetAll(a.DB)
	sp := SitePage{
		Template: "web/templates/admin/list-all-surveys.html",
		Surveys:  s,
		AppUser:  authProvider.GetUserFromContext(r.Context()),
	}
	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}

//HomePageHandler is a route handler
func (a App) HomePageHandler(w http.ResponseWriter, r *http.Request) {
	sp := SitePage{
		Template: "web/templates/index.html",
	}
	renderTemplate(w, sp)
}

// GetTeamsHandler gets a list of all teams
func (a App) GetTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := survey.Team{}.GetAll(a.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, teams)
}

// PostSurveyHandler process new surveys
func (a App) PostSurveyHandler(w http.ResponseWriter, r *http.Request) {
	var s survey.Survey

	err := json.NewDecoder(r.Body).Decode(&s)

	if err != nil {
		// Replace with a more sensible error message
		if err == io.EOF {
			err = fmt.Errorf("empty request body")
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if s, err = a.API.SaveSurvey(s); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, s)
}

// GetStatsHandler returns survey stats
func (a App) GetStatsHandler(w http.ResponseWriter, r *http.Request) {
	var s survey.Survey
	var err error

	if r.URL.Query().Get("share_code") != "" {
		s, err = a.API.GetSurveyByNameAndShareCode(getMuxVar(r, "survey"), r.URL.Query().Get("share_code"))
	} else {
		s, err = a.API.GetSurveyByName(getMuxVar(r, "survey"))
	}

	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	if s.ShareCode == "" {
		s.ShareCode = xid.New().String()
		_, _ = s.Save(a.DB)
	}

	var calculatedResults stats.CalculatedResult

	if calculatedResults, err = a.API.GetStats(s); err != nil {
		respondWithError(w, http.StatusFound, err.Error())
		return
	}

	calculatedResults.ShareCode = s.ShareCode

	respondWithJSON(w, http.StatusOK, calculatedResults)
}

// GetAllStatsHandler returns all survey results
func (a App) GetAllStatsHandler(w http.ResponseWriter, r *http.Request) {
	var s survey.Survey
	var err error

	dateStart := time.Now().Add(-720 * time.Hour).Local().Format("2006-01-02")
	dateEnd := time.Now().Local().Format("2006-01-02")

	v := r.URL.Query()

	if v.Get("dateStart") != "" {
		dateStart = v.Get("dateStart")
	}
	if v.Get("dateEnd") != "" {
		dateEnd = v.Get("dateEnd")
	}

	if _, err := utils.CheckDataBoundariesStr(dateStart, dateEnd); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	searchParams := survey.ResultSearchParameters{
		DateStart: dateStart,
		DateEnd:   dateEnd,
		Team:      v.Get("team"),
	}

	s, _ = survey.Survey{}.GetAllResults(a.DB, searchParams)

	s.Team.Name = searchParams.Team

	var calculatedResults stats.CalculatedResult

	if calculatedResults, err = a.API.GetStats(s); err != nil {
		log.Println(err)
	}

	calculatedResults.DateStart = dateStart
	calculatedResults.DateEnd = dateEnd

	respondWithJSON(w, http.StatusOK, calculatedResults)
}

//GetAllUsersHandler is a route handler
func (a App) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	u, _ := authProvider.GetAllUsers()
	sp := SitePage{
		Template: "web/templates/admin/list-all-users.html",
		AppUser:  authProvider.GetUserFromContext(r.Context()),
		Users:    u,
	}
	getSiteSections(&sp, r)
	renderTemplate(w, sp)
}
