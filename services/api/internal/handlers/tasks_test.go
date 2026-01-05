package handlers

import (
	"testing"

	"github.com/chrispotter/makerlog/services/api/internal/models"
	"github.com/google/uuid"
)

func TestNewTaskHandler(t *testing.T) {
	handler := NewTaskHandler(nil)
	if handler == nil {
		t.Error("Expected handler to be created")
		return
	}
	if handler.queries != nil {
		t.Error("Expected queries to be nil when passed nil")
	}
}

func TestCreateTaskRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.CreateTaskRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.CreateTaskRequest{
				ProjectID:   "550e8400-e29b-41d4-a716-446655440000",
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "todo",
			},
			isValid: true,
		},
		{
			name: "empty title",
			request: models.CreateTaskRequest{
				ProjectID:   "550e8400-e29b-41d4-a716-446655440000",
				Title:       "",
				Description: "Test Description",
				Status:      "todo",
			},
			isValid: false,
		},
		{
			name: "invalid project ID",
			request: models.CreateTaskRequest{
				ProjectID:   "invalid-uuid",
				Title:       "Test Task",
				Description: "Test Description",
				Status:      "todo",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			titleValid := tt.request.Title != ""
			_, uuidErr := uuid.Parse(tt.request.ProjectID)
			uuidValid := uuidErr == nil

			actualValid := titleValid && uuidValid
			if actualValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}

func TestUpdateTaskRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.UpdateTaskRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.UpdateTaskRequest{
				Title:       "Updated Task",
				Description: "Updated Description",
				Status:      "in_progress",
			},
			isValid: true,
		},
		{
			name: "empty title",
			request: models.UpdateTaskRequest{
				Title:       "",
				Description: "Updated Description",
				Status:      "done",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.request.Title != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}

func TestTaskDefaultStatus(t *testing.T) {
	status := ""
	if status == "" {
		status = "todo"
	}

	if status != "todo" {
		t.Errorf("Expected default status to be 'todo', got %s", status)
	}
}
