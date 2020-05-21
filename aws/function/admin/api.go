package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/onsdigital/team-survey/internal/survey"
)

func handleAPIGetSurveyStats(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	surveyID := req.PathParameters["survey"]
	if surveyID == "" {
		log.Println("missing surveyID in proxy request")
		return proxyError()
	}

	srvy, err := survey.Get(ctx, surveyID)
	if err != nil {
		log.Printf("error finding details for survey '%s': %v", surveyID, err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "survey not found"}`,
		}, nil
	}

	stats, err := survey.Calculate(srvy)
	if err != nil {
		log.Printf("error calculating stats for survey '%s': %v", surveyID, err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "failed to calculate stats"}`,
		}, nil
	}

	j, err := json.Marshal(stats)
	if err != nil {
		log.Printf("error marshaling JSON: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "failed to marshal stats JSON"}`,
		}, nil
	}
	return jsonProxyResponse(string(j))
}

func handleAPIReturnSurvey(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	surveyID := req.PathParameters["survey"]
	if surveyID == "" {
		log.Println("missing surveyID in proxy request")
		return proxyError()
	}

	srvy, err := survey.Get(ctx, surveyID)
	if err != nil {
		log.Printf("error finding details for survey '%s': %v", surveyID, err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "survey not found"}`,
		}, nil
	}

	j, err := json.Marshal(srvy)
	if err != nil {
		log.Printf("error marshaling JSON: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       `{"error": "failed to marshal response JSON"}`,
		}, nil
	}
	return jsonProxyResponse(string(j))
}
