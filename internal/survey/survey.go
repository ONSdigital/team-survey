// Package survey contains all the interactions related to creating and running
// a survey, including data persistence.
package survey

import (
	"time"
)

type (
	// Survey is a top level representation of a survey
	Survey struct {
		SurveyID  string    `json:"surveyId"`
		Name      string    `json:"name"`
		Team      string    `json:"team"`
		Results   []Result  `json:"results"`
		CreatedAt time.Time `json:"created_at"`
		Open      bool      `json:"open"`
	}

	// Result is a single respondent's result set for a survey. The fields mapped
	// in the `json` marshaling MUST match the fields defined in the survey yaml
	// configurations. Pay attention to indexes (such as Programming Languages
	// starts at zero!).
	Result struct {
		ResultID string    `json:"resultId"`
		Time     time.Time `json:"timestamp"`

		// ONS About
		CurrentRole               string `json:"current_role"`
		CurrentRoleOther          string `json:"current_role_other"`
		EmploymentStatus          string `json:"employment_status"`
		ProgrammingLanguages0     string `json:"programming_languages.0"`
		ProgrammingLanguages1     string `json:"programming_languages.1"`
		ProgrammingLanguages2     string `json:"programming_languages.2"`
		ProgrammingLanguages3     string `json:"programming_languages.3"`
		ProgrammingLanguages4     string `json:"programming_languages.4"`
		ProgrammingLanguages5     string `json:"programming_languages.5"`
		ProgrammingLanguages6     string `json:"programming_languages.6"`
		ProgrammingLanguages7     string `json:"programming_languages.7"`
		ProgrammingLanguages8     string `json:"programming_languages.8"`
		ProgrammingLanguages9     string `json:"programming_languages.9"`
		ProgrammingLanguages10    string `json:"programming_languages.10"`
		ProgrammingLanguagesOther string `json:"programming_languages_other"`

		// Culture Westrum
		AboutTheTeamInformationActivelySought   int `json:"about_the_team.information_actively_sought,string"`
		AboutTheTeamMessengersNotPunished       int `json:"about_the_team.messengers_not_punished,string"`
		AboutTheTeamResponsibilitiesShared      int `json:"about_the_team.responsibilities_shared,string"`
		AboutTheTeamCollaborationEncouraged     int `json:"about_the_team.collaboration_encouraged,string"`
		AboutTheTeamFailureCausesEnquiry        int `json:"about_the_team.failure_causes_enquiry,string"`
		AboutTheTeamNewIdeasWelcomed            int `json:"about_the_team.new_ideas_welcomed,string"`
		AboutTheTeamFailureTreatedAsOpportunity int `json:"about_the_team.failure_treated_as_opportunity,string"`

		// Cohesion Lencioni
		AboutTheTeamExtendedUnguardedDiscussion                int `json:"about_the_team_extended.unguarded_discussion,string"`
		AboutTheTeamExtendedCallOutUnproductiveBehaviour       int `json:"about_the_team_extended.call_out_unproductive_behaviour,string"`
		AboutTheTeamExtendedContributeToCollectiveGood         int `json:"about_the_team_extended.contribute_to_collective_good,string"`
		AboutTheTeamExtendedApologise                          int `json:"about_the_team_extended.apologise,string"`
		AboutTheTeamExtendedWillinglyMakeSacrifices            int `json:"about_the_team_extended.willingly_make_sacrifices,string"`
		AboutTheTeamExtendedAdmitMistakes                      int `json:"about_the_team_extended.admit_mistakes,string"`
		AboutTheTeamExtendedCompellingMeetings                 int `json:"about_the_team_extended.compelling_meetings,string"`
		AboutTheTeamExtendedLeaveMeetingsCommitted             int `json:"about_the_team_extended.leave_meetings_committed,string"`
		AboutTheTeamExtendedMoraleAffectedByFailure            int `json:"about_the_team_extended.morale_affected_by_failure,string"`
		AboutTheTeamExtendedDifficultIssuesRaised              int `json:"about_the_team_extended.difficult_issues_raised,string"`
		AboutTheTeamExtendedConcernedAboutLettingDownPeers     int `json:"about_the_team_extended.concerned_about_letting_down_peers,string"`
		AboutTheTeamExtendedComfortableDiscussingPersonalLives int `json:"about_the_team_extended.comfortable_discussing_personal_lives,string"`
		AboutTheTeamExtendedClearResolutionDiscussions         int `json:"about_the_team_extended.clear_resolution_discussions,string"`
		AboutTheTeamExtendedChallengeOneAnother                int `json:"about_the_team_extended.challenge_one_another,string"`
		AboutTheTeamExtendedSlowToSeekCredit                   int `json:"about_the_team_extended.slow_to_seek_credit,string"`

		// Metrics Accelerate
		LeadTime            string `json:"lead_time"`
		DeploymentFrequency string `json:"deployment_frequency"`
		Mttr                string `json:"mttr"`
		ChangeFailure       string `json:"change_failure"`
	}
)

// Exists returns whether the given survey exists and is validly open
func (s *Survey) Exists() bool {
	return s.SurveyID != "" && s.Open
}
