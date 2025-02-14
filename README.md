
# GoFast - Cloud Storage API

A fast and secure cloud storage API built with Go, offering seamless file management with Azure Blob Storage integration.

## Tech Stack

- **Backend**: Go 1.22 with Gin Framework
- **Database**: MongoDB
- **Storage**: Azure Blob Storage
- **Authentication**: JWT
- **Security**: reCAPTCHA
- **Container**: Docker
- **CI/CD**: GitHub Actions
- **Deployment**: Azure Kubernetes Service (AKS)

## Getting Started

### Prerequisites

- Docker and Docker Compose
- Make

Create a `.env` file in the root directory with the following variables:

```properties
# MongoDB Configuration
MONGO_USER=           # MongoDB username
MONGO_PASSWORD=       # MongoDB password
MONGO_DATABASE=       # MongoDB database name

# Azure Blob Storage
AZURE_STORAGE_ACCOUNT_NAME=  # Azure Storage account name
AZURE_STORAGE_ACCOUNT_KEY=   # Azure Storage account key

# Authentication & Security
JWT_SECRET=          # Secret key for JWT tokens
RECAPTCHA_SECRET_KEY=  # Google reCAPTCHA v2 secret key
```

### Development

```bash
# Build and start the project with hot reload
make dev-build

# Start without rebuilding
make dev

# Stop all containers
make down

# Run tests
make test

# Run linter
make lint
```


# Fix auto-fixable issues
```
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

## API Routes

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/health` | Check API health status | No |
| GET | `/` | API root health check | No |
| POST | `/api/v1/users/register` | Register new user | No |
| POST | `/api/v1/users/login` | Login user | No |
| GET | `/api/v1/users/me` | Get current user profile | Yes |
| DELETE | `/api/v1/users/:id` | Delete user account | Yes |
| POST | `/api/v1/captcha/verify-recaptcha` | Verify reCAPTCHA token | No |
| POST | `/api/v1/assets` | Upload file | Yes |
| POST | `/api/v1/assets/folder` | Create folder | Yes |
| GET | `/api/v1/assets` | List all assets | Yes |
| GET | `/api/v1/assets/recent` | Get recent files | Yes |
| GET | `/api/v1/assets/folder/recent` | Get recent folders | Yes |
| GET | `/api/v1/assets/:id` | Get asset by ID | Yes |
| PUT | `/api/v1/assets/:id` | Update asset | Yes |
| DELETE | `/api/v1/assets/:id` | Delete asset | Yes |

### Route Groups

#### Health Check
- `/health`: System health status
- `/`: Root endpoint health check

#### Authentication
- `/api/v1/users/*`: User management endpoints
- Required header for protected routes:
  ```bash
  Authorization: Bearer <your_jwt_token>
  ```

#### Security
- `/api/v1/captcha/*`: Security verification endpoints

#### Asset Management
- `/api/v1/assets/*`: File and folder management
- File upload requires multipart/form-data with:
  - `file`: File to upload
  - `containerId`: Target container ID
  - `blobName`: Desired file name