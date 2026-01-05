package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUserJSONSerialization(t *testing.T) {
	user := User{
		ID:           "test-id",
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
		Name:         "Test User",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test that password hash is not included in JSON
	jsonData, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, exists := result["password_hash"]; exists {
		t.Error("Password hash should not be included in JSON output")
	}

	if result["id"] != user.ID {
		t.Errorf("Expected id %s, got %v", user.ID, result["id"])
	}

	if result["email"] != user.Email {
		t.Errorf("Expected email %s, got %v", user.Email, result["email"])
	}
}

func TestProjectJSONSerialization(t *testing.T) {
	project := Project{
		ID:          "project-id",
		UserID:      "user-id",
		Name:        "Test Project",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		t.Fatalf("Failed to marshal project: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["name"] != project.Name {
		t.Errorf("Expected name %s, got %v", project.Name, result["name"])
	}
}

func TestTaskJSONSerialization(t *testing.T) {
	task := Task{
		ID:          "task-id",
		ProjectID:   "project-id",
		UserID:      "user-id",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "todo",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["title"] != task.Title {
		t.Errorf("Expected title %s, got %v", task.Title, result["title"])
	}

	if result["status"] != task.Status {
		t.Errorf("Expected status %s, got %v", task.Status, result["status"])
	}
}

func TestLogEntryJSONSerialization(t *testing.T) {
	taskID := "task-id"
	projectID := "project-id"
	logEntry := LogEntry{
		ID:        "log-id",
		UserID:    "user-id",
		TaskID:    &taskID,
		ProjectID: &projectID,
		Content:   "Test log entry",
		LogDate:   time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		t.Fatalf("Failed to marshal log entry: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if result["content"] != logEntry.Content {
		t.Errorf("Expected content %s, got %v", logEntry.Content, result["content"])
	}

	if result["task_id"] == nil {
		t.Error("Expected task_id to be present")
	}
}

func TestRegisterRequestDeserialization(t *testing.T) {
	jsonStr := `{"email":"test@example.com","password":"password123","name":"Test User"}`
	var req RegisterRequest

	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal register request: %v", err)
	}

	if req.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", req.Email)
	}

	if req.Password != "password123" {
		t.Errorf("Expected password password123, got %s", req.Password)
	}

	if req.Name != "Test User" {
		t.Errorf("Expected name Test User, got %s", req.Name)
	}
}

func TestLoginRequestDeserialization(t *testing.T) {
	jsonStr := `{"email":"test@example.com","password":"password123"}`
	var req LoginRequest

	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal login request: %v", err)
	}

	if req.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", req.Email)
	}

	if req.Password != "password123" {
		t.Errorf("Expected password password123, got %s", req.Password)
	}
}

func TestCreateProjectRequestDeserialization(t *testing.T) {
	jsonStr := `{"name":"Test Project","description":"Test Description"}`
	var req CreateProjectRequest

	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal create project request: %v", err)
	}

	if req.Name != "Test Project" {
		t.Errorf("Expected name Test Project, got %s", req.Name)
	}

	if req.Description != "Test Description" {
		t.Errorf("Expected description Test Description, got %s", req.Description)
	}
}

func TestCreateTaskRequestDeserialization(t *testing.T) {
	jsonStr := `{"project_id":"proj-123","title":"Test Task","description":"Test Description","status":"todo"}`
	var req CreateTaskRequest

	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal create task request: %v", err)
	}

	if req.ProjectID != "proj-123" {
		t.Errorf("Expected project_id proj-123, got %s", req.ProjectID)
	}

	if req.Title != "Test Task" {
		t.Errorf("Expected title Test Task, got %s", req.Title)
	}

	if req.Status != "todo" {
		t.Errorf("Expected status todo, got %s", req.Status)
	}
}

func TestCreateLogEntryRequestDeserialization(t *testing.T) {
	jsonStr := `{"task_id":"task-123","project_id":"proj-123","content":"Log content","log_date":"2024-01-05"}`
	var req CreateLogEntryRequest

	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		t.Fatalf("Failed to unmarshal create log entry request: %v", err)
	}

	if req.TaskID == nil || *req.TaskID != "task-123" {
		t.Errorf("Expected task_id task-123, got %v", req.TaskID)
	}

	if req.ProjectID == nil || *req.ProjectID != "proj-123" {
		t.Errorf("Expected project_id proj-123, got %v", req.ProjectID)
	}

	if req.Content != "Log content" {
		t.Errorf("Expected content Log content, got %s", req.Content)
	}

	if req.LogDate != "2024-01-05" {
		t.Errorf("Expected log_date 2024-01-05, got %s", req.LogDate)
	}
}
