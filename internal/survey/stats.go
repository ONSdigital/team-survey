package survey

import (
	"fmt"
	"math"
)

// CalculatedResult is the representation of Westrum/Lencioni etc.
type CalculatedResult struct {
	Team                        string  `json:"team"`
	LeadTime                    string  `json:"lead_time"`
	DeploymentFrequency         string  `json:"deployment_frequency"`
	Mttr                        string  `json:"mttr"`
	ChangeFailure               string  `json:"change_failure"`
	Westrum                     string  `json:"westrum_result"`
	WestrumScore                int     `json:"westrum_score"`
	LencioniTrust               string  `json:"lencioni_trust"`
	LencioniTrustScore          float64 `json:"lencioni_trust_score"`
	LencioniConflict            string  `json:"lencioni_conflict"`
	LencioniConflictScore       float64 `json:"lencioni_conflict_score"`
	LencioniCommitment          string  `json:"lencioni_commitment"`
	LencioniCommitmentScore     float64 `json:"lencioni_commitment_score"`
	LencioniAccountability      string  `json:"lencioni_accountability"`
	LencioniAccountabilityScore float64 `json:"lencioni_accountability_score"`
	LencioniResults             string  `json:"lencioni_results"`
	LencioniResultsScore        float64 `json:"lencioni_results_score"`
	ResponseTotal               int     `json:"total_responses"`
	DateStart                   string  `json:"date_start"`
	DateEnd                     string  `json:"date_end"`
	ShareCode                   string  `json:"share_code"`
}

// StringTotal is a struct for holding counts of a string value
type StringTotal struct {
	Text  string
	Count int
}

// GroupTotalStats holds the raw totals split into groups
type GroupTotalStats struct {
	LeadTime        []StringTotal
	Mttr            []StringTotal
	DeployFrequency []StringTotal
	ChangeFailure   []StringTotal
	Lencioni        struct {
		Trust          int
		Conflict       int
		Commitment     int
		Accountability int
		Results        int
	}
	Westrum   int
	Responses int
}

// Calculate will convert a survey into westrum/lencioni etc.
func Calculate(s *Survey) (CalculatedResult, error) {
	if len(s.Results) == 0 {
		c := CalculatedResult{
			Team:                s.Team,
			LeadTime:            "No Results",
			DeploymentFrequency: "No Results",
			Mttr:                "No Results",
			ChangeFailure:       "No Results",
		}
		return c, fmt.Errorf("no results to calculate")
	}

	questionGroups := SplitSurveyIntoQuestionGroups(s)
	result := CalculateResultFromGroupTotalStats(questionGroups)
	result.Team = s.Team
	return result, nil
}

// CountStringInstance counts the instances of a string and updates the map
func CountStringInstance(text string, store []StringTotal) []StringTotal {
	var elementFound bool
	for index, element := range store {
		if element.Text == text {
			store[index].Count++
			elementFound = true
		}
	}

	if !elementFound {
		store = append(store, StringTotal{
			Text:  text,
			Count: 1,
		})
	}
	return store
}

// SplitSurveyIntoQuestionGroups will take questions and split them out into the question group
func SplitSurveyIntoQuestionGroups(s *Survey) GroupTotalStats {

	var totalStats GroupTotalStats

	for _, result := range s.Results {
		totalStats.Responses++

		totalStats.LeadTime = CountStringInstance(result.LeadTime, totalStats.LeadTime)
		totalStats.Mttr = CountStringInstance(result.Mttr, totalStats.Mttr)
		totalStats.DeployFrequency = CountStringInstance(result.DeploymentFrequency, totalStats.DeployFrequency)
		totalStats.ChangeFailure = CountStringInstance(result.ChangeFailure, totalStats.ChangeFailure)

		totalStats.Lencioni.Trust += result.AboutTheTeamExtendedApologise
		totalStats.Lencioni.Trust += result.AboutTheTeamExtendedAdmitMistakes
		totalStats.Lencioni.Trust += result.AboutTheTeamExtendedComfortableDiscussingPersonalLives

		totalStats.Lencioni.Conflict += result.AboutTheTeamExtendedUnguardedDiscussion
		totalStats.Lencioni.Conflict += result.AboutTheTeamExtendedCompellingMeetings
		totalStats.Lencioni.Conflict += result.AboutTheTeamExtendedDifficultIssuesRaised

		totalStats.Lencioni.Commitment += result.AboutTheTeamExtendedContributeToCollectiveGood
		totalStats.Lencioni.Commitment += result.AboutTheTeamExtendedLeaveMeetingsCommitted
		totalStats.Lencioni.Commitment += result.AboutTheTeamExtendedClearResolutionDiscussions

		totalStats.Lencioni.Accountability += result.AboutTheTeamExtendedCallOutUnproductiveBehaviour
		totalStats.Lencioni.Accountability += result.AboutTheTeamExtendedConcernedAboutLettingDownPeers
		totalStats.Lencioni.Accountability += result.AboutTheTeamExtendedChallengeOneAnother

		totalStats.Lencioni.Results += result.AboutTheTeamExtendedWillinglyMakeSacrifices
		totalStats.Lencioni.Results += result.AboutTheTeamExtendedMoraleAffectedByFailure
		totalStats.Lencioni.Results += result.AboutTheTeamExtendedSlowToSeekCredit

		totalStats.Westrum += result.AboutTheTeamCollaborationEncouraged
		totalStats.Westrum += result.AboutTheTeamNewIdeasWelcomed
		totalStats.Westrum += result.AboutTheTeamMessengersNotPunished
		totalStats.Westrum += result.AboutTheTeamFailureTreatedAsOpportunity
		totalStats.Westrum += result.AboutTheTeamInformationActivelySought
		totalStats.Westrum += result.AboutTheTeamResponsibilitiesShared
		totalStats.Westrum += result.AboutTheTeamFailureCausesEnquiry
	}

	return totalStats

}

// CalculateResultFromGroupTotalStats will calculate the resuls from GroupTotalStats
func CalculateResultFromGroupTotalStats(s GroupTotalStats) CalculatedResult {
	var result CalculatedResult

	lencioniAmberFloor := 6.00
	lencioniGreenFloor := 8.00
	westrumAmberFloor := 3
	westrumGreenFloor := 5

	result.LencioniTrustScore = float64(s.Lencioni.Trust / s.Responses)
	result.LencioniTrust = GetRagResult(result.LencioniTrustScore, lencioniAmberFloor, lencioniGreenFloor)
	result.LencioniResultsScore = float64(s.Lencioni.Results / s.Responses)
	result.LencioniResults = GetRagResult(result.LencioniTrustScore, lencioniAmberFloor, lencioniGreenFloor)
	result.LencioniConflictScore = float64(s.Lencioni.Conflict / s.Responses)
	result.LencioniConflict = GetRagResult(result.LencioniConflictScore, lencioniAmberFloor, lencioniGreenFloor)
	result.LencioniCommitmentScore = float64(s.Lencioni.Commitment / s.Responses)
	result.LencioniCommitment = GetRagResult(result.LencioniCommitmentScore, lencioniAmberFloor, lencioniGreenFloor)
	result.LencioniAccountabilityScore = float64(s.Lencioni.Accountability / s.Responses)
	result.LencioniAccountability = GetRagResult(result.LencioniAccountabilityScore, lencioniAmberFloor, lencioniGreenFloor)

	result.WestrumScore = (s.Westrum / 7) / s.Responses
	result.Westrum = GetWestrumResult(result.WestrumScore, westrumAmberFloor, westrumGreenFloor)

	// Round long floats down to 2 decimal places
	result.LencioniTrustScore = (math.Floor(result.LencioniTrustScore*100) / 100)
	result.LencioniResultsScore = (math.Floor(result.LencioniResultsScore*100) / 100)
	result.LencioniConflictScore = (math.Floor(result.LencioniConflictScore*100) / 100)
	result.LencioniCommitmentScore = (math.Floor(result.LencioniCommitmentScore*100) / 100)
	result.LencioniAccountabilityScore = (math.Floor(result.LencioniAccountabilityScore*100) / 100)

	result.ResponseTotal = s.Responses

	result.Mttr = CalculateKeyMetrics(s.Mttr).Text
	result.LeadTime = CalculateKeyMetrics(s.LeadTime).Text
	result.DeploymentFrequency = CalculateKeyMetrics(s.DeployFrequency).Text
	result.ChangeFailure = CalculateKeyMetrics(s.ChangeFailure).Text

	return result
}

// GetRagResult returns a RAG status based on input
func GetRagResult(score float64, amberFloor float64, greenFloor float64) string {
	if score >= greenFloor {
		return "green"
	}

	if score >= amberFloor {
		return "amber"
	}

	return "red"
}

// GetWestrumResult returns a Westrum status based on input
func GetWestrumResult(score int, amberFloor int, greenFloor int) string {
	if score >= greenFloor {
		return "generative"
	}

	if score >= amberFloor {
		return "bureaucratic"
	}

	return "pathological"
}

// CalculateKeyMetrics gives a StringTotal of the most voted result
func CalculateKeyMetrics(m []StringTotal) StringTotal {

	var total StringTotal

	if len(m) == 0 {
		total.Text = "No metrics to calculate"
		total.Count = 0
		return total
	}

	highestMetric := m[0]

	if len(m) == 1 {
		total.Text = highestMetric.Text
		total.Count = highestMetric.Count

		return total
	}

	for _, metric := range m {
		if metric.Count > highestMetric.Count {
			highestMetric.Text = metric.Text
			highestMetric.Count = metric.Count
		} else if metric.Count == highestMetric.Count {
			highestMetric.Text = "Unknown"
			highestMetric.Count = (highestMetric.Count - m[0].Count) + metric.Count
		}
	}

	return highestMetric
}
