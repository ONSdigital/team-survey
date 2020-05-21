package main

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/onsdigital/team-survey-lambda/internal/page"
	"github.com/onsdigital/team-survey-lambda/internal/survey"
)

var (
	// Common environment vars
	stage        string
	assetsBucket string
)

func init() {
	stage = os.Getenv("STAGE") + "/"
	assetsBucket = os.Getenv("S3_ASSETS_BUCKET")
}

var surveyTableName string

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	surveyTableName = "SurveyTable-" + stage

	switch req.Resource {
	// Web routes
	case "/survey/{survey}":
		return handleSurvey(ctx, req)

	// API routes
	case "/api/v1/survey/{survey}":
		switch req.HTTPMethod {
		case http.MethodPost:
			return handleAPISurveyResponse(ctx, req)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
	}, nil
}

func handleAPISurveyResponse(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	surveyID := req.PathParameters["survey"]
	if surveyID == "" {
		log.Println("missing surveyID in proxy request")
		return proxyError()
	}
	log.Println("received survey response")

	var s survey.Survey
	err := json.Unmarshal([]byte(req.Body), &s)
	if err != nil {
		log.Println("error parsing survey response body", err)
		return proxyError()
	}

	if surveyID != s.SurveyID {
		log.Printf("survey ID in resource path '%s' does not match value in form '%s'", surveyID, s.SurveyID)
		return proxyError()
	}

	_, err = survey.AddResult(ctx, s.SurveyID, s.Results[0])
	if err != nil {
		log.Println("error adding survey response:", err)
		return proxyError()
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
	}, nil
}

func handleSurvey(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	surveyID := req.PathParameters["survey"]
	if surveyID == "" {
		log.Println("missing surveyID in proxy request")
		return proxyError()
	}
	log.Println("loading survey ID:", surveyID)

	srvy, err := survey.Get(ctx, surveyID)
	if err != nil {
		log.Println(err)
		return proxyError()
	}

	if !srvy.Exists() {
		return events.APIGatewayProxyResponse{
			Body:       "Survey not found or is closed",
			StatusCode: http.StatusBadRequest,
		}, nil
	}

	builder := yamlFileToSurveyJSBuilder{}
	builder.Survey = "assets/survey-components/survey.yml"
	builder.Pages = []string{
		"assets/survey-components/about-ons.yml",
		"assets/survey-components/culture-westrum.yml",
		"assets/survey-components/cohesion-lencioni.yml",
		"assets/survey-components/metrics-accelerate.yml",
	}
	sjs, err := buildSurvey(builder)
	if err != nil {
		log.Println("error creating survey template:", err)
		return proxyError()
	}

	pageData := page.SurveyData{
		SurveyJS: template.JS(sjs),
		Survey: survey.Survey{
			Team:     srvy.Team,
			SurveyID: surveyID,
		},
		CommonData: page.CommonData{
			Stage:      stage,
			StaticPath: assetsBucket,
		},
	}

	page, err := page.RenderToString(pageData, "survey", "survey")
	if err != nil {
		log.Println("error rendering survey template:", err)
		return proxyError()
	}

	return events.APIGatewayProxyResponse{
		Body:       page,
		Headers:    map[string]string{"Content-Type": "text/html"},
		StatusCode: http.StatusOK,
	}, nil
}

func proxyError() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func main() {
	lambda.Start(router)
}
