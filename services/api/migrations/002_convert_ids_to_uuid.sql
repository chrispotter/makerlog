-- +goose Up
-- +goose StatementBegin
-- WARNING: This migration is intended for fresh database schemas or when existing data can be lost.
-- Converting from auto-incrementing integers to UUIDs will generate new IDs for all existing records.
-- Foreign key relationships will need to be re-established after this migration for existing data.
-- For production databases with existing data, consider creating a custom migration that preserves
-- relationships by creating temporary mapping tables.

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Convert users table
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey CASCADE;
ALTER TABLE users ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
UPDATE users SET id_new = uuid_generate_v4();
ALTER TABLE users DROP COLUMN id;
ALTER TABLE users RENAME COLUMN id_new TO id;
ALTER TABLE users ADD PRIMARY KEY (id);

-- Convert projects table
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_pkey CASCADE;
ALTER TABLE projects ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE projects ADD COLUMN user_id_new UUID;
UPDATE projects SET id_new = uuid_generate_v4();
-- WARNING: user_id_new will remain NULL, breaking referential integrity for existing data
-- In a production environment, you would need to map old user IDs to new UUIDs
ALTER TABLE projects DROP COLUMN id;
ALTER TABLE projects DROP COLUMN user_id;
ALTER TABLE projects RENAME COLUMN id_new TO id;
ALTER TABLE projects RENAME COLUMN user_id_new TO user_id;
ALTER TABLE projects ADD PRIMARY KEY (id);
ALTER TABLE projects ADD CONSTRAINT fk_projects_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
CREATE INDEX idx_projects_user_id ON projects(user_id);

-- Convert tasks table
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS tasks_pkey CASCADE;
ALTER TABLE tasks ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE tasks ADD COLUMN user_id_new UUID;
ALTER TABLE tasks ADD COLUMN project_id_new UUID;
UPDATE tasks SET id_new = uuid_generate_v4();
-- WARNING: user_id_new and project_id_new will remain NULL, breaking referential integrity for existing data
-- In a production environment, you would need to map old IDs to new UUIDs
ALTER TABLE tasks DROP COLUMN id;
ALTER TABLE tasks DROP COLUMN user_id;
ALTER TABLE tasks DROP COLUMN project_id;
ALTER TABLE tasks RENAME COLUMN id_new TO id;
ALTER TABLE tasks RENAME COLUMN user_id_new TO user_id;
ALTER TABLE tasks RENAME COLUMN project_id_new TO project_id;
ALTER TABLE tasks ADD PRIMARY KEY (id);
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE tasks ADD CONSTRAINT fk_tasks_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE;
CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);

-- Convert log_entries table
ALTER TABLE log_entries DROP CONSTRAINT IF EXISTS log_entries_pkey CASCADE;
ALTER TABLE log_entries ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE log_entries ADD COLUMN user_id_new UUID;
ALTER TABLE log_entries ADD COLUMN task_id_new UUID;
ALTER TABLE log_entries ADD COLUMN project_id_new UUID;
UPDATE log_entries SET id_new = uuid_generate_v4();
-- WARNING: Foreign key columns will remain NULL, breaking referential integrity for existing data
-- In a production environment, you would need to map old IDs to new UUIDs
ALTER TABLE log_entries DROP COLUMN id;
ALTER TABLE log_entries DROP COLUMN user_id;
ALTER TABLE log_entries DROP COLUMN task_id;
ALTER TABLE log_entries DROP COLUMN project_id;
ALTER TABLE log_entries RENAME COLUMN id_new TO id;
ALTER TABLE log_entries RENAME COLUMN user_id_new TO user_id;
ALTER TABLE log_entries RENAME COLUMN task_id_new TO task_id;
ALTER TABLE log_entries RENAME COLUMN project_id_new TO project_id;
ALTER TABLE log_entries ADD PRIMARY KEY (id);
ALTER TABLE log_entries ADD CONSTRAINT fk_log_entries_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE log_entries ADD CONSTRAINT fk_log_entries_task_id FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE SET NULL;
ALTER TABLE log_entries ADD CONSTRAINT fk_log_entries_project_id FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE SET NULL;
CREATE INDEX idx_log_entries_user_id ON log_entries(user_id);
CREATE INDEX idx_log_entries_task_id ON log_entries(task_id);
CREATE INDEX idx_log_entries_project_id ON log_entries(project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Note: This down migration will result in data loss as we can't convert UUIDs back to sequential integers
-- This is a destructive change and should be carefully considered

DROP TABLE IF EXISTS log_entries;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

-- Recreate tables with original schema
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);

CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_user_id ON projects(user_id);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL DEFAULT 'todo',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_status ON tasks(status);

CREATE TABLE log_entries (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id INTEGER REFERENCES tasks(id) ON DELETE SET NULL,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    content TEXT NOT NULL,
    log_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_log_entries_user_id ON log_entries(user_id);
CREATE INDEX idx_log_entries_task_id ON log_entries(task_id);
CREATE INDEX idx_log_entries_project_id ON log_entries(project_id);
CREATE INDEX idx_log_entries_log_date ON log_entries(log_date);
-- +goose StatementEnd
