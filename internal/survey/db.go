package survey

// Example survey record: (using DynamoDB text view)
//
//   {
//      "Open": true,							# Denotes survey is open or closed
//      "PK": "SURVEY",							# Fixed text "SURVEY"
//      "SK": "SURVEY#<SurveyID>",				# Built from "SURVEY#" + <SurveyID>
//      "SurveyID": "<SurveyID>",				# Unique ID of the survey
//      "Team": "<Name_of_team_taking_survey"	# Name of the team
//   }
//
// If you need to manually create a new ID, use:
//
// package main
// import (
// 	 "fmt"
// 	 "github.com/lithammer/shortuuid/v3"
// )
// func main() {
// 	 u := shortuuid.New()
// 	 fmt.Println(u)
// }

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/lithammer/shortuuid"
	"github.com/pkg/errors"
)

var surveyTableName string
var svc *dynamodb.DynamoDB

func init() {
	// TODO Need a better way of determining common loading stuff like this
	surveyTableName = os.Getenv("SURVEY_TABLE_NAME")
}

const (
	surveyPrefix = "SURVEY#"
	resultPrefix = "RESULT#"
)

// TODO should be done in such a way that we can easily abstract this out for testing
func getSession() *dynamodb.DynamoDB {
	if svc != nil {
		return svc
	}
	svc = dynamodb.New(session.New())
	return svc
}

type dynamoKeyFields struct {
	PK string
	SK string
}

type insertSurvey struct {
	dynamoKeyFields
	Survey
}

// Create stores a new survey record for the given team. The returned `Survey`
// contains the
func Create(ctx context.Context, team string) (*Survey, error) {

	srvy := Survey{
		SurveyID:  shortuuid.New(),
		Team:      team,
		Open:      true,
		Name:      "CakeSurvey",
		CreatedAt: time.Now(),
	}

	item := insertSurvey{
		dynamoKeyFields: dynamoKeyFields{
			PK: surveyPrefix + srvy.SurveyID,
			SK: "SURVEY",
		},
		Survey: srvy,
	}

	err := insert(ctx, &item)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create new survey")
	}
	return &srvy, nil
}

type insertResult struct {
	dynamoKeyFields
	Result
}

// AddResult will add a result set to an existing survey, returning the ID of the
// new result record. The timestamp of it being added is stored into the record.
func AddResult(ctx context.Context, surveyID string, result Result) (string, error) {
	if surveyID == "" {
		return "", errors.New("must supply non-empty surveyID")
	}

	u := shortuuid.New()

	result.Time = time.Now()
	result.ResultID = u

	item := insertResult{
		dynamoKeyFields: dynamoKeyFields{
			PK: surveyPrefix + surveyID,
			SK: resultPrefix + u,
		},
		Result: result,
	}

	err := insert(ctx, &item)
	if err != nil {
		return "", errors.Wrap(err, "failed to add result record")
	}
	return u, nil
}

func insert(ctx context.Context, item interface{}) error {
	s := getSession() // TODO

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(surveyTableName),
	}

	_, err = s.PutItemWithContext(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// Get fetches the info and results for a survey by ID
func Get(ctx context.Context, id string) (*Survey, error) {
	s := getSession() // TODO

	result, err := s.GetItemWithContext(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(surveyTableName),
		Key: map[string]*dynamodb.AttributeValue{
			"PK": {
				// S: aws.String(surveyPrefix + id),
				S: aws.String("SURVEY"),
			},
			"SK": {
				S: aws.String(surveyPrefix + id),
				// S: aws.String("SURVEY"),
			},
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get survey '%s'", id)
	}

	// Fetch the parent survey record
	srvy := &Survey{}
	err = dynamodbattribute.UnmarshalMap(result.Item, srvy)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get survey '%s'", id)
	}

	// Fetch all the current results for the given survey
	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with( SK, :sk )"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String(surveyPrefix + id),
			},
			":sk": {
				S: aws.String(resultPrefix),
			},
		},
		TableName: aws.String(surveyTableName),
	}

	queryResult, err := s.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to return results for survey '%s'", id)
	}

	for _, r := range queryResult.Items {
		var surveyResult Result
		err := dynamodbattribute.UnmarshalMap(r, &surveyResult)
		if err != nil {
			return nil, errors.Wrap(err, "error reading survey results records")
		}
		srvy.Results = append(srvy.Results, surveyResult)
	}
	return srvy, nil
}

// List returns all the surveys currently in the data store and "open". Optionally
// set `closed` to `true` to include all surveys.
func List(ctx context.Context, closed bool) ([]*Survey, error) {
	s := getSession() // TODO

	queryInput := &dynamodb.QueryInput{
		KeyConditionExpression: aws.String("PK = :pk AND begins_with( SK, :sk )"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":pk": {
				S: aws.String("SURVEY"),
			},
			":sk": {
				S: aws.String(surveyPrefix),
			},
		},
		TableName: aws.String(surveyTableName),
	}

	queryResult, err := s.QueryWithContext(ctx, queryInput)
	if err != nil {
		return nil, errors.Wrap(err, "failed to return survey list")
	}

	surveys := []*Survey{}

	for _, r := range queryResult.Items {
		var srvySummary Survey
		err := dynamodbattribute.UnmarshalMap(r, &srvySummary)
		if err != nil {
			return nil, errors.Wrap(err, "error reading survey record")
		}

		// Fetch the full survey details
		srvy, err := Get(ctx, srvySummary.SurveyID)
		if err != nil {
			return nil, errors.Wrapf(err, "error fetching survey detail for survey '%s'", srvySummary.SurveyID)
		}

		surveys = append(surveys, srvy)
	}
	return surveys, nil
}
