package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		hasValue bool
	}{
		{
			name:     "valid user ID",
			userID:   "user-123",
			hasValue: true,
		},
		{
			name:     "empty user ID",
			userID:   "",
			hasValue: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), userIDKey, tt.userID)
			userID, ok := GetUserID(ctx)

			if !ok {
				t.Error("Expected ok to be true")
			}

			if userID != tt.userID {
				t.Errorf("Expected userID %s, got %s", tt.userID, userID)
			}
		})
	}
}

func TestGetUserIDNoValue(t *testing.T) {
	ctx := context.Background()
	userID, ok := GetUserID(ctx)

	if ok {
		t.Error("Expected ok to be false when no user ID in context")
	}

	if userID != "" {
		t.Errorf("Expected empty userID, got %s", userID)
	}
}

func TestAuthMiddleware(t *testing.T) {
	sessionStore := sessions.NewCookieStore([]byte("test-secret-key"))

	tests := []struct {
		name           string
		setupSession   func(*http.Request, http.ResponseWriter) error
		expectedStatus int
	}{
		{
			name: "valid session with user ID",
			setupSession: func(r *http.Request, w http.ResponseWriter) error {
				session, err := sessionStore.Get(r, "makerlog-session")
				if err != nil {
					return err
				}
				session.Values["user_id"] = "user-123"
				return session.Save(r, w)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "no session",
			setupSession: func(r *http.Request, w http.ResponseWriter) error {
				return nil
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "session without user ID",
			setupSession: func(r *http.Request, w http.ResponseWriter) error {
				session, err := sessionStore.Get(r, "makerlog-session")
				if err != nil {
					return err
				}
				return session.Save(r, w)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "session with empty user ID",
			setupSession: func(r *http.Request, w http.ResponseWriter) error {
				session, err := sessionStore.Get(r, "makerlog-session")
				if err != nil {
					return err
				}
				session.Values["user_id"] = ""
				return session.Save(r, w)
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify user ID is in context for successful cases
				if tt.expectedStatus == http.StatusOK {
					userID, ok := GetUserID(r.Context())
					if !ok {
						t.Error("Expected user ID in context")
					}
					if userID != "user-123" {
						t.Errorf("Expected userID user-123, got %s", userID)
					}
				}
				w.WriteHeader(http.StatusOK)
			})

			// Create auth middleware
			authMiddleware := Auth(sessionStore)
			handler := authMiddleware(nextHandler)

			// Create request and response recorder
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Setup session if needed
			if tt.setupSession != nil {
				if err := tt.setupSession(req, w); err != nil {
					t.Fatalf("Failed to setup session: %v", err)
				}
				// Copy cookies from setup response to request
				for _, cookie := range w.Result().Cookies() {
					req.AddCookie(cookie)
				}
				// Reset response recorder
				w = httptest.NewRecorder()
			}

			// Execute handler
			handler.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
