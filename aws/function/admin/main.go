package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/onsdigital/team-survey/internal/page"
	"github.com/onsdigital/team-survey/internal/survey"
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

func router(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// TODO Authenticate routes

	switch req.Resource {
	// Web routes
	case "/admin/survey":
		return handleListSurveys(ctx, req)
	case "/admin/survey/{survey}":
		return handleSurveyDashboard(ctx, req)
	case "/admin/login":
		return handleGetLogin(ctx, req)
	// case "/admin/survey/new":
	// 	switch req.HTTPMethod {
	// 	case http.MethodGet:
	// 		return handleGetNewSurveyForm(ctx, req)
	// 	case http.MethodPost:
	// 		return handlePostNewSurveyForm(ctx, req)
	// 	}

	// API routes
	case "/api/v1/survey/{survey}":
		switch req.HTTPMethod {
		case http.MethodGet:
			return handleAPIReturnSurvey(ctx, req)
		}
	case "/api/v1/survey/{survey}/stats":
		switch req.HTTPMethod {
		case http.MethodGet:
			return handleAPIGetSurveyStats(ctx, req)
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
	}, nil
}

func handleGetLogin(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// .. TODO

	pageData := page.CommonData{
		Stage:      stage,
		StaticPath: assetsBucket,
	}

	page, err := page.RenderToString(pageData, "admin", "_base", "login")
	if err != nil {
		log.Println("error rendering login template:", err)
		return proxyError()
	}
	return htmlProxyResponse(page)
}

func handleListSurveys(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("listing existing surveys")

	surveys, err := survey.List(ctx, false)
	if err != nil {
		log.Println("error retrieving surveys:", err)
		return proxyError()
	}

	pageData := page.SurveyListData{
		CommonData: page.CommonData{
			Stage:      stage,
			StaticPath: assetsBucket,
		},
		Surveys: surveys,
	}

	page, err := page.RenderToString(pageData, "admin", "_base", "list-surveys")
	if err != nil {
		log.Println("error rendering survey template:", err)
		return proxyError()
	}
	return htmlProxyResponse(page)
}

func handleSurveyDashboard(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	surveyID := req.PathParameters["survey"]
	if surveyID == "" {
		log.Println("missing surveyID in proxy request")
		return proxyError()
	}

	srvy, err := survey.Get(ctx, surveyID)
	if err != nil {
		log.Printf("failed to load survey: %v", surveyID)
		return proxyError()
	}

	pageData := page.SurveyData{
		CommonData: page.CommonData{
			Stage:      stage,
			StaticPath: assetsBucket,
		},
		Survey: *srvy,
	}

	page, err := page.RenderToString(pageData, "admin", "_base", "survey-dashboard")
	if err != nil {
		log.Println("error rendering survey dashboard template:", err)
		return proxyError()
	}
	return htmlProxyResponse(page)
}

func proxyError() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
	}, nil
}

func htmlProxyResponse(body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/html"},
		Body:       body,
	}, nil
}

func jsonProxyResponse(body string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       body,
	}, nil
}

func main() {
	lambda.Start(router)
}
