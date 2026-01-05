package handlers

import (
	"testing"
	"time"

	"github.com/chrispotter/makerlog/services/api/internal/models"
	"github.com/google/uuid"
)

func TestNewLogEntryHandler(t *testing.T) {
	handler := NewLogEntryHandler(nil)
	if handler == nil {
		t.Error("Expected handler to be created")
		return
	}
	if handler.queries != nil {
		t.Error("Expected queries to be nil when passed nil")
	}
}

func TestCreateLogEntryRequestValidation(t *testing.T) {
	taskID := "550e8400-e29b-41d4-a716-446655440000"
	projectID := "660e8400-e29b-41d4-a716-446655440000"
	invalidID := "invalid-uuid"

	tests := []struct {
		name    string
		request models.CreateLogEntryRequest
		isValid bool
	}{
		{
			name: "valid request with task and project",
			request: models.CreateLogEntryRequest{
				TaskID:    &taskID,
				ProjectID: &projectID,
				Content:   "Test log entry",
				LogDate:   "2024-01-05",
			},
			isValid: true,
		},
		{
			name: "valid request with only content",
			request: models.CreateLogEntryRequest{
				Content: "Test log entry",
				LogDate: "2024-01-05",
			},
			isValid: true,
		},
		{
			name: "empty content",
			request: models.CreateLogEntryRequest{
				Content: "",
				LogDate: "2024-01-05",
			},
			isValid: false,
		},
		{
			name: "invalid task ID",
			request: models.CreateLogEntryRequest{
				TaskID:  &invalidID,
				Content: "Test log entry",
				LogDate: "2024-01-05",
			},
			isValid: false,
		},
		{
			name: "invalid project ID",
			request: models.CreateLogEntryRequest{
				ProjectID: &invalidID,
				Content:   "Test log entry",
				LogDate:   "2024-01-05",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contentValid := tt.request.Content != ""

			taskIDValid := true
			if tt.request.TaskID != nil {
				_, err := uuid.Parse(*tt.request.TaskID)
				taskIDValid = err == nil
			}

			projectIDValid := true
			if tt.request.ProjectID != nil {
				_, err := uuid.Parse(*tt.request.ProjectID)
				projectIDValid = err == nil
			}

			actualValid := contentValid && taskIDValid && projectIDValid
			if actualValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}

func TestLogDateParsing(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		isValid bool
	}{
		{
			name:    "valid date format",
			date:    "2024-01-05",
			isValid: true,
		},
		{
			name:    "invalid date format",
			date:    "01/05/2024",
			isValid: false,
		},
		{
			name:    "invalid date",
			date:    "2024-13-45",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse("2006-01-02", tt.date)
			if (err == nil) != tt.isValid {
				t.Errorf("Expected isValid=%v for date=%s, err=%v", tt.isValid, tt.date, err)
			}
		})
	}
}

func TestUpdateLogEntryRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.UpdateLogEntryRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.UpdateLogEntryRequest{
				Content: "Updated log entry",
				LogDate: "2024-01-05",
			},
			isValid: true,
		},
		{
			name: "empty content",
			request: models.UpdateLogEntryRequest{
				Content: "",
				LogDate: "2024-01-05",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.request.Content != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}
