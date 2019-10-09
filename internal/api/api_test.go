package api

import (
	"fmt"
	"os"
	"testing"

	"github.com/ONSdigital/team-survey/internal/survey"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB
var testDBFileName string
var api API
var surveyCount int

func init() {
	testDBFileName = "test.db"
	deleteTestDB()
	var err error

	if db, err = gorm.Open("sqlite3", testDBFileName); err != nil {
		panic("failed to connect database")
	}
	api.DB = db
}

func deleteTestDB() {
	if _, err := os.Stat(testDBFileName); err == nil {
		os.Remove(testDBFileName)
		fmt.Println("=== INFO  Test DB Deleted")
	}
}

func createNewSurvey() survey.Survey {
	results := []survey.Result{}

	results = append(results, survey.Result{
		AboutTheTeamExtendedApologise:                          1,
		AboutTheTeamExtendedAdmitMistakes:                      2,
		AboutTheTeamExtendedComfortableDiscussingPersonalLives: 3,
		AboutTheTeamExtendedUnguardedDiscussion:                1,
		AboutTheTeamExtendedCompellingMeetings:                 2,
		AboutTheTeamExtendedDifficultIssuesRaised:              3,
		AboutTheTeamExtendedContributeToCollectiveGood:         1,
		AboutTheTeamExtendedLeaveMeetingsCommitted:             2,
		AboutTheTeamExtendedClearResolutionDiscussions:         3,
		AboutTheTeamExtendedCallOutUnproductiveBehaviour:       1,
		AboutTheTeamExtendedConcernedAboutLettingDownPeers:     2,
		AboutTheTeamExtendedChallengeOneAnother:                3,
		AboutTheTeamExtendedWillinglyMakeSacrifices:            1,
		AboutTheTeamExtendedMoraleAffectedByFailure:            2,
		AboutTheTeamExtendedSlowToSeekCredit:                   3,
		AboutTheTeamCollaborationEncouraged:                    1,
		AboutTheTeamNewIdeasWelcomed:                           2,
		AboutTheTeamMessengersNotPunished:                      3,
		AboutTheTeamFailureTreatedAsOpportunity:                1,
		AboutTheTeamInformationActivelySought:                  2,
		AboutTheTeamResponsibilitiesShared:                     3,
		AboutTheTeamFailureCausesEnquiry:                       1,
		LeadTime:                                               "Between one day and one week",
		Mttr:                                                   "More than six months",
		DeploymentFrequency:                                    "Between once per day and once per week",
		ChangeFailure:                                          "46%-60%",
	})

	surveyCount++

	return survey.Survey{
		Name: fmt.Sprintf("Test Survey %d", surveyCount),
		Team: survey.Team{
			Name: "Amazing Team",
		},
		Results: results,
	}
}

func TestAutoMigrateModels(t *testing.T) {
	err := survey.AutoMigrateModels(db)
	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}
}

func TestSaveSurvey(t *testing.T) {
	s, err := api.SaveSurvey(createNewSurvey())

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if s.ID == 0 {
		t.Errorf("Expected not 0, but got %d", s.ID)
	}
}

func TestGetStats(t *testing.T) {
	s, err := api.SaveSurvey(createNewSurvey())

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	result, err := api.GetStats(s)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if result.Westrum != "pathological" {
		t.Errorf("Expected 'pathological', but got %s", result.Westrum)
	}
	if result.WestrumScore != 1 {
		t.Errorf("Expected 1, but got %d", result.WestrumScore)
	}

	if result.LencioniTrustScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniTrustScore)
	}

	if result.LencioniTrust != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniTrust)
	}

	if result.LencioniResultsScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniTrustScore)
	}

	if result.LencioniResults != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniTrust)
	}

	if result.LencioniConflictScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniConflictScore)
	}

	if result.LencioniConflict != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniConflict)
	}

	if result.LencioniCommitmentScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniCommitmentScore)
	}

	if result.LencioniCommitment != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniCommitment)
	}

	if result.LencioniTrustScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniTrustScore)
	}

	if result.LencioniTrust != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniTrust)
	}

	if result.LencioniAccountabilityScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniAccountabilityScore)
	}

	if result.LencioniAccountability != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniAccountability)
	}

	if result.LencioniAccountabilityScore != 6 {
		t.Errorf("Expected 6, but got %f", result.LencioniAccountabilityScore)
	}

	if result.LencioniAccountability != "amber" {
		t.Errorf("Expected 'amber', but got %s", result.LencioniAccountability)
	}
}

func TestGetSurvey(t *testing.T) {
	s, err := api.SaveSurvey(createNewSurvey())

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if s.ID == 0 {
		t.Errorf("Expected not 0, but got %d", s.ID)
	}

	rs, err := api.GetSurvey(s.ID)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if rs.ID != s.ID {
		t.Errorf("Expected %d, but got %d", s.ID, rs.ID)
	}
}

func TestGetSurveyByName(t *testing.T) {
	s, err := api.SaveSurvey(createNewSurvey())

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if s.ID == 0 {
		t.Errorf("Expected not 0, but got %d", s.ID)
	}

	rs, err := api.GetSurveyByName(s.Name)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if rs.Name != s.Name {
		t.Errorf("Expected %s, but got %s", s.Name, rs.Name)
	}
}

func TestGetSurveyByNameInvalid(t *testing.T) {
	_, err := api.GetSurveyByName("SproutsGalore")

	if err.Error() != "record not found" {
		t.Errorf("Expected 'record not found', but got %s", err)
	}
}

func TestGetSurveyByNameAndShareCode(t *testing.T) {
	s, err := api.SaveSurvey(createNewSurvey())

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if s.ID == 0 {
		t.Errorf("Expected not 0, but got %d", s.ID)
	}

	rs, err := api.GetSurveyByNameAndShareCode(s.Name, s.ShareCode)

	if err != nil {
		t.Errorf("Expected nil, but got %s", err)
	}

	if rs.Name != s.Name {
		t.Errorf("Expected %s, but got %s", s.Name, rs.Name)
	}
}

func TestCleanup(t *testing.T) {
	deleteTestDB()
}
