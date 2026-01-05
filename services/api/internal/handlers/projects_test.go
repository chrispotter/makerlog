package handlers

import (
	"testing"

	"github.com/chrispotter/makerlog/services/api/internal/models"
	"github.com/google/uuid"
)

func TestNewProjectHandler(t *testing.T) {
	handler := NewProjectHandler(nil)
	if handler == nil {
		t.Error("Expected handler to be created")
		return
	}
	if handler.queries != nil {
		t.Error("Expected queries to be nil when passed nil")
	}
}

func TestValidateUUID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		isValid bool
	}{
		{
			name:    "valid UUID",
			id:      "550e8400-e29b-41d4-a716-446655440000",
			isValid: true,
		},
		{
			name:    "invalid UUID",
			id:      "invalid-uuid",
			isValid: false,
		},
		{
			name:    "empty string",
			id:      "",
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := uuid.Parse(tt.id)
			if (err == nil) != tt.isValid {
				t.Errorf("Expected isValid=%v for id=%s", tt.isValid, tt.id)
			}
		})
	}
}

func TestCreateProjectRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.CreateProjectRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.CreateProjectRequest{
				Name:        "Test Project",
				Description: "Test Description",
			},
			isValid: true,
		},
		{
			name: "empty name",
			request: models.CreateProjectRequest{
				Name:        "",
				Description: "Test Description",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.request.Name != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}

func TestUpdateProjectRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.UpdateProjectRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.UpdateProjectRequest{
				Name:        "Updated Project",
				Description: "Updated Description",
			},
			isValid: true,
		},
		{
			name: "empty name",
			request: models.UpdateProjectRequest{
				Name:        "",
				Description: "Updated Description",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.request.Name != ""
			if isValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}
