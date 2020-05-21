package survey_test

import (
	"testing"

	"github.com/onsdigital/team-survey/internal/survey"
	"github.com/stretchr/testify/assert"
)

func TestExists(t *testing.T) {
	cases := []struct {
		s      survey.Survey
		exists bool
	}{
		{s: survey.Survey{SurveyID: "someid", Open: true}, exists: true},
		{s: survey.Survey{SurveyID: "", Open: true}, exists: false},
		{s: survey.Survey{SurveyID: "someid", Open: false}, exists: false},
	}

	for _, tc := range cases {
		assert.Equal(t, tc.exists, tc.s.Exists())
	}
}
