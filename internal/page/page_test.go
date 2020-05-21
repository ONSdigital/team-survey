package page_test

import (
	"html/template"
	"os"
	"testing"

	"github.com/onsdigital/team-survey/internal/page"
	"github.com/onsdigital/team-survey/internal/survey"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Override the path to our templates so that
	// the test can actually see them!
	page.Path = "../../assets/templates/"
	os.Exit(m.Run())
}

func TestRenderToString(t *testing.T) {
	data := page.SurveyData{
		CommonData: page.CommonData{
			Stage: "test",
		},
		SurveyJS: template.JS("null"),
		Survey:   survey.Survey{},
	}
	_, err := page.RenderToString(data, "survey", "survey")
	assert.NoError(t, err)
}
