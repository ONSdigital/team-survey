package survey

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/jinzhu/gorm"
	"github.com/kennygrant/sanitize"
	"github.com/rs/xid"
)

// Model is a GORM compatible model
type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created"`
	UpdatedAt time.Time  `json:"updated"`
	DeletedAt *time.Time `json:"-"`
}

// Result Struct
type Result struct {
	Model
	Time                                                   time.Time `json:"timestamp"`
	CurrentRole                                            string    `json:"current_role"`
	CurrentRoleOther                                       string    `json:"current_role_other"`
	EmploymentStatus                                       string    `json:"employment_status"`
	Team                                                   string    `json:"team"`
	ProgrammingLanguages1                                  string    `json:"programming_languages.1"`
	ProgrammingLanguages2                                  string    `json:"programming_languages.2"`
	ProgrammingLanguages3                                  string    `json:"programming_languages.3"`
	ProgrammingLanguages4                                  string    `json:"programming_languages.4"`
	ProgrammingLanguages5                                  string    `json:"programming_languages.5"`
	ProgrammingLanguages6                                  string    `json:"programming_languages.6"`
	ProgrammingLanguages7                                  string    `json:"programming_languages.7"`
	ProgrammingLanguages8                                  string    `json:"programming_languages.8"`
	ProgrammingLanguages9                                  string    `json:"programming_languages.9"`
	ProgrammingLanguages10                                 string    `json:"programming_languages.10"`
	ProgrammingLanguagesOther                              string    `json:"programming_languages_other"`
	AboutTheTeamInformationActivelySought                  int       `json:"about_the_team.information_actively_sought,string"`
	AboutTheTeamMessengersNotPunished                      int       `json:"about_the_team.messengers_not_punished,string"`
	AboutTheTeamResponsibilitiesShared                     int       `json:"about_the_team.responsibilities_shared,string"`
	AboutTheTeamCollaborationEncouraged                    int       `json:"about_the_team.collaboration_encouraged,string"`
	AboutTheTeamFailureCausesEnquiry                       int       `json:"about_the_team.failure_causes_enquiry,string"`
	AboutTheTeamNewIdeasWelcomed                           int       `json:"about_the_team.new_ideas_welcomed,string"`
	AboutTheTeamFailureTreatedAsOpportunity                int       `json:"about_the_team.failure_treated_as_opportunity,string"`
	AboutTheTeamExtendedUnguardedDiscussion                int       `json:"about_the_team_extended.unguarded_discussion,string"`
	AboutTheTeamExtendedCallOutUnproductiveBehaviour       int       `json:"about_the_team_extended.call_out_unproductive_behaviour,string"`
	AboutTheTeamExtendedContributeToCollectiveGood         int       `json:"about_the_team_extended.contribute_to_collective_good,string"`
	AboutTheTeamExtendedApologise                          int       `json:"about_the_team_extended.apologise,string"`
	AboutTheTeamExtendedWillinglyMakeSacrifices            int       `json:"about_the_team_extended.willingly_make_sacrifices,string"`
	AboutTheTeamExtendedAdmitMistakes                      int       `json:"about_the_team_extended.admit_mistakes,string"`
	AboutTheTeamExtendedCompellingMeetings                 int       `json:"about_the_team_extended.compelling_meetings,string"`
	AboutTheTeamExtendedLeaveMeetingsCommitted             int       `json:"about_the_team_extended.leave_meetings_committed,string"`
	AboutTheTeamExtendedMoraleAffectedByFailure            int       `json:"about_the_team_extended.morale_affected_by_failure,string"`
	AboutTheTeamExtendedDifficultIssuesRaised              int       `json:"about_the_team_extended.difficult_issues_raised,string"`
	AboutTheTeamExtendedConcernedAboutLettingDownPeers     int       `json:"about_the_team_extended.concerned_about_letting_down_peers,string"`
	AboutTheTeamExtendedComfortableDiscussingPersonalLives int       `json:"about_the_team_extended.comfortable_discussing_personal_lives,string"`
	AboutTheTeamExtendedClearResolutionDiscussions         int       `json:"about_the_team_extended.clear_resolution_discussions,string"`
	AboutTheTeamExtendedChallengeOneAnother                int       `json:"about_the_team_extended.challenge_one_another,string"`
	AboutTheTeamExtendedSlowToSeekCredit                   int       `json:"about_the_team_extended.slow_to_seek_credit,string"`
	LeadTime                                               string    `json:"lead_time"`
	DeploymentFrequency                                    string    `json:"deployment_frequency"`
	Mttr                                                   string    `json:"mttr"`
	ChangeFailure                                          string    `json:"change_failure"`
	SurveyID                                               uint      `json:"-"`
}

// Team Struct
type Team struct {
	Model
	Name string `json:"name"`
}

// Survey Struct
type Survey struct {
	Model
	Results    []Result `json:"results"`
	Team       Team     `json:"team"`
	Name       string   `gorm:"unique;not null;unique_index" json:"name"`
	TeamID     uint     `json:"-"`
	ShareCode  string   `json:"share_code"`
	Components []Component
}

// Component is the a component of the survey e.g. Westrum
type Component struct {
	Model
	Category string `json:"category"`
	Name     string `json:"name"`
	Weight   int    `json:"weight"`
	SurveyID uint   `json:"-"`
}

// ResultSearchParameters is used for searching Results
type ResultSearchParameters struct {
	DateStart string
	DateEnd   string
	Team      string
}

// JSKVPair struct
type JSKVPair struct {
	Value interface{} `json:"value" yaml:"value"`
	Text  string      `json:"text" yaml:"text"`
}

// JSElement struct
type JSElement struct {
	Type       string     `json:"type" yaml:"type"`
	Name       string     `json:"name" yaml:"name"`
	Title      string     `json:"title" yaml:"title"`
	IsRequired bool       `json:"isRequired" yaml:"isRequired"`
	ColCount   int        `json:"colCount" yaml:"colCount"`
	Choices    []string   `json:"choices,omitempty" yaml:"choices"`
	VisibleIf  string     `json:"visibleIf,omitempty" yaml:"visibleIf"`
	Columns    []JSKVPair `json:"columns,omitempty" yaml:"columns"`
	Rows       []JSKVPair `json:"rows,omitempty" yaml:"rows"`
}

// JSPage struct
type JSPage struct {
	Name     string      `json:"name" yaml:"name"`
	Title    string      `json:"title" yaml:"title"`
	Elements []JSElement `json:"elements" yaml:"elements"`
}

// JSSurvey struct
type JSSurvey struct {
	CompletedHTML       template.HTML `json:"completedHtml" yaml:"completedHtml"`
	Pages               []JSPage      `json:"pages" yaml:"pages"`
	ShowQuestionNumbers string        `json:"showQuestionNumbers" yaml:"showQuestionNumbers"`
}

// YamlFileToSurveyJSBuilder holds the componets to build a SurveyJS survey from YAML
type YamlFileToSurveyJSBuilder struct {
	Survey string
	Pages  []string
}

// AutoMigrateModels imports the models into the DB
func AutoMigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(&Result{}, &Survey{}, &Team{}, &Component{}).Error
}

// Get will return a survey by ID
func (s Survey) Get(db *gorm.DB, id uint) (Survey, error) {

	if err := db.First(&s, id).Error; err != nil {
		return s, err
	}

	return s, db.Model(&s).Related(&s.Results).Related(&s.Team).Error
}

// GetByName will return a survey by name
func (s Survey) GetByName(db *gorm.DB, name string) (Survey, error) {

	if err := db.Where("name = ?", name).First(&s).Error; err != nil {
		return s, err
	}

	return s, db.Model(&s).Related(&s.Results).Related(&s.Team).Error
}

// GetByNameAndShareCode will return a survey by name and share code
func (s Survey) GetByNameAndShareCode(db *gorm.DB, name, shareCode string) (Survey, error) {

	if err := db.Where("name = ? AND share_code= ?", name, shareCode).First(&s).Error; err != nil {
		return s, err
	}

	packGenerationHook(&s)

	return s, db.Model(&s).Related(&s.Results).Related(&s.Team).Error
}

// BeforeSave will sanitize input data to be safe
func (s *Survey) BeforeSave() {
	s.Name = sanitize.HTML(s.Name)
	s.Name = strings.ToLower(s.Name)
	s.Name = strings.Replace(s.Name, " ", "-", -1)
	s.Name = strings.ToLower(s.Name)

	if s.ShareCode == "" {
		s.ShareCode = xid.New().String()
	}
}

// AfterSave has some post-save hooks
func (s *Survey) AfterSave() {
	packGenerationHook(s)
}

// PackGenerate posts to the pack generation service
func PackGenerate(url string) {
	log.Printf("Post-save hook: %s", url)
	_, err := http.Post(url, "application/json", bytes.NewBuffer([]byte("Survey Saved")))
	if err != nil {
		return
	}
}

func packGenerationHook(s *Survey) {
	// Generation hook goes here
}

// Save will commit to the database
func (s Survey) Save(db *gorm.DB) (Survey, error) {

	if s.Name == "" {
		return s, fmt.Errorf("survey has no name")
	}

	if s.Team.Name == "" {
		return s, fmt.Errorf("team not defined")
	}
	// Avoid duplicate team names. Allow update by name only.
	if s.Team.Name != "" {
		var team Team
		db.Where("name = ?", s.Team.Name).First(&team)

		if team.ID > 0 {
			s.Team = team
		}
	}

	es, _ := s.GetByName(db, s.Name)
	if es.ID > 0 {
		s.ShareCode = es.ShareCode
	}

	return s, db.Save(&s).Error
}

// GetAll Returns a list of all surveys
func (s Survey) GetAll(db *gorm.DB) ([]Survey, error) {
	var surveys []Survey
	db.Preload("Results").Preload("Team").Order("created_at desc").Find(&surveys)
	return surveys, db.Model(&s).Error
}

// GetAllResults returns a survey with all results
func (s Survey) GetAllResults(db *gorm.DB, p ResultSearchParameters) (Survey, error) {
	var results []Result
	if p.Team != "" && p.Team != "All" {
		db.Where("DATE(created_at) BETWEEN ? AND ? AND team = ?", p.DateStart, p.DateEnd, p.Team).Find(&results)
	} else {
		db.Where("DATE(created_at) BETWEEN ? AND ?", p.DateStart, p.DateEnd).Find(&results)
	}

	s.Results = results
	return s, nil
}

// GetAll Returns a list of teams
func (t Team) GetAll(db *gorm.DB) ([]Team, error) {
	var teams []Team
	return teams, db.Order("name asc").Find(&teams).Error
}

// GetByName returns a team by name or an error
func (t Team) GetByName(db *gorm.DB, name string) (Team, error) {
	return t, db.Where("name = ?", name).First(&t).Error
}

// BeforeSave executes before a team is persisted
func (t *Team) BeforeSave() {
	t.Name = sanitize.HTML(t.Name)
}

// PopulateSurveyJSFromYML yaml > JSSurvey
func PopulateSurveyJSFromYML(sjs *JSSurvey, yamlString string) error {
	return yaml.Unmarshal([]byte(yamlString), &sjs)
}

// GenerateSurveyJSON generates JSON without escaping HTML
func GenerateSurveyJSON(sjs JSSurvey) (string, error) {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err := enc.Encode(sjs)

	return strings.TrimSpace(buf.String()), err
}

// BuildSurveyFromYMLFiles build a JSSurvey from a YamlFileToSurveyJSBuilder
func BuildSurveyFromYMLFiles(builder YamlFileToSurveyJSBuilder) (string, error) {
	var sjs JSSurvey
	d, _ := ioutil.ReadFile(builder.Survey)
	_ = PopulateSurveyJSFromYML(&sjs, string(d))

	for _, v := range builder.Pages {
		var tempSjs JSSurvey
		d, _ := ioutil.ReadFile(v)
		_ = PopulateSurveyJSFromYML(&tempSjs, string(d))
		sjs.Pages = append(sjs.Pages, tempSjs.Pages...)
	}

	return GenerateSurveyJSON(sjs)
}
