# GoFast API

A fast and reliable Go API using Gin framework.

## Development

### Prerequisites

- Docker
- Docker Compose
- Make

### Running Tests

Tests are run using Docker to ensure consistent environments:

```bash
# Run all tests
make test

# Run tests for a specific package
docker compose run --rm dev go test -v ./routes/...
```

### Linting

We use golangci-lint through Docker for code quality checks:

```bash
# Run linter
make lint

# Fix auto-fixable issues
make lint-fix
```

The linter checks include:

* Code formatting (gofmt)
* Code simplification (gosimple)
* Error handling (errcheck)
* Import formatting (goimports)
* And more configurations in [.golangci.yml](vscode-file://vscode-app/Applications/Visual%20Studio%20Code.app/Contents/Resources/app/out/vs/code/electron-sandbox/workbench/workbench.html)

### Git Workflow & Commit Conventions

#### Branch Naming

* Always include the Jira ticket ID at the start
* Format: `GOFAST-<number>/<description>`
* Example: `GOFAST-8/add-health-check`

#### Pull Requests

* Title must start with ticket ID
* Format: `GOFAST-<number>: <description>`
* Example: `GOFAST-8: Implement health check endpoint`

#### Commit Convention

We follow Conventional Commits specification:

`<type>`(optional scope): `<description>`

Types:

* `feat`: New features
* `fix`: Bug fixes
* `docs`: Documentation changes
* `style`: Code style changes
* `refactor`: Code refactoring
* `test`: Adding/modifying tests
* `chore`: Maintenance tasks
