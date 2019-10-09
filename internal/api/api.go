package api

import (
	"github.com/ONSdigital/team-survey/internal/stats"
	"github.com/ONSdigital/team-survey/internal/survey"

	"github.com/jinzhu/gorm"
)

// API is the wrapper for interaction with the Team Survey
type API struct {
	DB *gorm.DB
}

// SaveSurvey will save a new Survey if it doesn't exist
func (a API) SaveSurvey(s survey.Survey) (survey.Survey, error) {
	return s.Save(a.DB)
}

// GetStats will return survey stats
func (a API) GetStats(s survey.Survey) (stats.CalculatedResult, error) {
	return stats.Calculate(s)
}

// GetSurvey get a survey
func (a API) GetSurvey(id uint) (survey.Survey, error) {
	return survey.Survey{}.Get(a.DB, id)
}

// GetSurveyByName get a survey by name
func (a API) GetSurveyByName(name string) (survey.Survey, error) {
	return survey.Survey{}.GetByName(a.DB, name)
}

// GetSurveyByNameAndShareCode get a survey by name and share code
func (a API) GetSurveyByNameAndShareCode(name string, shareCode string) (survey.Survey, error) {
	return survey.Survey{}.GetByNameAndShareCode(a.DB, name, shareCode)
}
