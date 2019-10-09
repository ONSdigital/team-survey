package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/ONSdigital/team-survey/internal/api"
	"github.com/ONSdigital/team-survey/internal/auth"
	"github.com/ONSdigital/team-survey/internal/survey"
	"github.com/ONSdigital/team-survey/internal/utils"

	cfenv "github.com/cloudfoundry-community/go-cfenv"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/google"
	"github.com/dghubble/sessions"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/urfave/negroni"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
)

const (
	sessionName    = "team-survey"
	sessionSecret  = "team-survey"
	sessionUserKey = "team-survey"
)

// App is the wrapper for everything
type App struct {
	Router *mux.Router
	DB     *gorm.DB
	API    api.API
	IsTest bool
}

// SitePage holds template contents
type SitePage struct {
	Template     string
	TemplateName string
	Section      string
	SubSection   string
	Survey       survey.Survey
	Surveys      []survey.Survey
	Team         survey.Team
	Teams        []survey.Team
	SurveyJS     template.JS
	AppUser      auth.User
	Users        []auth.User
}

// DBConnectionCredentials are for determining the database type
type DBConnectionCredentials struct {
	Dialect          string
	ConnectionString string
}

// AuthConfig configures the main ServeMux.
type AuthConfig struct {
	ClientID     string
	ClientSecret string
}

var authProvider auth.Provider

// Initialize sets up the app
func (a *App) Initialize(d DBConnectionCredentials, config AuthConfig) error {
	db, err := gorm.Open(d.Dialect, d.ConnectionString)

	if err != nil {
		return err
	}

	a.DB = db
	a.API.DB = a.DB

	authProvider = auth.CreateProvider(db, sessions.NewCookieStore([]byte(sessionSecret), nil), sessionName, sessionSecret, sessionUserKey)

	a.Router = mux.NewRouter().StrictSlash(true)

	if err := survey.AutoMigrateModels(app.DB); err != nil {
		return err
	}

	if err := auth.AutoMigrateModels(app.DB); err != nil {
		return err
	}

	a.Router.HandleFunc("/api/v1/teams/", http.HandlerFunc(a.GetTeamsHandler)).Methods("GET")
	a.Router.HandleFunc("/api/v1/stats/", http.HandlerFunc(a.GetAllStatsHandler)).Methods("GET")
	a.Router.HandleFunc("/api/v1/stats/all/", http.HandlerFunc(a.GetAllStatsHandler)).Methods("GET")
	a.Router.HandleFunc("/api/v1/survey/{survey}/stats/", http.HandlerFunc(a.GetStatsHandler)).Methods("GET")
	a.Router.HandleFunc("/api/v1/survey/", a.PostSurveyHandler).Methods("POST")

	a.Router.PathPrefix("/static/").Handler(http.FileServer(http.Dir("assets")))

	a.Router.HandleFunc("/survey/{name}/", http.HandlerFunc(a.GetSurveyQuestionnaireHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/dashboard/", http.HandlerFunc(a.GetAdminDashboardHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/dashboard/team/{team}/", http.HandlerFunc(a.GetTeamAdminDashboardHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/dashboard/survey/{survey}/", http.HandlerFunc(a.GetSurveyAdminDashboardHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/dashboard/survey/{survey}/delete", http.HandlerFunc(a.GetSurveyAdminDeleteHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/survey/new/", http.HandlerFunc(a.GetCreateNewSurveyHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/survey/", http.HandlerFunc(a.GetAllSurveysHandler)).Methods("GET")
	a.Router.HandleFunc("/admin/users/", http.HandlerFunc(a.GetAllUsersHandler)).Methods("GET")

	a.Router.HandleFunc("/public/dashboard/survey/{survey}/", http.HandlerFunc(a.GetSurveyAdminDashboardHandler)).Methods("GET")
	a.Router.HandleFunc("/public/survey/new/", http.HandlerFunc(a.GetPublicCreateNewSurveyHandler)).Methods("GET")

	a.Router.HandleFunc("/logout", authProvider.LogoutHandler)

	oauth2Config := &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  fmt.Sprintf("%s/google/callback", os.Getenv("APP_URL")),
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	stateConfig := gologin.DebugOnlyCookieConfig
	a.Router.Handle("/google/login", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	a.Router.Handle("/google/callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, authProvider.IssueSession(), nil)))

	a.Router.HandleFunc("/", a.HomePageHandler).Methods("GET")

	return nil
}

// Run makes this thing fly
func (a *App) Run(port string) {
	n := negroni.Classic()
	sirMuxalot := http.NewServeMux()
	sirMuxalot.Handle("/", app.Router)
	if os.Getenv("DISABLE_AUTH") != "True" {
		log.Println("AUTH IS ENABLED")
		sirMuxalot.Handle("/admin/", negroni.New(
			negroni.HandlerFunc(DashboardMiddleware),
			negroni.Wrap(app.Router),
		))
		sirMuxalot.Handle("/api/v1/survey/", negroni.New(
			negroni.HandlerFunc(ShareCodeOrSSOMiddleware),
			negroni.Wrap(app.Router),
		))
		sirMuxalot.Handle("/public/dashboard/", negroni.New(
			negroni.HandlerFunc(ShareCodeMiddleware),
			negroni.Wrap(app.Router),
		))
		sirMuxalot.Handle("/api/v1/", negroni.New(
			negroni.HandlerFunc(DashboardMiddleware),
			negroni.Wrap(app.Router),
		))
	}

	n.UseHandler(sirMuxalot)

	srv := &http.Server{
		Handler:      authProvider.AddAppUserContext(n),
		Addr:         fmt.Sprintf(":%s", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("App ready... http://localhost:%s", port)
	log.Fatal(srv.ListenAndServe())
}

func responseWithStatusCode(w http.ResponseWriter, code int) {
	w.WriteHeader(code)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithErrorCheck(w http.ResponseWriter, err error, message string, code int) {
	if err != nil {
		if err.Error() == message {
			responseWithStatusCode(w, code)
		}
		responseWithStatusCode(w, http.StatusBadRequest)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

func getSiteSections(p *SitePage, r *http.Request) string {
	var matches []string
	var re = regexp.MustCompile(`^(\/.*?\/[a-z]+)`)
	m := re.FindAllStringSubmatch(r.RequestURI, -1)

	if len(m) > 0 {
		matches = m[0]
	}

	if len(matches) > 1 {
		p.Section = matches[1]
	}

	p.SubSection = r.RequestURI
	return ""
}

func renderTemplate(w http.ResponseWriter, sp SitePage) {
	basepath := filepath.Dir(sp.Template)
	basetemplate := fmt.Sprintf("%s/%s", basepath, "_template.html")

	sp.TemplateName = filepath.Base(sp.Template)

	templates, err := template.ParseFiles(
		sp.Template,
	)

	if err != nil {
		log.Printf("%s", err)
	}

	if _, err := os.Stat(basetemplate); !os.IsNotExist(err) {
		templates, err = template.ParseFiles(
			basetemplate,
			sp.Template,
		)
		if err != nil {
			log.Printf("%s", err)
		}
	}
	if err := templates.Execute(w, sp); err != nil {
		log.Printf("%s", err)
	}

}

func getMuxVar(r *http.Request, v string) string {
	vars := mux.Vars(r)
	return vars[v]
}

// GetDBConnectionType will give a best effort to determine available databases
func (a *App) GetDBConnectionType() (DBConnectionCredentials, error) {
	dbType := "sqlite3"
	dbConnectionString := "team-survey.db"

	if a.IsTest {
		dbConnectionString = "test.db"
	}

	if cfenv.IsRunningOnCF() {
		app, err := cfenv.Current()

		if err != nil {
			return DBConnectionCredentials{}, err
		}

		if len(app.Services["elephantsql"]) > 0 {
			dbType = "postgres"
			dbConnectionString = utils.ConvertCFPostgresConnectionString(fmt.Sprintf("%s", app.Services["elephantsql"][0].Credentials["uri"]))
		}
	}

	return DBConnectionCredentials{Dialect: dbType, ConnectionString: dbConnectionString}, nil
}
