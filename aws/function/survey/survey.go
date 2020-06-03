package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type (
	// JSKVPair struct
	JSKVPair struct {
		Value interface{} `json:"value" yaml:"value"`
		Text  string      `json:"text" yaml:"text"`
	}

	// JSElement struct
	JSElement struct {
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
	JSPage struct {
		Name     string      `json:"name" yaml:"name"`
		Title    string      `json:"title" yaml:"title"`
		Elements []JSElement `json:"elements" yaml:"elements"`
	}

	// JSSurvey struct
	JSSurvey struct {
		CompletedHTML       template.HTML `json:"completedHtml" yaml:"completedHtml"`
		Pages               []JSPage      `json:"pages" yaml:"pages"`
		ShowQuestionNumbers string        `json:"showQuestionNumbers" yaml:"showQuestionNumbers"`
	}
)

// YamlFileToSurveyJSBuilder holds the componets to build a SurveyJS survey from YAML
type yamlFileToSurveyJSBuilder struct {
	Survey string
	Pages  []string
}

func buildSurvey(builder yamlFileToSurveyJSBuilder) (string, error) {
	var sjs JSSurvey
	d, err := ioutil.ReadFile(builder.Survey)
	if err != nil {
		return "", err
	}
	err = yaml.Unmarshal(d, &sjs)
	if err != nil {
		return "", err
	}

	for _, v := range builder.Pages {
		d, err := ioutil.ReadFile(v)
		if err != nil {
			return "", err
		}

		var tempSjs JSSurvey
		err = yaml.Unmarshal(d, &tempSjs)
		if err != nil {
			return "", err
		}
		sjs.Pages = append(sjs.Pages, tempSjs.Pages...)
	}

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	err = enc.Encode(sjs)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(buf.String()), err
}
