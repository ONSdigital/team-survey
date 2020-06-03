package page

import (
	"html/template"
	"path"
	"strings"

	"github.com/onsdigital/team-survey/internal/survey"
)

// Path defines the root at which to look for templates.
var Path = "assets/templates"

type (
	// DataBlock defines a block of data that can be loaded into a page template.
	// Included for making code easier to read.
	DataBlock interface{}

	// CommonData is a block of data that is common to any page
	CommonData struct {
		Stage      string
		StaticPath string
	}

	// AdminData is a block of data that is common to admin pages
	AdminData struct {
		Section string
		AppUser struct { // TODO - bound with authentication bits!
			Email           string
			ProfileImageURL string
		}
	}

	// TODO Should these definitions live here? Seems more reasonable that they
	// should live with the code that creates them unless they need to be more
	// widely shared across responsibilities

	// SurveyData is a block of data to be interpolated into a survey page
	SurveyData struct {
		CommonData
		SurveyJS template.JS
		Survey   survey.Survey
	}

	// SurveyListData is a block of data to be interpolated into a page listing
	// a set of surveys
	SurveyListData struct {
		CommonData
		AdminData
		Surveys []*survey.Survey
	}
)

// RenderToString attempts to render a template (assets/templates/<name>.html.tmpl)
// to a string
func RenderToString(data DataBlock, group string, names ...string) (string, error) {

	files := make([]string, len(names))

	for i, n := range names {
		files[i] = path.Join(Path, group, n+".html.tmpl")
	}

	t, err := template.ParseFiles(files...)
	if err != nil {
		return "", err
	}

	b := &strings.Builder{}
	err = t.Execute(b, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
