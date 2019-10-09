package survey

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var testDBFileName string
var surveyCount int

func init() {
	testDBFileName = "test.db"
	deleteTestDB()
	var err error

	if db, err = gorm.Open("sqlite3", testDBFileName); err != nil {
		panic("failed to connect database")
	}
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

func checkErrorMessage(t *testing.T, r *httptest.ResponseRecorder, expected string) {
	var m map[string]string
	json.Unmarshal(r.Body.Bytes(), &m)
	assertString(t, expected, m["error"])
}

func deleteTestDB() {
	if _, err := os.Stat(testDBFileName); err == nil {
		os.Remove(testDBFileName)
		fmt.Println("=== INFO  Test DB Deleted")
	}
}

func createNewSurvey() Survey {
	results := []Result{}

	results = append(results, Result{
		CurrentRole: "Tester",
	})

	surveyCount++

	return Survey{
		Name: fmt.Sprintf("Test Survey %d", surveyCount),
		Team: Team{
			Name: "Amazing Team",
		},
		Results: results,
	}
}

func TestAutoMigrateModels(t *testing.T) {
	err := AutoMigrateModels(db)
	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}
}

func TestSurveyGet(t *testing.T) {
	savedSurvey, _ := createNewSurvey().Save(db)

	survey, err := Survey{}.Get(db, savedSurvey.ID)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if savedSurvey.ID == 0 {
		t.Errorf("Expected ID not to be 0, but got %d", savedSurvey.ID)
	}

	if survey.ID != savedSurvey.ID {
		t.Errorf("Expected ID to be 1, but got %d", survey.ID)
	}
}

func TestSurveyGetNonExistent(t *testing.T) {
	savedSurvey, _ := createNewSurvey().Save(db)

	survey, err := Survey{}.Get(db, savedSurvey.ID+1)

	if err.Error() != "record not found" {
		t.Errorf("Expected record not found, but got %s", err)
	}

	if survey.ID != 0 {
		t.Errorf("Expected ID to be 0, but got %d", survey.ID)
	}
}

func TestSurveyGetByName(t *testing.T) {
	savedSurvey, _ := createNewSurvey().Save(db)

	survey, err := Survey{}.GetByName(db, savedSurvey.Name)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if savedSurvey.ID == 0 {
		t.Errorf("Expected ID not to be 0, but got %d", savedSurvey.ID)
	}

	if survey.Name != savedSurvey.Name {
		t.Errorf("Expected Name to be %s, but got %s", survey.Name, savedSurvey.Name)
	}
}

func TestSurveyGetByNameNonExisting(t *testing.T) {
	_, err := Survey{}.GetByName(db, "CabbageBacon")

	if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found', but got %s", err)
	}
}

func TestSurveySaveCreateNew(t *testing.T) {
	savedSurvey, err := createNewSurvey().Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if savedSurvey.ID == 0 {
		t.Errorf("Expected ID not to be 0, but got %d", savedSurvey.ID)
	}
}

func TestSurveySaveUpdate(t *testing.T) {
	savedSurvey, err := createNewSurvey().Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	savedSurvey.Name = "Test Survey 1 Updated"
	_, err = savedSurvey.Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	// Fetch the survey from the database after it's been "updated" to be sure
	fetchedSurvey, err := Survey{}.Get(db, savedSurvey.ID)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if fetchedSurvey.ID != savedSurvey.ID {
		t.Errorf("Expected ID to be %d, but got %d", savedSurvey.ID, fetchedSurvey.ID)
	}

	if fetchedSurvey.Name != "test-survey-1-updated" {
		t.Errorf("Expected 'test-survey-1-updated', but got %s", fetchedSurvey.Name)
	}
}

func TestSurveySaveWithResultsNew(t *testing.T) {
	savedSurvey, err := createNewSurvey().Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if savedSurvey.Results[0].ID == 0 {
		t.Errorf("Expected not 0, but got %d", savedSurvey.Results[0].ID)
	}

	fetchedSurvey, _ := Survey{}.Get(db, savedSurvey.ID)

	if len(fetchedSurvey.Results) != len(savedSurvey.Results) {
		t.Errorf("Expected %d result(s), but got %d", len(savedSurvey.Results), len(fetchedSurvey.Results))
	}

	if fetchedSurvey.Results[0].CurrentRole != "Tester" {
		t.Errorf("Expected 'Tester', but got %s", fetchedSurvey.Results[0].CurrentRole)
	}
}

func TestSurveySaveWithResultsMultiple(t *testing.T) {
	resultsToSave := []Result{}

	resultsToSave = append(resultsToSave, Result{
		CurrentRole: "QA",
	})

	resultsToSave = append(resultsToSave, Result{
		CurrentRole: "Engineer",
	})
	resultsToSave = append(resultsToSave, Result{
		CurrentRole: "Swimmer",
	})

	surveyToSave := Survey{
		Name:    "Test Survey 205",
		Results: resultsToSave,
		Team: Team{
			Name: "Amazing Team",
		},
	}

	savedSurvey, err := surveyToSave.Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	fetchedSurvey, _ := Survey{}.Get(db, savedSurvey.ID)

	if len(fetchedSurvey.Results) != len(savedSurvey.Results) {
		t.Errorf("Expected %d result(s), but got %d", len(savedSurvey.Results), len(fetchedSurvey.Results))
	}

	if fetchedSurvey.Results[0].CurrentRole != "QA" {
		t.Errorf("Expected 'QA', but got %s", fetchedSurvey.Results[0].CurrentRole)
	}

	if fetchedSurvey.Results[1].CurrentRole != "Engineer" {
		t.Errorf("Expected 'Engineer', but got %s", fetchedSurvey.Results[1].CurrentRole)
	}

	if fetchedSurvey.Results[2].CurrentRole != "Swimmer" {
		t.Errorf("Expected 'Swimmer', but got %s", fetchedSurvey.Results[2].CurrentRole)
	}
}

func TestSurveySaveWithNoName(t *testing.T) {
	s := createNewSurvey()
	s.Name = ""

	_, err := s.Save(db)

	if err.Error() != "survey has no name" {
		t.Errorf("Expected 'survey has no name, but got %s", err)
	}
}

func TestSurveySaveWithNoTeam(t *testing.T) {
	s := createNewSurvey()
	s.Team = Team{}

	_, err := s.Save(db)

	if err.Error() != "team not defined" {
		t.Errorf("Expected 'team not defined', but got %s", err)
	}
}

func TestSurveySaveWithTeamNew(t *testing.T) {

	savedSurvey, err := createNewSurvey().Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	fetchedSurvey, err := Survey{}.Get(db, savedSurvey.ID)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if fetchedSurvey.Team.ID == 0 {
		t.Errorf("Expected not 0, but got %d", fetchedSurvey.Team.ID)
	}

	if fetchedSurvey.Team.Name != "Amazing Team" {
		t.Errorf("Expected Amazing Team, but got %s", fetchedSurvey.Team.Name)
	}
}

func TestSurveySaveWithTeamUpdate(t *testing.T) {
	savedSurvey, err := createNewSurvey().Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if savedSurvey.Team.ID == 0 {
		t.Errorf("Expected not 0, but got %d", savedSurvey.Team.ID)
	}

	if savedSurvey.Team.Name != "Amazing Team" {
		t.Errorf("Expected Amazing Team, but got %s", savedSurvey.Team.Name)
	}

	originalTeam := savedSurvey.Team

	savedSurvey.Team = Team{
		Name: "Bad Team",
	}

	savedSurvey, err = savedSurvey.Save(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	fetchedSurvey, _ := Survey{}.Get(db, savedSurvey.ID)

	if savedSurvey.Team.Name != fetchedSurvey.Team.Name {
		t.Errorf("Expected %s, but got %s", savedSurvey.Team.Name, fetchedSurvey.Team.Name)
	}

	if savedSurvey.Team.ID == originalTeam.ID {
		t.Errorf("Expected %d, but got %d", savedSurvey.Team.ID, originalTeam.ID)
	}

	savedSurvey.Team = Team{
		Name: "Amazing Team",
	}

	savedSurvey, _ = savedSurvey.Save(db)

	if savedSurvey.Team.Name != "Amazing Team" {
		t.Errorf("Expected Amazing Team, but got %s", savedSurvey.Team.Name)
	}

	if savedSurvey.Team.ID != originalTeam.ID {
		t.Errorf("Expected %d, but got %d", savedSurvey.Team.ID, originalTeam.ID)
	}
}

func TestGetAllTeams(t *testing.T) {
	teams, err := Team{}.GetAll(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if len(teams) == 0 {
		t.Errorf("Expected not 0, but got %d", len(teams))
	}
}

func TestGetTeamByNameInvalid(t *testing.T) {
	_, err := Team{}.GetByName(db, "the-brisketeers")

	if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found', but got %s", err)
	}
}

func TestGetTeamByName(t *testing.T) {
	s := createNewSurvey()
	s.Team.Name = "the-brisketeers"
	s.Save(db)

	team, err := Team{}.GetByName(db, "the-brisketeers")

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if team.Name != "the-brisketeers" {
		t.Errorf("Expected 'the-brisketeers', but got %s", team.Name)
	}
}

func TestGetAll(t *testing.T) {
	surveys, err := Survey{}.GetAll(db)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if len(surveys) == 0 {
		t.Errorf("Expected not 0, but got %d", len(surveys))
	}
}

func TestGetAllResults(t *testing.T) {
	dateStart := time.Now().Add(-720 * time.Hour).Local().Format("2006-01-02")
	dateEnd := time.Now().Local().Format("2006-01-02")
	searchParams := ResultSearchParameters{
		DateStart: dateStart,
		DateEnd:   dateEnd,
	}
	survey, err := Survey{}.GetAllResults(db, searchParams)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if len(survey.Results) == 0 {
		t.Errorf("Expected not 0, but got %d", len(survey.Results))
	}

	searchParams.Team = "THISTEAMDOESNOTEXIST"

	survey, err = Survey{}.GetAllResults(db, searchParams)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if len(survey.Results) != 0 {
		t.Errorf("Expected 0, but got %d", len(survey.Results))
	}
}

func TestPopulateSurveyJSFromYML(t *testing.T) {
	var sjs JSSurvey
	yamlString := `---
completedHtml: my test survey`
	assertErrorNotNil(t, PopulateSurveyJSFromYML(&sjs, yamlString))
	assertString(t, "my test survey", string(sjs.CompletedHTML))

	yamlString = `---
pages:
- name: "page1"`
	assertErrorNotNil(t, PopulateSurveyJSFromYML(&sjs, yamlString))
	if len(sjs.Pages) != 1 {
		t.Errorf("Expected 1, got %d", len(sjs.Pages))
	}
	assertString(t, "page1", sjs.Pages[0].Name)

	yamlString = `---
pages:
- name: "page2"`
	assertErrorNotNil(t, PopulateSurveyJSFromYML(&sjs, yamlString))
	if len(sjs.Pages) != 1 {
		t.Errorf("Expected 1, got %d", len(sjs.Pages))
	}
	assertString(t, "page2", sjs.Pages[0].Name)
}

func TestGenerateSurveyJSON(t *testing.T) {
	var sjs = JSSurvey{
		CompletedHTML: template.HTML("Hello"),
		Pages: []JSPage{
			JSPage{
				Name: "page1",
			},
		},
	}
	j, err := GenerateSurveyJSON(sjs)
	assertErrorNotNil(t, err)
	assertString(t, `{"completedHtml":"Hello","pages":[{"name":"page1","title":"","elements":null}],"showQuestionNumbers":""}`, j)
}

func TestBuildSurveyFromYMLFiles(t *testing.T) {
	builder := YamlFileToSurveyJSBuilder{}
	builder.Survey = "assets/survey-components/survey.yml"
	builder.Pages = []string{
		"../assets/survey-components/about/ons.yml",
		"../assets/survey-components/culture/westrum.yml",
		"../assets/survey-components/cohesion/lencioni.yml",
		"../assets/survey-components/metrics/accelerate.yml",
	}

	_, err := BuildSurveyFromYMLFiles(builder)
	assertErrorNotNil(t, err)
}

func TestBuildSurveyFromYMLFilesInvalid(t *testing.T) {
	builder := YamlFileToSurveyJSBuilder{}
	builder.Survey = "assets/survey-components/survey.yml"
	builder.Pages = []string{
		"../assets/survey-components/about/does-not-exist.yml",
	}

	_, err := BuildSurveyFromYMLFiles(builder)
	assertErrorNotNil(t, err)
}

func TestShareCodeImmutability(t *testing.T) {
	s, _ := createNewSurvey().Save(db)

	if s.ShareCode == "" {
		t.Errorf("Expected share code, got %s", s.ShareCode)
	}

	s2 := createNewSurvey()
	s2.Name = s.Name
	s2.ID = s.ID
	s2, _ = s2.Save(db)

	if s.ShareCode != s2.ShareCode {
		t.Errorf("Expected %s, got %s", s.ShareCode, s2.ShareCode)
	}
}

func TestAttemptToChangeID(t *testing.T) {
	s, _ := createNewSurvey().Save(db)

	if s.ShareCode == "" {
		t.Errorf("Expected share code, got %s", s.ShareCode)
	}

	s2 := createNewSurvey()
	s2.Name = s.Name
	s2.ID = 987

	t.Logf("%d", s2.ID)

	_, err := s2.Save(db)

	if err == nil {
		t.Errorf("Expected 'survey with that name already exists', got %s", err)
	}
}

func TestPackGenerateValid(t *testing.T) {
	PackGenerate("http://example.com")
}

func TestPackGenerateInvalid(t *testing.T) {
	PackGenerate("http://thisdoesnotoiajfisufh329hriu23bij23by8dfysdiufhksjdfbkjsbdf.com")
}

func TestGetByNameAndShareCode(t *testing.T) {
	s, _ := createNewSurvey().Save(db)

	ss, err := s.GetByNameAndShareCode(db, s.Name, s.ShareCode)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	if ss.ID != s.ID {
		t.Errorf("Expected %d, got %d", s.ID, ss.ID)
	}
}

func TestGetByNameAndShareCodeInvalid(t *testing.T) {
	s, _ := createNewSurvey().Save(db)

	_, err := s.GetByNameAndShareCode(db, s.Name, "invalidShareCode")

	if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found', got %s", err)
	}
}

func TestDuplicateSurveyName(t *testing.T) {
	s1 := createNewSurvey()
	s2 := createNewSurvey()

	_, err := s1.Save(db)

	if err != nil {
		t.Errorf("Expected nil, got %s", err)
	}

	s2.Name = s1.Name

	if _, err := s2.Save(db); err.Error() != "UNIQUE constraint failed: surveys.name" {
		t.Errorf("Expected 'UNIQUE constraint failed: surveys.name', got %s", err)
	}
}

func TestCleanup(t *testing.T) {
	deleteTestDB()
}
