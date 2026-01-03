export interface User {
  id: number;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Project {
  id: number;
  user_id: number;
  name: string;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: number;
  user_id: number;
  project_id: number;
  title: string;
  description: string;
  status: 'todo' | 'in_progress' | 'done';
  created_at: string;
  updated_at: string;
}

export interface LogEntry {
  id: number;
  user_id: number;
  task_id?: number;
  project_id?: number;
  content: string;
  log_date: string;
  created_at: string;
  updated_at: string;
}

export interface RegisterData {
  email: string;
  password: string;
  name: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface CreateProjectData {
  name: string;
  description: string;
}

export interface CreateTaskData {
  project_id: number;
  title: string;
  description: string;
  status?: string;
}

export interface CreateLogEntryData {
  task_id?: number;
  project_id?: number;
  content: string;
  log_date?: string;
}
