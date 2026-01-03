import type {
  User,
  Project,
  Task,
  LogEntry,
  RegisterData,
  LoginData,
  CreateProjectData,
  CreateTaskData,
  CreateLogEntryData,
} from './types';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const config: RequestInit = {
      ...options,
      credentials: 'include', // Important for cookie-based sessions
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    };

    const response = await fetch(url, config);

    if (!response.ok) {
      const error = await response.text();
      throw new Error(error || `HTTP error! status: ${response.status}`);
    }

    // Handle empty responses (like 204 No Content)
    if (response.status === 204) {
      return {} as T;
    }

    return response.json();
  }

  // Auth endpoints
  async register(data: RegisterData): Promise<User> {
    return this.request<User>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async login(data: LoginData): Promise<User> {
    return this.request<User>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async logout(): Promise<void> {
    return this.request<void>('/api/auth/logout', {
      method: 'POST',
    });
  }

  async getMe(): Promise<User> {
    return this.request<User>('/api/auth/me');
  }

  // Projects endpoints
  async getProjects(): Promise<Project[]> {
    return this.request<Project[]>('/api/projects');
  }

  async getProject(id: number): Promise<Project> {
    return this.request<Project>(`/api/projects/${id}`);
  }

  async createProject(data: CreateProjectData): Promise<Project> {
    return this.request<Project>('/api/projects', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateProject(id: number, data: Partial<CreateProjectData>): Promise<Project> {
    return this.request<Project>(`/api/projects/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteProject(id: number): Promise<void> {
    return this.request<void>(`/api/projects/${id}`, {
      method: 'DELETE',
    });
  }

  // Tasks endpoints
  async getTasks(projectId?: number): Promise<Task[]> {
    const query = projectId ? `?project_id=${projectId}` : '';
    return this.request<Task[]>(`/api/tasks${query}`);
  }

  async getTask(id: number): Promise<Task> {
    return this.request<Task>(`/api/tasks/${id}`);
  }

  async createTask(data: CreateTaskData): Promise<Task> {
    return this.request<Task>('/api/tasks', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateTask(id: number, data: Partial<CreateTaskData>): Promise<Task> {
    return this.request<Task>(`/api/tasks/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteTask(id: number): Promise<void> {
    return this.request<void>(`/api/tasks/${id}`, {
      method: 'DELETE',
    });
  }

  // Log entries endpoints
  async getLogEntries(projectId?: number): Promise<LogEntry[]> {
    const query = projectId ? `?project_id=${projectId}` : '';
    return this.request<LogEntry[]>(`/api/log-entries${query}`);
  }

  async getLogEntry(id: number): Promise<LogEntry> {
    return this.request<LogEntry>(`/api/log-entries/${id}`);
  }

  async createLogEntry(data: CreateLogEntryData): Promise<LogEntry> {
    return this.request<LogEntry>('/api/log-entries', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  }

  async updateLogEntry(id: number, data: Partial<CreateLogEntryData>): Promise<LogEntry> {
    return this.request<LogEntry>(`/api/log-entries/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async deleteLogEntry(id: number): Promise<void> {
    return this.request<void>(`/api/log-entries/${id}`, {
      method: 'DELETE',
    });
  }

  // Today endpoint
  async getTodayLogEntries(): Promise<LogEntry[]> {
    return this.request<LogEntry[]>('/api/today');
  }
}

export const apiClient = new ApiClient();
