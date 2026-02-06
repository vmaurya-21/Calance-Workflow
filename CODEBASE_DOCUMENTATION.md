# Calance Workflow Backend - Comprehensive Codebase Documentation

## Table of Contents
1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Project Structure](#project-structure)
4. [Core Packages](#core-packages)
5. [File-by-File Documentation](#file-by-file-documentation)
6. [Function Reference](#function-reference)

---

## Project Overview

**Calance Workflow Backend** is a Go-based backend server that provides GitHub workflow management capabilities with OAuth authentication, JWT-based sessions, and PostgreSQL database integration.

### Key Features
- GitHub OAuth 2.0 authentication
- JWT-based session management
- PostgreSQL database with GORM ORM
- RESTful API architecture
- Clean three-layer architecture (API handlers, domain services, infrastructure)
- CORS support for frontend integration
- GitHub Actions workflow creation and management
- EC2 and Kubernetes deployment workflow generation

### Technology Stack
- **Language**: Go 1.25.4
- **Web Framework**: Gin (v1.11.0)
- **Database**: PostgreSQL with GORM (v1.31.1)
- **Authentication**: OAuth 2.0 (GitHub) + JWT (golang-jwt/jwt v5.3.0)
- **Logging**: Zerolog (v1.34.0)
- **Configuration**: godotenv (v1.5.1)

---

## Architecture

The application follows a **clean three-layer architecture**:

```
┌─────────────────────────────────────────────┐
│         API Layer (Handlers)                │
│  - HTTP request/response handling           │
│  - Input validation                         │
│  - Response formatting                      │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│         Domain Layer (Services)             │
│  - Business logic                           │
│  - Domain models                            │
│  - Service orchestration                    │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│      Infrastructure Layer                   │
│  - GitHub API client                        │
│  - Database repositories                    │
│  - Template generators                      │
│  - External service integrations            │
└─────────────────────────────────────────────┘
```

### Data Flow
1. **Request** → API Handler (validates input)
2. **Handler** → Domain Service (executes business logic)
3. **Service** → Infrastructure (interacts with external systems)
4. **Infrastructure** → External Systems (GitHub API, Database)
5. **Response** ← Returns through the layers

---

## Project Structure

```
Calance-Workflow-backend/
├── cmd/
│   └── server/
│       └── main.go                     # Application entry point
├── internal/
│   ├── api/
│   │   └── handlers/                   # HTTP request handlers
│   │       ├── auth/                   # Authentication handlers
│   │       ├── organization/           # Organization handlers
│   │       ├── repository/             # Repository handlers
│   │       └── workflow/               # Workflow handlers
│   ├── config/
│   │   └── config.go                   # Configuration management
│   ├── database/
│   │   └── database.go                 # Database initialization
│   ├── domain/                         # Business logic layer
│   │   ├── auth/                       # Auth domain
│   │   ├── organization/               # Organization domain
│   │   ├── repository/                 # Repository domain
│   │   └── workflow/                   # Workflow domain
│   ├── infrastructure/                 # External integrations
│   │   ├── database/
│   │   │   └── repositories/           # Data access layer
│   │   ├── github/                     # GitHub API client
│   │   └── template/                   # Workflow template generators
│   ├── logger/
│   │   ├── logger.go                   # Logging configuration
│   │   └── middleware.go               # HTTP logging middleware
│   ├── middleware/
│   │   └── auth_middleware.go          # JWT & CORS middleware
│   ├── pkg/
│   │   ├── http/                       # HTTP utilities
│   │   └── logger/                     # Logger utilities
│   ├── router/
│   │   └── router.go                   # Route configuration
│   └── utils/
│       ├── jwt.go                      # JWT utilities
│       └── response.go                 # API response helpers
├── db/
│   └── migrations/                     # Database migrations
├── .env                                # Environment variables
├── .env.example                        # Environment template
├── docker-compose.yml                  # Docker configuration
├── go.mod                              # Go module dependencies
└── README.md                           # Project documentation
```

---

## Core Packages

### 1. **cmd/server** - Application Entry Point
Contains the main application entry point that initializes all components.

### 2. **internal/config** - Configuration Management
Handles loading and validation of environment variables and application configuration.

### 3. **internal/database** - Database Layer
Manages PostgreSQL connection, migrations, and database lifecycle.

### 4. **internal/router** - Routing
Configures all HTTP routes and connects handlers with middleware.

### 5. **internal/logger** - Logging
Provides structured logging using Zerolog with different output formats.

### 6. **internal/middleware** - HTTP Middleware
Contains authentication (JWT validation) and CORS middleware.

### 7. **internal/domain** - Business Logic
Core business logic organized by domain:
- **auth**: User authentication and authorization
- **workflow**: GitHub workflow generation and management
- **organization**: GitHub organization operations
- **repository**: Repository and package management

### 8. **internal/infrastructure** - External Integrations
Handles communication with external systems:
- **github**: GitHub API client
- **database/repositories**: Data access layer
- **template**: Workflow template generators

### 9. **internal/api/handlers** - HTTP Handlers
HTTP request handlers organized by domain.

### 10. **internal/utils** - Utilities
Common utilities for JWT, HTTP responses, etc.

---

## File-by-File Documentation

### **cmd/server/main.go**

**Purpose**: Application entry point that bootstraps the server.

**Key Responsibilities**:
- Load configuration from environment variables
- Initialize logger with configured level and format
- Initialize database connection and run migrations
- Set up HTTP router with all routes and middleware
- Start HTTP server with graceful shutdown support

**Main Components**:
- Configuration loading
- Logger initialization
- Database initialization
- Router setup
- Graceful shutdown handling

---

### **internal/config/config.go**

**Purpose**: Centralized configuration management for the application.

**Structures**:

#### `Config`
Main configuration struct containing all app settings.

**Fields**:
- `Server`: Server configuration (port, mode, environment)
- `Database`: Database connection settings
- `GitHub`: GitHub OAuth credentials
- `JWT`: JWT secret and expiration
- `Frontend`: Frontend URL and allowed origins
- `Log`: Logging level and format

#### `ServerConfig`
- `Port`: HTTP server port (default: 8080)
- `GinMode`: Gin framework mode (debug/release)
- `Env`: Environment (development/production)

#### `DatabaseConfig`
- `Host`: Database host
- `Port`: Database port (default: 5432)
- `User`: Database user
- `Password`: Database password
- `DBName`: Database name
- `SSLMode`: SSL mode (disable/require)

#### `GitHubConfig`
- `ClientID`: GitHub OAuth app client ID
- `ClientSecret`: GitHub OAuth app client secret
- `RedirectURL`: OAuth callback URL

#### `JWTConfig`
- `Secret`: JWT signing secret
- `ExpirationHours`: Token expiration time in hours

#### `FrontendConfig`
- `URL`: Frontend application URL
- `AllowedOrigins`: CORS allowed origins

#### `LogConfig`
- `Level`: Log level (debug/info/warn/error/fatal)
- `Format`: Log format (json/console)

**Functions**:

#### `LoadConfig() (*Config, error)`
Loads configuration from environment variables using godotenv.

**Returns**:
- `*Config`: Populated configuration struct
- `error`: Error if required fields are missing

**Process**:
1. Loads `.env` file (if exists)
2. Reads environment variables with defaults
3. Validates required fields
4. Returns configuration or error

#### `Validate() error`
Validates that all required configuration fields are set.

**Validates**:
- GitHub Client ID (required)
- GitHub Client Secret (required)
- JWT Secret (required)
- Database Password (warning if empty)

#### `GetDatabaseDSN() string`
Constructs PostgreSQL connection string from database config.

**Returns**: DSN string in format: `host=X port=X user=X password=X dbname=X sslmode=X`

#### `getEnv(key, defaultValue string) string`
Helper function to get environment variable with fallback default.

#### `getEnvAsInt(key string, defaultValue int) int`
Helper function to get environment variable as integer.

#### `getEnvAsSlice(key string, defaultValue []string) []string`
Helper function to get comma-separated environment variable as string slice.

---

### **internal/database/database.go**

**Purpose**: Manages PostgreSQL database connection and migrations.

**Global Variables**:
- `DB *gorm.DB`: Global database instance

**Functions**:

#### `InitDatabase(cfg *config.Config) error`
Initializes database connection and runs migrations.

**Parameters**:
- `cfg`: Application configuration

**Returns**:
- `error`: Error if connection or migration fails

**Process**:
1. Constructs DSN from config
2. Configures GORM logger
3. Opens PostgreSQL connection
4. Runs auto-migrations
5. Returns error or nil

#### `runMigrations() error`
Runs database migrations using GORM AutoMigrate.

**Migrates**:
- `auth.User` table
- `auth.Token` table

**Note**: For production, use golang-migrate CLI with SQL migration files.

#### `GetDB() *gorm.DB`
Returns the global database instance.

**Returns**: `*gorm.DB` database instance

#### `CloseDatabase() error`
Closes the database connection gracefully.

**Returns**: `error` if closing fails

---

### **internal/router/router.go**

**Purpose**: Configures all HTTP routes and middleware.

**Functions**:

#### `SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine`
Sets up and configures the Gin router with all routes and middleware.

**Parameters**:
- `db`: Database instance
- `cfg`: Application configuration

**Returns**: `*gin.Engine` configured router

**Configuration**:
1. Sets Gin mode (debug/release)
2. Creates router with default middleware
3. Applies logging middleware
4. Applies CORS middleware
5. Initializes repositories
6. Initializes domain services
7. Initializes handlers
8. Configures routes

**Routes**:

**Public Routes**:
- `GET /ping` - Health check
- `GET /api/auth/github` - Initiate GitHub OAuth
- `GET /api/auth/github/callback` - OAuth callback

**Protected Routes** (require JWT):
- `GET /api/auth/me` - Get current user profile
- `POST /api/auth/logout` - Logout user
- `GET /api/auth/organizations` - List user organizations
- `GET /api/auth/organizations/:org/repositories` - Get org repositories
- `GET /api/auth/repositories` - List user repositories
- `GET /api/repositories/:owner/:repo/branches` - Get repository branches
- `GET /api/repositories/:owner/:repo/branches/:branch/commits` - Get branch commits
- `GET /api/repositories/:owner/:repo/tags` - Get repository tags
- `POST /api/repositories/tags` - Create new tag
- `GET /api/repositories/:owner/:repo/actions/runs` - Get workflow runs
- `GET /api/repositories/:owner/:repo/actions/runs/:run_id` - Get workflow run details
- `GET /api/repositories/:owner/:repo/actions/jobs/:job_id/logs` - Get job logs
- `GET /api/packages/user` - Get user packages
- `GET /api/packages/org/:org` - Get organization packages
- `GET /api/workflows/:owner/:repo` - List workflows
- `POST /api/workflows/create` - Create workflow
- `POST /api/workflows/preview` - Preview workflow
- `GET /api/workflows/:owner/:repo/file` - Get workflow content
- `PUT /api/workflows/:owner/:repo/file` - Update workflow

---

### **internal/logger/logger.go**

**Purpose**: Provides structured logging using Zerolog.

**Global Variables**:
- `Logger zerolog.Logger`: Global logger instance

**Functions**:

#### `InitLogger(logLevel, logFormat string)`
Initializes the global logger with configuration.

**Parameters**:
- `logLevel`: Log level (debug/info/warn/error/fatal)
- `logFormat`: Output format (json/console)

**Behavior**:
- **Console format**: Pretty-printed output for development
- **JSON format**: Structured JSON logs for production
- Adds timestamp and caller information to all logs

#### `parseLogLevel(level string) zerolog.Level`
Converts string log level to zerolog.Level.

**Parameters**:
- `level`: String log level

**Returns**: `zerolog.Level` corresponding to the string

#### `GetLogger() *zerolog.Logger`
Returns the global logger instance.

#### `Debug() *zerolog.Event`
Creates a debug-level log event.

#### `Info() *zerolog.Event`
Creates an info-level log event.

#### `Warn() *zerolog.Event`
Creates a warning-level log event.

#### `Error() *zerolog.Event`
Creates an error-level log event.

#### `Fatal() *zerolog.Event`
Creates a fatal-level log event (exits after logging).

---

### **internal/logger/middleware.go**

**Purpose**: Gin middleware for HTTP request logging.

**Functions**:

#### `GinMiddleware() gin.HandlerFunc`
Returns a Gin middleware function for logging HTTP requests.

**Logs**:
- Request ID (UUID)
- HTTP method
- Request path
- Query parameters
- Status code
- Response latency
- Client IP
- User agent
- Errors (if any)

**Log Levels**:
- **Info**: Successful requests (2xx, 3xx)
- **Warn**: Client errors (4xx)
- **Error**: Server errors (5xx) or requests with errors

---

### **internal/middleware/auth_middleware.go**

**Purpose**: JWT authentication and CORS middleware.

**Functions**:

#### `AuthMiddleware() gin.HandlerFunc`
Validates JWT tokens and sets user context.

**Process**:
1. Extracts Authorization header
2. Validates Bearer token format
3. Validates JWT signature and expiration
4. Sets `user_id` and `username` in context
5. Aborts with 401 if validation fails

**Context Values Set**:
- `user_id`: User UUID as string
- `username`: Username string

#### `OptionalAuthMiddleware() gin.HandlerFunc`
Validates JWT tokens but doesn't require them.

**Behavior**:
- If token present and valid: sets user context
- If token absent or invalid: continues without user context

#### `GetUserIDFromContext(c *gin.Context) (uuid.UUID, bool)`
Extracts user ID from Gin context.

**Parameters**:
- `c`: Gin context

**Returns**:
- `uuid.UUID`: User ID
- `bool`: True if user ID exists and is valid

#### `CORSMiddleware(allowedOrigins []string) gin.HandlerFunc`
Handles CORS headers for cross-origin requests.

**Parameters**:
- `allowedOrigins`: List of allowed origin URLs

**Headers Set**:
- `Access-Control-Allow-Origin`: Matched origin
- `Access-Control-Allow-Credentials`: true
- `Access-Control-Allow-Headers`: Content-Type, Authorization, etc.
- `Access-Control-Allow-Methods`: POST, GET, PUT, DELETE, PATCH, OPTIONS

**Behavior**:
- Checks if request origin is in allowed list
- Handles OPTIONS preflight requests
- Returns 204 for OPTIONS requests

---

### **internal/utils/jwt.go**

**Purpose**: JWT token generation and validation utilities.

**Structures**:

#### `JWTClaims`
Custom JWT claims structure.

**Fields**:
- `UserID uuid.UUID`: User's unique identifier
- `Username string`: User's username
- `jwt.RegisteredClaims`: Standard JWT claims (exp, iat, nbf, iss)

**Functions**:

#### `GenerateToken(userID uuid.UUID, username string) (string, error)`
Generates a new JWT token for a user.

**Parameters**:
- `userID`: User's UUID
- `username`: User's username

**Returns**:
- `string`: Signed JWT token
- `error`: Error if generation fails

**Token Properties**:
- **Algorithm**: HS256 (HMAC with SHA-256)
- **Issuer**: "calance-workflow"
- **Expiration**: Configured hours from config
- **Issued At**: Current time
- **Not Before**: Current time

#### `ValidateToken(tokenString string) (*JWTClaims, error)`
Validates a JWT token and returns claims.

**Parameters**:
- `tokenString`: JWT token string

**Returns**:
- `*JWTClaims`: Parsed claims if valid
- `error`: `ErrExpiredToken` or `ErrInvalidToken`

**Validation**:
- Verifies signature using JWT secret
- Checks signing method is HMAC
- Validates expiration time
- Returns parsed claims

#### `ExtractTokenFromHeader(authHeader string) (string, error)`
Extracts JWT token from Authorization header.

**Parameters**:
- `authHeader`: Authorization header value

**Returns**:
- `string`: Token without "Bearer " prefix
- `error`: Error if format is invalid

**Expected Format**: `Bearer <token>`

---

### **internal/utils/response.go**

**Purpose**: Standardized API response helpers.

**Structures**:

#### `Response`
Standard API response structure.

**Fields**:
- `Success bool`: Indicates if request was successful
- `Message string`: Human-readable message
- `Data interface{}`: Response data (optional)
- `Error string`: Error message (optional)

**Functions**:

#### `SuccessResponse(c *gin.Context, statusCode int, message string, data interface{})`
Sends a success response.

**Parameters**:
- `c`: Gin context
- `statusCode`: HTTP status code
- `message`: Success message
- `data`: Response data

#### `ErrorResponse(c *gin.Context, statusCode int, message string, err error)`
Sends an error response.

**Parameters**:
- `c`: Gin context
- `statusCode`: HTTP status code
- `message`: Error message
- `err`: Error object (optional)

#### `ValidationErrorResponse(c *gin.Context, errors interface{})`
Sends a validation error response (400).

**Parameters**:
- `c`: Gin context
- `errors`: Validation error details

#### `UnauthorizedResponse(c *gin.Context, message string)`
Sends an unauthorized response (401).

#### `NotFoundResponse(c *gin.Context, message string)`
Sends a not found response (404).

#### `InternalServerErrorResponse(c *gin.Context, message string, err error)`
Sends an internal server error response (500).

#### `BadRequestResponse(c *gin.Context, message string)`
Sends a bad request response (400).

---

### **internal/domain/auth/models.go**

**Purpose**: Authentication domain models.

**Structures**:

#### `User`
Represents a user in the system.

**Fields**:
- `ID uuid.UUID`: Primary key (UUID)
- `GitHubID int64`: GitHub user ID (unique)
- `Username string`: GitHub username
- `Email string`: User email
- `AvatarURL string`: Profile avatar URL
- `Name string`: Full name
- `Bio string`: User bio
- `Location string`: User location
- `Company string`: User company
- `CreatedAt time.Time`: Creation timestamp
- `UpdatedAt time.Time`: Last update timestamp
- `DeletedAt gorm.DeletedAt`: Soft delete timestamp

**Methods**:

##### `BeforeCreate(tx *gorm.DB) error`
GORM hook that generates UUID before creating record.

##### `ToResponse() map[string]interface{}`
Converts User to response format (excludes sensitive data).

#### `Token`
Represents a GitHub access token.

**Fields**:
- `ID uuid.UUID`: Primary key (UUID)
- `UserID uuid.UUID`: Foreign key to User (unique)
- `AccessToken string`: GitHub access token
- `TokenType string`: Token type (usually "Bearer")
- `Scope string`: OAuth scopes
- `ExpiresAt *time.Time`: Token expiration
- `CreatedAt time.Time`: Creation timestamp
- `UpdatedAt time.Time`: Last update timestamp
- `DeletedAt gorm.DeletedAt`: Soft delete timestamp

**Methods**:

##### `BeforeCreate(tx *gorm.DB) error`
GORM hook that generates UUID before creating record.

---

### **internal/domain/auth/service.go**

**Purpose**: Authentication business logic service.

**Structure**:

#### `Service`
Authentication service.

**Fields**:
- `githubAuth *github.AuthClient`: GitHub OAuth client

**Functions**:

#### `NewService(clientID, clientSecret, redirectURL string, scopes []string) *Service`
Creates a new auth service.

**Parameters**:
- `clientID`: GitHub OAuth client ID
- `clientSecret`: GitHub OAuth client secret
- `redirectURL`: OAuth callback URL
- `scopes`: OAuth permission scopes

**Returns**: `*Service` instance

#### `GetAuthURL(state string) string`
Returns the GitHub OAuth authorization URL.

**Parameters**:
- `state`: CSRF protection state token

**Returns**: Authorization URL string

#### `ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)`
Exchanges authorization code for access token.

**Parameters**:
- `ctx`: Context
- `code`: Authorization code from GitHub

**Returns**:
- `*oauth2.Token`: Access token
- `error`: Error if exchange fails

#### `GetGitHubUser(ctx context.Context, token string) (*User, error)`
Fetches user information from GitHub.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token

**Returns**:
- `*User`: User domain model
- `error`: Error if fetch fails

#### `GenerateJWT(userID uuid.UUID, username string) (string, error)`
Generates a JWT token for a user.

**Parameters**:
- `userID`: User UUID
- `username`: Username

**Returns**:
- `string`: JWT token
- `error`: Error if generation fails

#### `ValidateJWT(tokenString string) (*utils.JWTClaims, error)`
Validates a JWT token.

**Parameters**:
- `tokenString`: JWT token

**Returns**:
- `*utils.JWTClaims`: Token claims
- `error`: Error if validation fails

---

### **internal/domain/auth/errors.go**

**Purpose**: Authentication domain errors.

**Error Variables**:
- `ErrInvalidCredentials`: Invalid login credentials
- `ErrUserNotFound`: User not found in database
- `ErrTokenNotFound`: Token not found in database
- `ErrTokenExpired`: Token has expired
- `ErrInvalidToken`: Token is invalid

---

### **internal/infrastructure/database/repositories/user.go**

**Purpose**: User data access layer.

**Structure**:

#### `UserRepository`
Handles user database operations.

**Fields**:
- `db *gorm.DB`: Database instance

**Functions**:

#### `NewUserRepository(db *gorm.DB) *UserRepository`
Creates a new user repository.

#### `FindByID(id uuid.UUID) (*auth.User, error)`
Finds a user by UUID.

**Parameters**:
- `id`: User UUID

**Returns**:
- `*auth.User`: User if found
- `error`: Error if not found or query fails

#### `FindByGitHubID(githubID int64) (*auth.User, error)`
Finds a user by GitHub ID.

**Parameters**:
- `githubID`: GitHub user ID

**Returns**:
- `*auth.User`: User if found
- `error`: Error if not found or query fails

#### `CreateOrUpdate(user *auth.User) error`
Creates a new user or updates existing user.

**Parameters**:
- `user`: User to create/update

**Returns**: `error` if operation fails

**Behavior**:
- Checks if user exists by GitHub ID
- Creates new record if not found
- Updates existing record if found
- Preserves existing UUID on update

---

### **internal/infrastructure/database/repositories/token.go**

**Purpose**: Token data access layer.

**Structure**:

#### `TokenRepository`
Handles token database operations.

**Fields**:
- `db *gorm.DB`: Database instance

**Functions**:

#### `NewTokenRepository(db *gorm.DB) *TokenRepository`
Creates a new token repository.

#### `FindByUserID(userID uuid.UUID) (*auth.Token, error)`
Finds a token by user ID.

**Parameters**:
- `userID`: User UUID

**Returns**:
- `*auth.Token`: Token if found, nil if not found
- `error`: Error if query fails

#### `CreateOrUpdate(token *auth.Token) error`
Creates a new token or updates existing token.

**Parameters**:
- `token`: Token to create/update

**Returns**: `error` if operation fails

**Behavior**:
- Checks if token exists for user
- Creates new record if not found
- Updates existing record if found
- Preserves existing UUID on update

---

### **internal/infrastructure/github/auth.go**

**Purpose**: GitHub OAuth operations.

**Structure**:

#### `AuthClient`
Handles GitHub OAuth authentication.

**Fields**:
- `*Client`: Embedded GitHub API client
- `oauthConfig *oauth2.Config`: OAuth configuration

**Functions**:

#### `NewAuthClient(clientID, clientSecret, redirectURL string, scopes []string) *AuthClient`
Creates a new GitHub OAuth client.

**Parameters**:
- `clientID`: GitHub OAuth app client ID
- `clientSecret`: GitHub OAuth app client secret
- `redirectURL`: OAuth callback URL
- `scopes`: OAuth permission scopes

**Returns**: `*AuthClient` instance

#### `GetAuthURL(state string) string`
Returns the GitHub OAuth authorization URL.

**Parameters**:
- `state`: CSRF protection state

**Returns**: Authorization URL

#### `ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error)`
Exchanges authorization code for access token.

**Parameters**:
- `ctx`: Context
- `code`: Authorization code

**Returns**:
- `*oauth2.Token`: Access token
- `error`: Error if exchange fails

#### `GetUser(ctx context.Context, token string) (*User, error)`
Fetches GitHub user information.

**Parameters**:
- `ctx`: Context (with 30s timeout)
- `token`: GitHub access token

**Returns**:
- `*User`: GitHub user data
- `error`: Error if fetch fails

**Process**:
1. Makes GET request to `/user` endpoint
2. If email is empty, fetches from `/user/emails`
3. Returns user with primary verified email

#### `getPrimaryEmail(ctx context.Context, token string) (string, error)`
Fetches user's primary verified email.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token

**Returns**:
- `string`: Primary verified email
- `error`: Error if no verified email found

**Priority**:
1. Primary verified email
2. First verified email
3. Error if none found

---

### **internal/infrastructure/github/client.go**

**Purpose**: Base GitHub API client.

**Structure**:

#### `Client`
Handles GitHub API interactions.

**Fields**:
- `httpClient *pkghttp.Client`: HTTP client
- `baseURL string`: GitHub API base URL

**Functions**:

#### `NewClient() *Client`
Creates a new GitHub API client.

**Returns**: `*Client` with base URL `https://api.github.com`

#### `doRequest(ctx context.Context, token, method, path string, body interface{}) (*pkghttp.Response, error)`
Performs a GitHub API request.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `method`: HTTP method
- `path`: API endpoint path
- `body`: Request body (optional)

**Returns**:
- `*pkghttp.Response`: API response
- `error`: Error if request fails

**Headers Set**:
- `Authorization`: Bearer token
- `Accept`: application/vnd.github+json
- `X-GitHub-Api-Version`: 2022-11-28

#### `checkResponse(resp *pkghttp.Response) error`
Checks if API response is successful.

**Parameters**:
- `resp`: HTTP response

**Returns**: `error` based on status code

**Error Mapping**:
- 401 → `ErrUnauthorized`
- 403 → `ErrForbidden`
- 404 → `ErrNotFound`
- Other → `ErrAPIFailed` with message

---

### **internal/infrastructure/github/types.go**

**Purpose**: GitHub API response types.

**Structures**:

#### `Repository`
GitHub repository data.

**Fields**:
- `ID int64`: Repository ID
- `Name string`: Repository name
- `FullName string`: Full name (owner/repo)
- `Description string`: Repository description
- `Private bool`: Is private repository
- `HTMLURL string`: Repository web URL
- `DefaultBranch string`: Default branch name
- `Owner Owner`: Repository owner

#### `Owner`
Repository or organization owner.

**Fields**:
- `Login string`: Username or org name
- `ID int64`: Owner ID
- `AvatarURL string`: Avatar image URL
- `Type string`: "User" or "Organization"

#### `Branch`
Repository branch.

**Fields**:
- `Name string`: Branch name
- `Commit Commit`: Latest commit
- `Protected bool`: Is protected branch

#### `Commit`
Git commit.

**Fields**:
- `SHA string`: Commit SHA
- `URL string`: Commit API URL

#### `Ref`
Git reference (tag/branch).

**Fields**:
- `Ref string`: Reference name
- `NodeID string`: Node ID
- `URL string`: Reference URL
- `Object`: Object details (SHA, type, URL)

#### `PullRequest`
GitHub pull request.

**Fields**:
- `ID int64`: PR ID
- `Number int`: PR number
- `State string`: PR state (open/closed)
- `Title string`: PR title
- `Body string`: PR description
- `HTMLURL string`: PR web URL
- `Head`: Head branch (ref, SHA)
- `Base`: Base branch (ref, SHA)

#### `Content`
File or directory content.

**Fields**:
- `Name string`: File/directory name
- `Path string`: Full path
- `SHA string`: Content SHA
- `Size int`: Size in bytes
- `URL string`: API URL
- `HTMLURL string`: Web URL
- `GitURL string`: Git URL
- `DownloadURL string`: Download URL
- `Type string`: "file" or "dir"
- `Content string`: Base64 encoded content
- `Encoding string`: Content encoding

#### `User`
GitHub user.

**Fields**:
- `ID int64`: User ID
- `Login string`: Username
- `Email string`: Email address
- `AvatarURL string`: Avatar URL
- `Name string`: Full name
- `Bio string`: User bio
- `Location string`: Location
- `Company string`: Company

#### `Organization`
GitHub organization.

**Fields**:
- `ID int64`: Organization ID
- `Login string`: Organization name
- `AvatarURL string`: Avatar URL
- `Description string`: Organization description

#### `Package`
GitHub package.

**Fields**:
- `ID int64`: Package ID
- `Name string`: Package name
- `PackageType string`: Package type (npm, docker, etc.)
- `Visibility string`: public/private
- `URL string`: API URL
- `HTMLURL string`: Web URL
- `CreatedAt string`: Creation timestamp
- `UpdatedAt string`: Update timestamp
- `Owner Owner`: Package owner
- `Repository *Repository`: Associated repository

#### `PackageVersion`
Package version.

**Fields**:
- `ID int64`: Version ID
- `Name string`: Version name
- `URL string`: Version URL
- `CreatedAt string`: Creation timestamp
- `UpdatedAt string`: Update timestamp

---

### **internal/infrastructure/github/errors.go**

**Purpose**: GitHub API error definitions.

**Error Variables**:
- `ErrUnauthorized`: 401 Unauthorized
- `ErrForbidden`: 403 Forbidden
- `ErrNotFound`: 404 Not Found
- `ErrAPIFailed`: General API failure

---

### **internal/api/handlers/auth/handler.go**

**Purpose**: Auth handler initialization.

**Structure**:

#### `Handler`
Auth HTTP request handler.

**Fields**:
- `authService *auth.Service`: Auth domain service
- `userRepository *database.UserRepository`: User repository
- `tokenRepository *database.TokenRepository`: Token repository

**Functions**:

#### `NewHandler(authService, userRepo, tokenRepo) *Handler`
Creates a new auth handler.

**Parameters**:
- `authService`: Auth service
- `userRepo`: User repository
- `tokenRepo`: Token repository

**Returns**: `*Handler` instance

#### `getUserID(c interface{}) (uuid.UUID, error)`
Helper to extract user ID from context (placeholder).

---

### **internal/api/handlers/auth/login.go**

**Purpose**: GitHub OAuth login handler.

**Functions**:

#### `Login(c *gin.Context)`
Redirects user to GitHub OAuth authorization page.

**Route**: `GET /api/auth/github`

**Process**:
1. Generates random CSRF state token
2. Stores state in secure HTTP-only cookie (5 min expiration)
3. Gets GitHub OAuth authorization URL
4. Redirects user to GitHub

**Cookie**:
- Name: `oauth_state`
- Max-Age: 300 seconds
- Path: `/`
- HttpOnly: true
- SameSite: Lax

#### `generateRandomState() (string, error)`
Generates random state for CSRF protection.

**Returns**:
- `string`: Random state string
- `error`: Error if generation fails

---

### **internal/api/handlers/auth/callback.go**

**Purpose**: GitHub OAuth callback handler.

**Functions**:

#### `Callback(c *gin.Context)`
Handles OAuth callback from GitHub.

**Route**: `GET /api/auth/github/callback`

**Query Parameters**:
- `state`: CSRF state token
- `code`: Authorization code

**Process**:
1. Verifies state matches cookie (CSRF protection)
2. Clears state cookie
3. Exchanges authorization code for access token
4. Fetches GitHub user information
5. Creates or updates user in database
6. Creates or updates access token in database
7. Generates JWT token
8. Redirects to frontend with JWT token

**Redirect URL**: `{FRONTEND_URL}/auth/callback?token={JWT}`

**Error Responses**:
- 400: Invalid state or missing code
- 500: Token exchange, user fetch, or database errors

---

### **internal/api/handlers/auth/profile.go**

**Purpose**: User profile handlers.

**Functions**:

#### `GetProfile(c *gin.Context)`
Returns the current authenticated user.

**Route**: `GET /api/auth/me`

**Authentication**: Required (JWT)

**Process**:
1. Extracts user ID from context (set by auth middleware)
2. Fetches user from database
3. Returns user data in response format

**Response**:
```json
{
  "success": true,
  "message": "User fetched successfully",
  "data": {
    "id": "uuid",
    "github_id": 12345,
    "username": "user",
    "email": "user@example.com",
    "avatar_url": "https://...",
    "name": "Full Name",
    "bio": "Bio text",
    "location": "Location",
    "company": "Company",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Error Responses**:
- 401: User not found in context
- 404: User not found in database
- 500: Database error

#### `Logout(c *gin.Context)`
Logs out the user.

**Route**: `POST /api/auth/logout`

**Authentication**: Required (JWT)

**Note**: JWT is stateless, so logout is handled client-side by removing the token.

**Response**:
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

### **internal/domain/workflow/models.go**

**Purpose**: Workflow domain models and request/response structures.

**Constants**:

#### `DeploymentType`
Deployment type enum.

**Values**:
- `DeploymentTypeEC2`: "ec2"
- `DeploymentTypeKubernetes`: "kubernetes"

**Structures**:

#### `Request`
Workflow creation request.

**Fields**:
- `Owner string`: Repository owner
- `Repository string`: Repository name
- `WorkflowName string`: Workflow name
- `DeploymentType DeploymentType`: ec2 or kubernetes
- `Projects []Project`: Common project configurations
- `EC2CommonFields *EC2CommonFields`: EC2 shared config
- `EC2Projects []EC2Project`: EC2 project configs
- `KubernetesCommonFields *KubernetesCommonFields`: K8s shared config
- `KubernetesProjects []KubernetesProject`: K8s project configs

**Methods**:

##### `Validate() error`
Validates the request based on deployment type.

**Validation**:
- EC2: Requires `EC2CommonFields` and `EC2Projects`
- Kubernetes: Requires `KubernetesCommonFields` and `KubernetesProjects`

#### `Project`
Common project configuration.

**Fields**:
- `ID string`: Project ID
- `Name string`: Project name
- `DockerContextPath string`: Docker build context
- `DockerfilePath string`: Dockerfile path
- `DotEnvTesting string`: Testing environment variables
- `DotEnvProduction string`: Production environment variables

#### `EC2CommonFields`
Shared EC2 configuration.

**Fields**:
- `CredentialID string`: AWS credential ID
- `AWSRegion string`: AWS region
- `JenkinsJobs string`: Jenkins job names
- `ReleaseTag string`: Release tag pattern
- `CodeownersEmails string`: Code owner emails
- `DevopsStakeholdersEmails string`: DevOps stakeholder emails

#### `EC2Project`
EC2-specific project configuration.

**Fields**:
- `ID string`: Project ID
- `Name string`: Project name
- `Command string`: Container command
- `Port string`: Container port
- `DockerNetwork string`: Docker network
- `MountPath string`: Volume mount path
- `EnableGPU bool`: Enable GPU support
- `LogDriver string`: Docker log driver
- `LogDriverOptions string`: Log driver options

#### `KubernetesCommonFields`
Shared Kubernetes configuration.

**Fields**:
- `JenkinsJobName string`: Jenkins job name
- `ReleaseTag string`: Release tag pattern
- `HelmValuesRepository string`: Helm values repo
- `CodeownersEmailIds string`: Code owner emails
- `DevopsStakeholdersEmailIds string`: DevOps stakeholder emails

#### `KubernetesProject`
Kubernetes-specific project configuration.

**Fields**:
- `ID string`: Project ID
- `Name string`: Project name

#### `Response`
Workflow creation response.

**Fields**:
- `Owner string`: Repository owner
- `Repository string`: Repository name
- `WorkflowName string`: Workflow name
- `FilePath string`: Workflow file path
- `FileURL string`: Pull request URL
- `ContentSHA string`: File content SHA
- `Message string`: Success message
- `CreatedAt time.Time`: Creation timestamp

#### `File`
Workflow file metadata.

**Fields**:
- `Name string`: File name
- `Path string`: File path
- `SHA string`: Content SHA
- `Size int`: File size
- `URL string`: Web URL
- `DownloadURL string`: Download URL

#### `History`
Workflow creation history record.

**Fields**:
- `ID uuid.UUID`: Record ID
- `UserID uuid.UUID`: User ID
- `Owner string`: Repository owner
- `Repository string`: Repository name
- `WorkflowName string`: Workflow name
- `DeploymentType DeploymentType`: Deployment type
- `FilePath string`: File path
- `ContentSHA string`: Content SHA
- `Status string`: Creation status
- `ErrorMessage *string`: Error message (if failed)
- `CreatedAt time.Time`: Creation timestamp
- `UpdatedAt time.Time`: Update timestamp
- `DeletedAt gorm.DeletedAt`: Soft delete timestamp

**Methods**:

##### `BeforeCreate(tx *gorm.DB) error`
GORM hook to generate UUID.

#### `FileContentResponse`
Workflow file content response.

**Fields**:
- `Name string`: File name
- `Path string`: File path
- `SHA string`: Content SHA
- `Size int`: Content size
- `Content string`: File content (decoded)

#### `UpdateWorkflowRequest`
Workflow update request.

**Fields**:
- `Owner string`: Repository owner
- `Repository string`: Repository name
- `FilePath string`: Workflow file path
- `Content string`: New file content
- `SHA string`: Current file SHA
- `CommitMessage string`: Commit message (optional)

---

### **internal/domain/workflow/service.go**

**Purpose**: Workflow business logic service.

**Structure**:

#### `Service`
Workflow service.

**Fields**:
- `githubClient *github.WorkflowClient`: GitHub workflow client
- `ec2Template *template.EC2Generator`: EC2 template generator
- `k8sTemplate *template.KubernetesGenerator`: K8s template generator

**Functions**:

#### `NewService() *Service`
Creates a new workflow service.

**Returns**: `*Service` instance

#### `GenerateWorkflow(req *Request) (string, error)`
Generates workflow YAML based on request.

**Parameters**:
- `req`: Workflow request

**Returns**:
- `string`: Generated YAML content
- `error`: Error if generation fails

**Process**:
1. Validates workflow name (alphanumeric, dash, underscore, max 255 chars)
2. Validates request based on deployment type
3. Generates YAML using appropriate template generator
4. Returns YAML or error

#### `CreateWorkflow(ctx, token, owner, repo, workflowName, content) (*Response, error)`
Creates a workflow in GitHub repository.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `workflowName`: Workflow name
- `content`: Workflow YAML content

**Returns**:
- `*Response`: Creation response with PR details
- `error`: Error if creation fails

**Process**:
1. Verifies repository exists
2. Gets default branch
3. Creates new branch: `workflow/{name}-{timestamp}`
4. Gets base branch SHA
5. Creates workflow file: `.github/workflows/{name}.yml`
6. Creates pull request to default branch
7. Returns response with PR URL and number

#### `GetWorkflows(ctx, token, owner, repo) ([]File, error)`
Retrieves all workflows from a repository.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name

**Returns**:
- `[]File`: List of workflow files
- `error`: Error if fetch fails

#### `GetWorkflowContent(ctx, token, owner, repo, filePath) (*FileContentResponse, error)`
Retrieves the content of a workflow file.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `filePath`: Workflow file path

**Returns**:
- `*FileContentResponse`: File content and metadata
- `error`: Error if fetch fails

**Validation**:
- File path must start with `.github/workflows/`
- File must have `.yml` or `.yaml` extension

#### `UpdateWorkflow(ctx, token, req) (*Response, error)`
Updates an existing workflow file and creates a PR.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `req`: Update request

**Returns**:
- `*Response`: Update response with PR details
- `error`: Error if update fails

**Process**:
1. Validates file path is a workflow file
2. Verifies repository exists
3. Gets default branch
4. Creates new branch: `update-workflow/{name}-{timestamp}`
5. Updates workflow file on new branch
6. Creates pull request to default branch
7. Returns response with PR URL and number

#### `isValidWorkflowName(name string) bool`
Validates workflow name format.

**Parameters**:
- `name`: Workflow name

**Returns**: `bool` true if valid

**Rules**:
- Only alphanumeric, dash, underscore
- Length: 1-255 characters

---

### **internal/domain/workflow/errors.go**

**Purpose**: Workflow domain errors.

**Error Variables**:
- `ErrInvalidWorkflowName`: Invalid workflow name format
- `ErrInvalidDeploymentType`: Invalid deployment type
- `ErrEC2CommonFieldsRequired`: EC2 common fields missing
- `ErrEC2ProjectsRequired`: EC2 projects missing
- `ErrKubernetesCommonFieldsRequired`: K8s common fields missing
- `ErrKubernetesProjectsRequired`: K8s projects missing
- `ErrTemplateGenerationFailed`: Template generation failed

---

### **internal/domain/organization/models.go**

**Purpose**: Organization domain models.

**Structures**:

#### `Organization`
GitHub organization.

**Fields**:
- `ID int64`: Organization ID
- `Login string`: Organization name
- `AvatarURL string`: Avatar URL
- `Description string`: Organization description

#### `Repository`
Repository under organization.

**Fields**:
- `ID int64`: Repository ID
- `Name string`: Repository name
- `FullName string`: Full name (owner/repo)
- `Description string`: Repository description
- `Private bool`: Is private
- `HTMLURL string`: Web URL
- `DefaultBranch string`: Default branch

---

### **internal/domain/organization/service.go**

**Purpose**: Organization business logic service.

**Structure**:

#### `Service`
Organization service.

**Fields**:
- `githubOrg *github.OrganizationClient`: GitHub org client

**Functions**:

#### `NewService() *Service`
Creates a new organization service.

#### `GetUserOrganizations(ctx, token) ([]Organization, error)`
Retrieves user's organizations.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token

**Returns**:
- `[]Organization`: List of organizations
- `error`: Error if fetch fails

#### `GetOrganizationRepositories(ctx, token, orgName) ([]Repository, error)`
Retrieves organization repositories.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `orgName`: Organization name

**Returns**:
- `[]Repository`: List of repositories
- `error`: Error if fetch fails

#### `GetUserRepositories(ctx, token) (map[string][]Repository, error)`
Retrieves all user repositories grouped by organization.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token

**Returns**:
- `map[string][]Repository`: Repositories by organization
- `error`: Error if fetch fails

---

### **internal/domain/repository/models.go**

**Purpose**: Repository domain models.

**Structures**:

#### `Branch`
Repository branch.

**Fields**:
- `Name string`: Branch name
- `CommitSHA string`: Latest commit SHA
- `Protected bool`: Is protected

#### `TagReference`
Git tag reference.

**Fields**:
- `Ref string`: Reference name
- `ObjectSHA string`: Object SHA
- `URL string`: Reference URL

#### `Package`
GitHub package.

**Fields**:
- `ID int64`: Package ID
- `Name string`: Package name
- `PackageType string`: Package type
- `Visibility string`: public/private
- `URL string`: API URL
- `HTMLURL string`: Web URL
- `CreatedAt string`: Creation timestamp
- `UpdatedAt string`: Update timestamp
- `OwnerLogin string`: Owner username
- `RepositoryName string`: Repository name

---

### **internal/domain/repository/service.go**

**Purpose**: Repository business logic service.

**Structure**:

#### `Service`
Repository service.

**Fields**:
- `githubRepo *github.RepositoryClient`: GitHub repo client

**Functions**:

#### `NewService() *Service`
Creates a new repository service.

#### `GetBranches(ctx, token, owner, repo) ([]Branch, error)`
Retrieves repository branches.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name

**Returns**:
- `[]Branch`: List of branches
- `error`: Error if fetch fails

#### `GetCommits(ctx, token, owner, repo, branch, perPage) ([]interface{}, error)`
Retrieves commits for a branch.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `branch`: Branch name
- `perPage`: Results per page

**Returns**:
- `[]interface{}`: List of commits
- `error`: Error if fetch fails

#### `GetTags(ctx, token, owner, repo) ([]interface{}, error)`
Retrieves repository tags.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name

**Returns**:
- `[]interface{}`: List of tags
- `error`: Error if fetch fails

#### `CreateTag(ctx, token, owner, repo, tagName, commitSHA) (*TagReference, error)`
Creates a new tag.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `tagName`: Tag name
- `commitSHA`: Commit SHA to tag

**Returns**:
- `*TagReference`: Created tag reference
- `error`: Error if creation fails

#### `GetWorkflowRuns(ctx, token, owner, repo, perPage) ([]interface{}, error)`
Retrieves workflow runs.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `perPage`: Results per page

**Returns**:
- `[]interface{}`: List of workflow runs
- `error`: Error if fetch fails

#### `GetWorkflowRunDetail(ctx, token, owner, repo, runID) (interface{}, []interface{}, error)`
Retrieves workflow run details.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `runID`: Workflow run ID

**Returns**:
- `interface{}`: Run details
- `[]interface{}`: Run jobs
- `error`: Error if fetch fails

#### `GetJobLogs(ctx, token, owner, repo, jobID) (string, error)`
Retrieves job logs.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `owner`: Repository owner
- `repo`: Repository name
- `jobID`: Job ID

**Returns**:
- `string`: Job logs
- `error`: Error if fetch fails

#### `GetUserPackages(ctx, token, packageType) ([]Package, error)`
Retrieves packages for the authenticated user.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `packageType`: Package type filter (npm, docker, etc.)

**Returns**:
- `[]Package`: List of packages
- `error`: Error if fetch fails

#### `GetOrgPackages(ctx, token, org, packageType) ([]Package, error)`
Retrieves packages for an organization.

**Parameters**:
- `ctx`: Context
- `token`: GitHub access token
- `org`: Organization name
- `packageType`: Package type filter

**Returns**:
- `[]Package`: List of packages
- `error`: Error if fetch fails

---

## Function Reference

### Authentication Flow

1. **User initiates login**: `GET /api/auth/github`
   - `auth.Login()` generates state and redirects to GitHub

2. **GitHub redirects back**: `GET /api/auth/github/callback?code=...&state=...`
   - `auth.Callback()` validates state, exchanges code for token
   - Creates/updates user and token in database
   - Generates JWT and redirects to frontend

3. **Frontend stores JWT**: Client stores token in localStorage/cookie

4. **Authenticated requests**: Include `Authorization: Bearer <JWT>` header
   - `AuthMiddleware()` validates JWT
   - Sets user context for handlers

5. **Get user profile**: `GET /api/auth/me`
   - `auth.GetProfile()` returns user data

6. **Logout**: `POST /api/auth/logout`
   - `auth.Logout()` returns success (client removes token)

### Workflow Creation Flow

1. **Generate workflow YAML**: `POST /api/workflows/preview`
   - `workflow.Preview()` generates YAML without creating file
   - Returns YAML for user review

2. **Create workflow**: `POST /api/workflows/create`
   - `workflow.Create()` validates request
   - Generates YAML using template generator
   - Creates branch and file in GitHub
   - Creates pull request
   - Returns PR details

3. **List workflows**: `GET /api/workflows/:owner/:repo`
   - `workflow.List()` fetches workflow files from `.github/workflows/`

4. **Get workflow content**: `GET /api/workflows/:owner/:repo/file?path=...`
   - `workflow.GetWorkflowContent()` fetches file content

5. **Update workflow**: `PUT /api/workflows/:owner/:repo/file`
   - `workflow.UpdateWorkflow()` creates branch
   - Updates file and creates PR

### Repository Operations Flow

1. **List organizations**: `GET /api/organizations`
   - `organization.List()` fetches user's organizations

2. **List org repositories**: `GET /api/organizations/:org/repositories`
   - `organization.GetRepositories()` fetches org repos

3. **List user repositories**: `GET /api/repositories`
   - `organization.GetUserRepositories()` fetches all repos by org

4. **Get branches**: `GET /api/repositories/:owner/:repo/branches`
   - `repository.GetBranches()` fetches branches

5. **Get commits**: `GET /api/repositories/:owner/:repo/branches/:branch/commits`
   - `repository.GetCommits()` fetches branch commits

6. **Get tags**: `GET /api/repositories/:owner/:repo/tags`
   - `repository.GetTags()` fetches repository tags

7. **Create tag**: `POST /api/repositories/tags`
   - `repository.CreateTag()` creates new tag

8. **Get workflow runs**: `GET /api/repositories/:owner/:repo/actions/runs`
   - `repository.GetWorkflowRuns()` fetches GitHub Actions runs

9. **Get run details**: `GET /api/repositories/:owner/:repo/actions/runs/:run_id`
   - `repository.GetWorkflowRunDetail()` fetches run details and jobs

10. **Get job logs**: `GET /api/repositories/:owner/:repo/actions/jobs/:job_id/logs`
    - `repository.GetJobLogs()` fetches job logs

11. **Get user packages**: `GET /api/packages/user?package_type=...`
    - `repository.GetUserPackages()` fetches user's packages

12. **Get org packages**: `GET /api/packages/org/:org?package_type=...`
    - `repository.GetOrgPackages()` fetches organization packages

---

## Environment Variables

### Required Variables
- `GITHUB_CLIENT_ID`: GitHub OAuth app client ID
- `GITHUB_CLIENT_SECRET`: GitHub OAuth app client secret
- `JWT_SECRET`: Secret key for signing JWT tokens
- `DB_PASSWORD`: PostgreSQL database password

### Optional Variables (with defaults)
- `PORT`: Server port (default: 8080)
- `GIN_MODE`: Gin mode (default: debug)
- `ENVIRONMENT`: Environment (default: development)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (default: postgres)
- `DB_NAME`: Database name (default: calance_workflow)
- `DB_SSLMODE`: SSL mode (default: disable)
- `GITHUB_REDIRECT_URL`: OAuth callback URL (default: http://localhost:8080/api/auth/github/callback)
- `JWT_EXPIRATION_HOURS`: JWT expiration (default: 24)
- `FRONTEND_URL`: Frontend URL (default: http://localhost:3000)
- `ALLOWED_ORIGINS`: CORS allowed origins (default: http://localhost:3000)
- `LOG_LEVEL`: Log level (default: info)
- `LOG_FORMAT`: Log format (default: json)

---

## Database Schema

### Users Table
```sql
CREATE TABLE users (
  id UUID PRIMARY KEY,
  github_id BIGINT UNIQUE NOT NULL,
  username VARCHAR NOT NULL,
  email VARCHAR,
  avatar_url VARCHAR,
  name VARCHAR,
  bio VARCHAR,
  location VARCHAR,
  company VARCHAR,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
```

### Tokens Table
```sql
CREATE TABLE tokens (
  id UUID PRIMARY KEY,
  user_id UUID UNIQUE NOT NULL REFERENCES users(id),
  access_token VARCHAR NOT NULL,
  token_type VARCHAR,
  scope VARCHAR,
  expires_at TIMESTAMP,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  deleted_at TIMESTAMP
);
```

---

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error message",
  "error": "Detailed error"
}
```

---

## Deployment

### Development
```bash
# Load environment variables
cp .env.example .env
# Edit .env with your credentials

# Run database
docker-compose up -d

# Run server
go run cmd/server/main.go
```

### Production
```bash
# Set environment variables
export GIN_MODE=release
export ENVIRONMENT=production
export LOG_FORMAT=json

# Build binary
go build -o calance-server cmd/server/main.go

# Run server
./calance-server
```

---

## Testing

### Health Check
```bash
curl http://localhost:8080/ping
```

### Authenticated Request
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/auth/me
```

---

## Security Considerations

1. **JWT Secret**: Use a strong, random secret in production
2. **HTTPS**: Always use HTTPS in production
3. **CORS**: Configure allowed origins carefully
4. **Database**: Use SSL mode in production
5. **Tokens**: GitHub tokens are stored encrypted in database
6. **State Parameter**: CSRF protection for OAuth flow
7. **HTTP-Only Cookies**: State cookie is HTTP-only
8. **Token Expiration**: JWT tokens expire after configured hours

---

## Contributing

1. Follow Go best practices
2. Maintain three-layer architecture
3. Add tests for new features
4. Update documentation
5. Use conventional commits

---

## License

---

**Generated**: 2026-02-06  
**Version**: 1.0.0  
**Maintainer**: Vagish Maurya
