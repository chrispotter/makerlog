# Maker Log

A full-stack monorepo application for tracking daily progress, projects, tasks, and log entries.

## Architecture

This is a monorepo containing:

- **Backend API** (`/services/api`): Go 1.22 REST API with chi router, PostgreSQL, sqlx, cookie sessions, bcrypt authentication, and goose migrations
- **Frontend Web** (`/apps/web`): Next.js 14 app with TypeScript and Tailwind CSS

## Features

### Backend (`/services/api`)
- **Authentication**: Cookie-based sessions with bcrypt password hashing
- **CRUD Operations**:
  - Projects: Create, read, update, delete projects
  - Tasks: Manage tasks within projects with status tracking (todo, in_progress, done)
  - Log Entries: Track daily work logs linked to projects and tasks
- **Special Endpoints**:
  - `GET /api/today`: Retrieve today's log entries
- **Database**: PostgreSQL with goose migrations
- **Router**: Chi router with middleware support
- **CORS**: Configured for frontend communication

### Frontend (`/apps/web`)
- **Pages**:
  - `/`: Home page with project listing and authentication
  - `/projects/[id]`: Project detail page with tasks and log entries
  - `/today`: Today's log entries view
- **Features**:
  - User authentication (login/register)
  - Project management
  - Task tracking with status updates
  - Daily log entries
  - Responsive design with Tailwind CSS

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for local development)
- Node.js 20+ (for local development)
- PostgreSQL 16 (for local development without Docker)
- Make (optional, for using Makefile commands)

### Quick Start with Docker

1. Clone the repository:
```bash
git clone https://github.com/chrispotter/makerlog.git
cd makerlog
```

2. Start all services:
```bash
make up
# or
docker-compose up -d
```

3. Run database migrations:
```bash
make migrate-up
```

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

### Local Development (without Docker)

#### Backend API

1. Start PostgreSQL:
```bash
make db-start
```

2. Create a `.env` file in `services/api`:
```bash
cp services/api/.env.example services/api/.env
```

3. Run migrations:
```bash
make migrate-up
```

4. Install dependencies and run the API:
```bash
cd services/api
go mod download
go run cmd/api/main.go
```

The API will be available at http://localhost:8080

#### Frontend Web

1. Create a `.env.local` file in `apps/web`:
```bash
cp apps/web/.env.example apps/web/.env.local
```

2. Install dependencies and run the development server:
```bash
cd apps/web
npm install
npm run dev
```

The web app will be available at http://localhost:3000

## API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login
- `POST /api/auth/logout` - Logout
- `GET /api/auth/me` - Get current user

### Projects
- `GET /api/projects` - List all projects
- `POST /api/projects` - Create a project
- `GET /api/projects/:id` - Get a project
- `PUT /api/projects/:id` - Update a project
- `DELETE /api/projects/:id` - Delete a project

### Tasks
- `GET /api/tasks` - List all tasks (optional `?project_id=` filter)
- `POST /api/tasks` - Create a task
- `GET /api/tasks/:id` - Get a task
- `PUT /api/tasks/:id` - Update a task
- `DELETE /api/tasks/:id` - Delete a task

### Log Entries
- `GET /api/log-entries` - List all log entries (optional `?project_id=` filter)
- `POST /api/log-entries` - Create a log entry
- `GET /api/log-entries/:id` - Get a log entry
- `PUT /api/log-entries/:id` - Update a log entry
- `DELETE /api/log-entries/:id` - Delete a log entry
- `GET /api/today` - Get today's log entries

## Database Schema

### Users
- `id` (serial, primary key)
- `email` (varchar, unique)
- `password_hash` (varchar)
- `name` (varchar)
- `created_at`, `updated_at` (timestamp)

### Projects
- `id` (serial, primary key)
- `user_id` (foreign key → users)
- `name` (varchar)
- `description` (text)
- `created_at`, `updated_at` (timestamp)

### Tasks
- `id` (serial, primary key)
- `user_id` (foreign key → users)
- `project_id` (foreign key → projects)
- `title` (varchar)
- `description` (text)
- `status` (varchar: todo, in_progress, done)
- `created_at`, `updated_at` (timestamp)

### Log Entries
- `id` (serial, primary key)
- `user_id` (foreign key → users)
- `task_id` (foreign key → tasks, nullable)
- `project_id` (foreign key → projects, nullable)
- `content` (text)
- `log_date` (date)
- `created_at`, `updated_at` (timestamp)

## Makefile Commands

Run `make help` to see all available commands:

- `make install` - Install all dependencies
- `make build` - Build Docker containers
- `make up` - Start all services
- `make down` - Stop all services
- `make logs` - View logs from all services
- `make clean` - Stop services and remove volumes
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create name=<name>` - Create a new migration
- `make api-dev` - Run API in development mode
- `make web-dev` - Run web app in development mode
- `make db-shell` - Open PostgreSQL shell

## Environment Variables

### Backend (`services/api/.env`)
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/makerlog?sslmode=disable
SESSION_SECRET=your-secret-key-change-this-in-production
PORT=8080
FRONTEND_URL=http://localhost:3000
```

### Frontend (`apps/web/.env.local`)
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Technologies Used

### Backend
- Go 1.22
- Chi (router)
- PostgreSQL
- sqlx (database)
- Gorilla Sessions (cookie sessions)
- bcrypt (password hashing)
- goose (database migrations)

### Frontend
- Next.js 14
- TypeScript
- Tailwind CSS
- React

### Infrastructure
- Docker & Docker Compose
- Make

## Project Structure

```
makerlog/
├── services/
│   └── api/                 # Backend Go API
│       ├── cmd/
│       │   └── api/
│       │       └── main.go  # Entry point
│       ├── internal/
│       │   ├── database/    # Database queries
│       │   ├── handlers/    # HTTP handlers
│       │   ├── middleware/  # Middleware
│       │   └── models/      # Data models
│       ├── migrations/      # SQL migrations
│       ├── Dockerfile
│       ├── go.mod
│       └── .env.example
├── apps/
│   └── web/                 # Frontend Next.js app
│       ├── app/             # Next.js app directory
│       │   ├── page.tsx     # Home page
│       │   ├── today/       # Today page
│       │   └── projects/    # Project pages
│       ├── components/      # React components
│       ├── lib/            # Utilities and API client
│       ├── Dockerfile
│       └── .env.example
├── docker-compose.yml
├── Makefile
└── README.md
```

## License

MIT

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request
