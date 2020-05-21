package survey_test

import (
	"testing"

	// "github.com/onsdigital/team-survey-lambda/internal/survey"

	"github.com/google/go-cmp/cmp"
	"github.com/onsdigital/team-survey-lambda/internal/survey"
)

func createNewSurvey() *survey.Survey {
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

	return &survey.Survey{
		Name:    "Test Survey 1",
		Team:    "Amazing Team",
		Results: results,
	}
}

func TestCountStringInstance(t *testing.T) {
	testData := []survey.StringTotal{
		{
			Text:  "Hello",
			Count: 1,
		},
		{
			Text:  "Goodbye",
			Count: 1,
		},
	}

	counts := survey.CountStringInstance("Hello", testData)

	expected := []survey.StringTotal{
		{
			Text:  "Hello",
			Count: 2,
		},
		{
			Text:  "Goodbye",
			Count: 1,
		},
	}

	if !cmp.Equal(counts, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, counts)
	}

	counts = survey.CountStringInstance("Hello", testData)

	expected = []survey.StringTotal{
		{
			Text:  "Hello",
			Count: 3,
		},
		{
			Text:  "Goodbye",
			Count: 1,
		},
	}

	if !cmp.Equal(counts, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, counts)
	}

	counts = survey.CountStringInstance("Good Afternoon", testData)

	expected = []survey.StringTotal{
		{
			Text:  "Hello",
			Count: 3,
		},
		{
			Text:  "Goodbye",
			Count: 1,
		},
		{
			Text:  "Good Afternoon",
			Count: 1,
		},
	}

	if !cmp.Equal(counts, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, counts)
	}
}

func TestCalculateWhenNoResults(t *testing.T) {
	s := createNewSurvey()
	s.Results = nil

	_, err := survey.Calculate(s)

	if err == nil || err.Error() != "no results to calculate" {
		t.Errorf("Expected 'no results to calculate', but got %s", err)
	}
}

func TestCalculateWhenOneResult(t *testing.T) {
	s := createNewSurvey()

	_, err := survey.Calculate(s)

	if err != nil {
		t.Errorf("Expected 'nil', but got %s", err)
	}
}

func TestSplitSurveyIntoQuestionGroupsDefaultValues(t *testing.T) {
	s := createNewSurvey()
	s.Results = nil

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	if groups.Lencioni.Trust != 0 {
		t.Errorf("Expected 0, but got %d", groups.Lencioni.Trust)
	}

	if groups.Lencioni.Conflict != 0 {
		t.Errorf("Expected 0, but got %d", groups.Lencioni.Conflict)
	}

	if groups.Lencioni.Commitment != 0 {
		t.Errorf("Expected 0, but got %d", groups.Lencioni.Commitment)
	}

	if groups.Lencioni.Accountability != 0 {
		t.Errorf("Expected 0, but got %d", groups.Lencioni.Accountability)
	}

	if groups.Lencioni.Results != 0 {
		t.Errorf("Expected 0, but got %d", groups.Lencioni.Results)
	}

	if groups.Westrum != 0 {
		t.Errorf("Expected 0, but got %d", groups.Westrum)
	}
}

func TestSplitSurveyIntoQuestionGroupsLencioni(t *testing.T) {
	s := createNewSurvey()

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	if groups.Lencioni.Trust != 6 {
		t.Errorf("Expected 6, but got %d", groups.Lencioni.Trust)
	}

	if groups.Lencioni.Conflict != 6 {
		t.Errorf("Expected 6, but got %d", groups.Lencioni.Conflict)
	}

	if groups.Lencioni.Commitment != 6 {
		t.Errorf("Expected 6, but got %d", groups.Lencioni.Commitment)
	}

	if groups.Lencioni.Accountability != 6 {
		t.Errorf("Expected 6, but got %d", groups.Lencioni.Accountability)
	}

	if groups.Lencioni.Results != 6 {
		t.Errorf("Expected 6, but got %d", groups.Lencioni.Results)
	}
}

func TestSplitSurveyIntoQuestionGroupsWestrum(t *testing.T) {
	s := createNewSurvey()

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	if groups.Westrum != 13 {
		t.Errorf("Expected 6, but got %d", groups.Westrum)
	}
}

func TestSplitSurveyIntoQuestionGroupsLeadTime(t *testing.T) {
	s := createNewSurvey()

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	expected := []survey.StringTotal{
		{
			Text:  "Between one day and one week",
			Count: 1,
		},
	}

	if !cmp.Equal(groups.LeadTime, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, groups.LeadTime)
	}

	s.Results = append(s.Results, survey.Result{
		LeadTime: "Between one day and one week",
	})

	groups = survey.SplitSurveyIntoQuestionGroups(s)

	expected = []survey.StringTotal{
		{
			Text:  "Between one day and one week",
			Count: 2,
		},
	}

	if !cmp.Equal(groups.LeadTime, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, groups.LeadTime)
	}

	s.Results = append(s.Results, survey.Result{
		LeadTime: "More than 6 months",
	})

	groups = survey.SplitSurveyIntoQuestionGroups(s)

	expected = []survey.StringTotal{
		{
			Text:  "Between one day and one week",
			Count: 2,
		},
		{
			Text:  "More than 6 months",
			Count: 1,
		},
	}

	if !cmp.Equal(groups.LeadTime, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, groups.LeadTime)
	}
}

func TestSplitSurveyIntoQuestionGroupsMTTR(t *testing.T) {
	s := createNewSurvey()

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	expected := []survey.StringTotal{
		{
			Text:  "More than six months",
			Count: 1,
		},
	}

	if !cmp.Equal(groups.Mttr, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, groups.Mttr)
	}

	s.Results = append(s.Results, survey.Result{
		Mttr: "More than six months",
	})

	groups = survey.SplitSurveyIntoQuestionGroups(s)

	expected = []survey.StringTotal{
		{
			Text:  "More than six months",
			Count: 2,
		},
	}

	if !cmp.Equal(groups.Mttr, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, groups.Mttr)
	}

	s.Results = append(s.Results, survey.Result{
		Mttr: "More than 6 months",
	})

	groups = survey.SplitSurveyIntoQuestionGroups(s)

	expected = []survey.StringTotal{
		{
			Text:  "More than six months",
			Count: 2,
		},
		{
			Text:  "More than 6 months",
			Count: 1,
		},
	}

	if !cmp.Equal(groups.Mttr, expected) {
		t.Errorf("Expected %+v, but got %+v", expected, groups.Mttr)
	}
}

func TestGetRagResult(t *testing.T) {
	if s := survey.GetRagResult(3, 4, 5); s != "red" {
		t.Errorf("Expected red, but got %s", s)
	}
	if s := survey.GetRagResult(4, 4, 5); s != "amber" {
		t.Errorf("Expected amber, but got %s", s)
	}
	if s := survey.GetRagResult(4.5, 4, 5); s != "amber" {
		t.Errorf("Expected amber, but got %s", s)
	}
	if s := survey.GetRagResult(5, 4, 5); s != "green" {
		t.Errorf("Expected green, but got %s", s)
	}
	if s := survey.GetRagResult(5.5, 4, 5); s != "green" {
		t.Errorf("Expected green, but got %s", s)
	}
}

func TestGetWestrumResult(t *testing.T) {
	if s := survey.GetWestrumResult(3, 16, 31); s != "pathological" {
		t.Errorf("Expected pathological, but got %s", s)
	}
	if s := survey.GetWestrumResult(15, 16, 31); s != "pathological" {
		t.Errorf("Expected pathological, but got %s", s)
	}
	if s := survey.GetWestrumResult(16, 16, 31); s != "bureaucratic" {
		t.Errorf("Expected bureaucratic, but got %s", s)
	}
	if s := survey.GetWestrumResult(30, 16, 31); s != "bureaucratic" {
		t.Errorf("Expected bureaucratic, but got %s", s)
	}
	if s := survey.GetWestrumResult(31, 16, 31); s != "generative" {
		t.Errorf("Expected generative, but got %s", s)
	}
	if s := survey.GetWestrumResult(50, 16, 31); s != "generative" {
		t.Errorf("Expected generative, but got %s", s)
	}
}

func TestCalculateResultFromGroupTotalStatsLencioni(t *testing.T) {
	s := createNewSurvey()

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	result := survey.CalculateResultFromGroupTotalStats(groups)

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

func TestCalculateResultFromGroupTotalStatsWestrum(t *testing.T) {
	s := createNewSurvey()

	groups := survey.SplitSurveyIntoQuestionGroups(s)

	result := survey.CalculateResultFromGroupTotalStats(groups)

	if result.WestrumScore != 1 {
		t.Errorf("Expected 1, but got %d", result.WestrumScore)
	}

	if result.Westrum != "pathological" {
		t.Errorf("Expected 'pathological', but got %s", result.Westrum)
	}
}

func TestCalculateKeyMetricsEmpty(t *testing.T) {
	emptyMetric := survey.CalculateKeyMetrics([]survey.StringTotal{})

	if emptyMetric.Text != "No metrics to calculate" {
		t.Errorf("Expected 'No metrics to calculate', got %s", emptyMetric.Text)
	}

	if emptyMetric.Count != 0 {
		t.Errorf("Expected '0', got %d", emptyMetric.Count)
	}
}

func TestCalculateKeyMetrics(t *testing.T) {
	s := createNewSurvey()
	groups := survey.SplitSurveyIntoQuestionGroups(s)

	mttr := survey.CalculateKeyMetrics(groups.Mttr)

	if mttr.Text != "More than six months" {
		t.Errorf("Expected 'More than six months', got %s", mttr.Text)
	}
	if mttr.Count != 1 {
		t.Errorf("Expected 1, got %d", mttr.Count)
	}
}

func TestCalculateKeyMetricsUncertainty(t *testing.T) {
	s := createNewSurvey()
	s.Results = append(s.Results, survey.Result{
		LeadTime:            "Between one day and one week",
		Mttr:                "Less than six months",
		DeploymentFrequency: "Between once per day and once per week",
		ChangeFailure:       "46%-60%",
	})
	groups := survey.SplitSurveyIntoQuestionGroups(s)

	mttr := survey.CalculateKeyMetrics(groups.Mttr)

	if mttr.Text != "Unknown" {
		t.Errorf("Expected 'Unknown', got %s", mttr.Text)
	}

	if mttr.Count != 1 {
		t.Errorf("Expected 1, got %d", mttr.Count)
	}
}

func TestCalculateKeyMetricsCertainty(t *testing.T) {
	s := createNewSurvey()
	s.Results = nil
	s.Results = append(s.Results, survey.Result{
		LeadTime:            "Between one day and one week",
		Mttr:                "More than six months",
		DeploymentFrequency: "Between once per day and once per week",
		ChangeFailure:       "46%-60%",
	})
	s.Results = append(s.Results, survey.Result{
		LeadTime:            "Between one day and one week",
		Mttr:                "Less than six months",
		DeploymentFrequency: "Between once per day and once per week",
		ChangeFailure:       "46%-60%",
	})
	s.Results = append(s.Results, survey.Result{
		LeadTime:            "Between one day and one week",
		Mttr:                "Less than six months",
		DeploymentFrequency: "Between once per day and once per week",
		ChangeFailure:       "46%-60%",
	})
	s.Results = append(s.Results, survey.Result{
		LeadTime:            "Between one day and one week",
		Mttr:                "Less than six months",
		DeploymentFrequency: "Between once per day and once per week",
		ChangeFailure:       "46%-60%",
	})
	groups := survey.SplitSurveyIntoQuestionGroups(s)

	mttr := survey.CalculateKeyMetrics(groups.Mttr)

	if mttr.Text != "Less than six months" {
		t.Errorf("Expected 'Less than six months', got %s", mttr.Text)
	}

	if mttr.Count != 3 {
		t.Errorf("Expected 2, got %d", mttr.Count)
	}
}
