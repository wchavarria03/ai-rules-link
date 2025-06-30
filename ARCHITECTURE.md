# Architecture Overview

## Clean Architecture Principles

This project is structured according to Clean Architecture and modern Go best practices. The main goals are:
- **Separation of concerns**: Business logic, CLI, and utilities are decoupled.
- **Testability**: All business logic is in pure Go packages, easily unit tested.
- **Extensibility**: New features can be added with minimal changes to existing code.

## Directory Layout

```
ai-rules-link/
  cmd/                  # CLI entrypoints (cobra commands)
  internal/
    domain/             # Core business logic, entities, interfaces
    service/            # Use cases, orchestration
    utils/              # Pure utility functions
  rules/                # Markdown rule files for symlinking
  main.go
  README.md
  ARCHITECTURE.md
  Dockerfile
  docker-compose.yml
  go.mod
  go.sum
```

## Key Concepts

- **Domain Layer (`internal/domain/`)**: Defines interfaces and core business logic. No dependencies on other layers.
- **Service Layer (`internal/service/`)**: Implements use cases, orchestrates domain logic, and is called by CLI/API.
- **CLI Layer (`cmd/`)**: Thin wrappers that parse arguments and call the service layer.
- **Utilities (`internal/utils/`)**: Pure, reusable helpers. No side effects or logging.
- **Rules (`rules/`)**: Markdown files that define coding, commit, and project standards for symlinking into projects.

## Extending the App

1. Add new business logic as interfaces/types in `internal/domain/`.
2. Implement new use cases in `internal/service/`.
3. Add or update CLI handlers in `cmd/`.
4. Add or update rule files in `rules/`.

## Best Practices
- All exported functions/types must have GoDoc comments.
- No business logic in CLI handlers.
- All errors must be wrapped with context.
- Use dependency injection for testability.
- Prefer interfaces for dependencies.
- Use context.Context for all public methods.

---

## ROADMAP

The following features/directories are planned for future development:
- **repository/**: Data access abstractions for persistent storage.
- **observability/**: Tracing, logging, and metrics helpers for improved monitoring.
- **api/**: REST/gRPC handlers and transport layer for exposing services.
- **configs/**: Configuration schemas and loaders for flexible deployments.
- **test/**: Test helpers, mocks, and integration/E2E tests for robust quality assurance. 