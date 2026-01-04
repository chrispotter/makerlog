package database

import (
	"database/sql"
	"time"

	"github.com/chrispotter/makerlog/services/api/internal/models"
	"github.com/google/uuid"
)

type Queries struct {
	db *sql.DB
}

func New(db *sql.DB) *Queries {
	return &Queries{db: db}
}

// User queries
func (q *Queries) CreateUser(email, passwordHash, name string) (*models.User, error) {
	var user models.User
	id := uuid.New().String()
	err := q.db.QueryRow(`
		INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, email, password_hash, name, created_at, updated_at
	`, id, email, passwordHash, name).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	return &user, err
}

func (q *Queries) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := q.db.QueryRow(`
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (q *Queries) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := q.db.QueryRow(`
		SELECT id, email, password_hash, name, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

// Project queries
func (q *Queries) CreateProject(userID string, name, description string) (*models.Project, error) {
	var project models.Project
	id := uuid.New().String()
	err := q.db.QueryRow(`
		INSERT INTO projects (id, user_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, user_id, name, description, created_at, updated_at
	`, id, userID, name, description).Scan(
		&project.ID, &project.UserID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt,
	)
	return &project, err
}

func (q *Queries) GetProject(id, userID string) (*models.Project, error) {
	var project models.Project
	err := q.db.QueryRow(`
		SELECT id, user_id, name, description, created_at, updated_at
		FROM projects WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&project.ID, &project.UserID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &project, err
}

func (q *Queries) ListProjects(userID string) ([]models.Project, error) {
	rows, err := q.db.Query(`
		SELECT id, user_id, name, description, created_at, updated_at
		FROM projects WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(
			&project.ID, &project.UserID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt,
		); err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}
	return projects, rows.Err()
}

func (q *Queries) UpdateProject(id, userID string, name, description string) (*models.Project, error) {
	var project models.Project
	err := q.db.QueryRow(`
		UPDATE projects
		SET name = $1, description = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, description, created_at, updated_at
	`, name, description, id, userID).Scan(
		&project.ID, &project.UserID, &project.Name, &project.Description, &project.CreatedAt, &project.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &project, err
}

func (q *Queries) DeleteProject(id, userID string) error {
	result, err := q.db.Exec(`
		DELETE FROM projects WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Task queries
func (q *Queries) CreateTask(userID, projectID string, title, description, status string) (*models.Task, error) {
	var task models.Task
	id := uuid.New().String()
	err := q.db.QueryRow(`
		INSERT INTO tasks (id, user_id, project_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, user_id, project_id, title, description, status, created_at, updated_at
	`, id, userID, projectID, title, description, status).Scan(
		&task.ID, &task.UserID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	return &task, err
}

func (q *Queries) GetTask(id, userID string) (*models.Task, error) {
	var task models.Task
	err := q.db.QueryRow(`
		SELECT id, user_id, project_id, title, description, status, created_at, updated_at
		FROM tasks WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&task.ID, &task.UserID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (q *Queries) ListTasks(userID string, projectID *string) ([]models.Task, error) {
	var rows *sql.Rows
	var err error

	if projectID != nil {
		rows, err = q.db.Query(`
			SELECT id, user_id, project_id, title, description, status, created_at, updated_at
			FROM tasks WHERE user_id = $1 AND project_id = $2
			ORDER BY created_at DESC
		`, userID, *projectID)
	} else {
		rows, err = q.db.Query(`
			SELECT id, user_id, project_id, title, description, status, created_at, updated_at
			FROM tasks WHERE user_id = $1
			ORDER BY created_at DESC
		`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(
			&task.ID, &task.UserID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func (q *Queries) UpdateTask(id, userID string, title, description, status string) (*models.Task, error) {
	var task models.Task
	err := q.db.QueryRow(`
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = NOW()
		WHERE id = $4 AND user_id = $5
		RETURNING id, user_id, project_id, title, description, status, created_at, updated_at
	`, title, description, status, id, userID).Scan(
		&task.ID, &task.UserID, &task.ProjectID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (q *Queries) DeleteTask(id, userID string) error {
	result, err := q.db.Exec(`
		DELETE FROM tasks WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Log entry queries
func (q *Queries) CreateLogEntry(userID string, taskID, projectID *string, content string, logDate time.Time) (*models.LogEntry, error) {
	var logEntry models.LogEntry
	id := uuid.New().String()
	err := q.db.QueryRow(`
		INSERT INTO log_entries (id, user_id, task_id, project_id, content, log_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, user_id, task_id, project_id, content, log_date, created_at, updated_at
	`, id, userID, taskID, projectID, content, logDate).Scan(
		&logEntry.ID, &logEntry.UserID, &logEntry.TaskID, &logEntry.ProjectID, &logEntry.Content, &logEntry.LogDate, &logEntry.CreatedAt, &logEntry.UpdatedAt,
	)
	return &logEntry, err
}

func (q *Queries) GetLogEntry(id, userID string) (*models.LogEntry, error) {
	var logEntry models.LogEntry
	err := q.db.QueryRow(`
		SELECT id, user_id, task_id, project_id, content, log_date, created_at, updated_at
		FROM log_entries WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&logEntry.ID, &logEntry.UserID, &logEntry.TaskID, &logEntry.ProjectID, &logEntry.Content, &logEntry.LogDate, &logEntry.CreatedAt, &logEntry.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &logEntry, err
}

func (q *Queries) ListLogEntries(userID string, projectID *string) ([]models.LogEntry, error) {
	var rows *sql.Rows
	var err error

	if projectID != nil {
		rows, err = q.db.Query(`
			SELECT id, user_id, task_id, project_id, content, log_date, created_at, updated_at
			FROM log_entries WHERE user_id = $1 AND project_id = $2
			ORDER BY log_date DESC, created_at DESC
		`, userID, *projectID)
	} else {
		rows, err = q.db.Query(`
			SELECT id, user_id, task_id, project_id, content, log_date, created_at, updated_at
			FROM log_entries WHERE user_id = $1
			ORDER BY log_date DESC, created_at DESC
		`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logEntries []models.LogEntry
	for rows.Next() {
		var logEntry models.LogEntry
		if err := rows.Scan(
			&logEntry.ID, &logEntry.UserID, &logEntry.TaskID, &logEntry.ProjectID, &logEntry.Content, &logEntry.LogDate, &logEntry.CreatedAt, &logEntry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		logEntries = append(logEntries, logEntry)
	}
	return logEntries, rows.Err()
}

func (q *Queries) GetTodayLogEntries(userID string, date time.Time) ([]models.LogEntry, error) {
	rows, err := q.db.Query(`
		SELECT id, user_id, task_id, project_id, content, log_date, created_at, updated_at
		FROM log_entries
		WHERE user_id = $1 AND DATE(log_date) = DATE($2)
		ORDER BY created_at DESC
	`, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logEntries []models.LogEntry
	for rows.Next() {
		var logEntry models.LogEntry
		if err := rows.Scan(
			&logEntry.ID, &logEntry.UserID, &logEntry.TaskID, &logEntry.ProjectID, &logEntry.Content, &logEntry.LogDate, &logEntry.CreatedAt, &logEntry.UpdatedAt,
		); err != nil {
			return nil, err
		}
		logEntries = append(logEntries, logEntry)
	}
	return logEntries, rows.Err()
}

func (q *Queries) UpdateLogEntry(id, userID string, content string, logDate time.Time) (*models.LogEntry, error) {
	var logEntry models.LogEntry
	err := q.db.QueryRow(`
		UPDATE log_entries
		SET content = $1, log_date = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, task_id, project_id, content, log_date, created_at, updated_at
	`, content, logDate, id, userID).Scan(
		&logEntry.ID, &logEntry.UserID, &logEntry.TaskID, &logEntry.ProjectID, &logEntry.Content, &logEntry.LogDate, &logEntry.CreatedAt, &logEntry.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &logEntry, err
}

func (q *Queries) DeleteLogEntry(id, userID string) error {
	result, err := q.db.Exec(`
		DELETE FROM log_entries WHERE id = $1 AND user_id = $2
	`, id, userID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
