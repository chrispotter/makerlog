package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chrispotter/makerlog/services/api/internal/database"
	"github.com/chrispotter/makerlog/services/api/internal/middleware"
	"github.com/chrispotter/makerlog/services/api/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LogEntryHandler struct {
	queries *database.Queries
}

func NewLogEntryHandler(queries *database.Queries) *LogEntryHandler {
	return &LogEntryHandler{queries: queries}
}

func (h *LogEntryHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Optional project_id filter
	var projectID *string
	if projectIDStr := r.URL.Query().Get("project_id"); projectIDStr != "" {
		if _, err := uuid.Parse(projectIDStr); err != nil {
			http.Error(w, "Invalid project_id format", http.StatusBadRequest)
			return
		}
		projectID = &projectIDStr
	}

	logEntries, err := h.queries.ListLogEntries(userID, projectID)
	if err != nil {
		http.Error(w, "Failed to list log entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntries)
}

func (h *LogEntryHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.CreateLogEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Validate optional UUID fields
	if req.TaskID != nil {
		if _, err := uuid.Parse(*req.TaskID); err != nil {
			http.Error(w, "Invalid task_id format", http.StatusBadRequest)
			return
		}
	}
	if req.ProjectID != nil {
		if _, err := uuid.Parse(*req.ProjectID); err != nil {
			http.Error(w, "Invalid project_id format", http.StatusBadRequest)
			return
		}
	}

	// Parse log date
	var logDate time.Time
	var err error
	if req.LogDate != "" {
		logDate, err = time.Parse("2006-01-02", req.LogDate)
		if err != nil {
			http.Error(w, "Invalid log date format. Use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		logDate = time.Now()
	}

	logEntry, err := h.queries.CreateLogEntry(userID, req.TaskID, req.ProjectID, req.Content, logDate)
	if err != nil {
		http.Error(w, "Failed to create log entry", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}

func (h *LogEntryHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid log entry ID format", http.StatusBadRequest)
		return
	}

	logEntry, err := h.queries.GetLogEntry(id, userID)
	if err != nil {
		http.Error(w, "Failed to get log entry", http.StatusInternalServerError)
		return
	}
	if logEntry == nil {
		http.Error(w, "Log entry not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

func (h *LogEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid log entry ID format", http.StatusBadRequest)
		return
	}

	var req models.UpdateLogEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	// Parse log date
	logDate, err := time.Parse("2006-01-02", req.LogDate)
	if err != nil {
		http.Error(w, "Invalid log date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	logEntry, err := h.queries.UpdateLogEntry(id, userID, req.Content, logDate)
	if err != nil {
		http.Error(w, "Failed to update log entry", http.StatusInternalServerError)
		return
	}
	if logEntry == nil {
		http.Error(w, "Log entry not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntry)
}

func (h *LogEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "Invalid log entry ID format", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteLogEntry(id, userID); err != nil {
		http.Error(w, "Failed to delete log entry", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *LogEntryHandler) Today(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get today's log entries
	today := time.Now()
	logEntries, err := h.queries.GetTodayLogEntries(userID, today)
	if err != nil {
		http.Error(w, "Failed to get today's log entries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logEntries)
}
