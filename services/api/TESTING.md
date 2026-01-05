# Testing Guide

This document describes the testing infrastructure for the Makerlog API.

## Running Tests Locally

### Run all tests
```bash
cd services/api
go test ./...
```

### Run tests with verbose output
```bash
cd services/api
go test -v ./...
```

### Run tests with coverage
```bash
cd services/api
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -html=coverage.out  # View coverage in browser
```

### Run tests with race detection
```bash
cd services/api
go test -race ./...
```

## Test Structure

- **`cmd/api/main_test.go`**: Tests for main package (environment variable handling)
- **`internal/models/models_test.go`**: Tests for model serialization/deserialization
- **`internal/middleware/auth_test.go`**: Tests for authentication middleware
- **`internal/handlers/*_test.go`**: Tests for HTTP handler validation logic
- **`internal/database/queries_test.go`**: Tests for database queries initialization

## Linting

### Run golangci-lint locally
```bash
cd services/api
golangci-lint run --timeout=5m
```

### Install golangci-lint
```bash
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## CI/CD Pipeline

The GitHub Actions workflow (`.github/workflows/go.yml`) runs on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

### Pipeline Jobs

1. **Test Job**
   - Runs all tests with race detection
   - Generates coverage report
   - Uploads coverage to Codecov (optional)
   - Warns if coverage is below 50%

2. **Lint Job**
   - Runs golangci-lint with configured linters
   - Checks code quality and style

3. **Build Job**
   - Builds the application
   - Only runs if tests and linting pass

### PR Blocking

Pull requests will be blocked if:
- Any tests fail
- Linting fails
- Build fails

## Coverage Requirements

While we don't enforce a minimum coverage threshold, we aim for:
- Critical validation logic: 100%
- Handler methods: High coverage with integration tests
- Database queries: Integration tests with real database

Current coverage focuses on:
- ✅ Request/response validation
- ✅ Authentication middleware
- ✅ JSON serialization
- ✅ Environment variable handling
- ⚠️ Database operations (requires integration tests)
- ⚠️ HTTP handler flows (requires integration tests)

## Best Practices

1. **Write tests for new code**: All new features should include tests
2. **Table-driven tests**: Use table-driven tests for multiple test cases
3. **Clear test names**: Test names should describe what they're testing
4. **Minimal mocking**: Use real implementations where possible
5. **Fast tests**: Keep unit tests fast (< 1s per test)

## Future Improvements

- Add integration tests with test database
- Add end-to-end API tests
- Increase coverage for HTTP handlers
- Add benchmark tests for critical paths
