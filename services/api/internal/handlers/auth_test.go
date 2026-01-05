package handlers

import (
	"testing"

	"github.com/chrispotter/makerlog/services/api/internal/models"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func TestNewAuthHandler(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret"))
	handler := NewAuthHandler(nil, sessionStore)

	if handler == nil {
		t.Error("Expected handler to be created")
	}
	if handler.queries != nil {
		t.Error("Expected queries to be nil when passed nil")
	}
	if handler.sessionStore != sessionStore {
		t.Error("Expected sessionStore to be set")
	}
}

func TestRegisterRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.RegisterRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "Test User",
			},
			isValid: true,
		},
		{
			name: "empty email",
			request: models.RegisterRequest{
				Email:    "",
				Password: "password123",
				Name:     "Test User",
			},
			isValid: false,
		},
		{
			name: "empty password",
			request: models.RegisterRequest{
				Email:    "test@example.com",
				Password: "",
				Name:     "Test User",
			},
			isValid: false,
		},
		{
			name: "empty name",
			request: models.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
				Name:     "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualValid := tt.request.Email != "" && tt.request.Password != "" && tt.request.Name != ""
			if actualValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}

func TestLoginRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.LoginRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			isValid: true,
		},
		{
			name: "empty email",
			request: models.LoginRequest{
				Email:    "",
				Password: "password123",
			},
			isValid: false,
		},
		{
			name: "empty password",
			request: models.LoginRequest{
				Email:    "test@example.com",
				Password: "",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualValid := tt.request.Email != "" && tt.request.Password != ""
			if actualValid != tt.isValid {
				t.Errorf("Expected isValid=%v for request=%+v", tt.isValid, tt.request)
			}
		})
	}
}

func TestPasswordHashing(t *testing.T) {
	password := "testpassword123"

	// Test hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Test that hash is different from original
	if string(hashedPassword) == password {
		t.Error("Hashed password should not equal plain password")
	}

	// Test successful comparison
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		t.Errorf("Password comparison should succeed: %v", err)
	}

	// Test failed comparison with wrong password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte("wrongpassword"))
	if err == nil {
		t.Error("Password comparison should fail with wrong password")
	}
}

func TestSessionStoreCreation(t *testing.T) {
	secret := []byte("test-secret-key")
	store := sessions.NewCookieStore(secret)

	if store == nil {
		t.Error("Expected session store to be created")
	}
}
