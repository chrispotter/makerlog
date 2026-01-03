# Copilot Instructions for Makerlog

## Project Overview

Makerlog is a Go-based project. This file provides guidance to GitHub Copilot when working with this codebase.

## Tech Stack

- **Language**: Go
- **Build System**: Standard Go tooling (`go build`, `go test`)

## Coding Conventions

### General Guidelines

- Follow standard Go conventions and idioms
- Use `gofmt` for code formatting
- Follow effective Go best practices
- Use meaningful variable and function names
- Keep functions focused and single-purpose

### Go-Specific Conventions

- Use camelCase for variable names
- Use PascalCase for exported functions and types
- Group imports into standard library, third-party, and local packages
- Handle errors explicitly - don't ignore them
- Use defer for cleanup operations
- Prefer table-driven tests

### Code Quality

- Write clear, self-documenting code
- Add comments for complex logic or non-obvious decisions
- Keep cyclomatic complexity low
- Avoid deep nesting

## Development Workflow

### Building

```bash
go build
```

### Testing

```bash
go test ./...
```

### Running

```bash
go run .
```

## File Organization

- Keep related functionality together
- Use packages to organize code logically
- Avoid circular dependencies

## Error Handling

- Always handle errors appropriately
- Return errors to the caller when appropriate
- Use custom error types for domain-specific errors
- Wrap errors with context using `fmt.Errorf` with `%w`

## Dependencies

- Minimize external dependencies
- Use Go modules for dependency management
- Keep dependencies up to date
- Vendor dependencies if needed for reproducibility

## Testing Guidelines

- Write tests for all public functions
- Use table-driven tests for multiple test cases
- Mock external dependencies
- Aim for high test coverage on critical paths
- Use meaningful test names that describe what's being tested

## Documentation

- Document all exported functions, types, and packages
- Use godoc-style comments
- Keep documentation up to date with code changes
- Include examples in documentation where helpful

## Security Considerations

- Never commit sensitive information (API keys, passwords, etc.)
- Use environment variables for configuration
- Validate all inputs
- Follow Go security best practices

## Performance

- Profile before optimizing
- Use benchmarks to measure performance improvements
- Be mindful of memory allocations
- Use appropriate data structures

## Git Workflow

- Write clear, descriptive commit messages
- Keep commits focused and atomic
- Review changes before committing
- Don't commit generated files or binaries
