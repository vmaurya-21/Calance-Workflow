# Calance Workflow Backend

A Go backend server with GitHub OAuth authentication, JWT-based sessions, and PostgreSQL database.

## 🚀 Features

- ✅ GitHub OAuth 2.0 authentication
- ✅ JWT-based session management
- ✅ PostgreSQL database with GORM
- ✅ RESTful API architecture
- ✅ Clean architecture (controllers, services, repositories)
- ✅ CORS support for frontend integration
- ✅ Environment-based configuration

## 📋 Prerequisites

- Go 1.25.4 or higher
- PostgreSQL 12 or higher
- GitHub account (for OAuth app creation)

## 🛠️ Installation

1. **Clone the repository**:
   ```bash
   git clone <your-repo-url>
   cd Calance-Workflow-backend
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up PostgreSQL**:
   ```bash
   # Create database
   psql -U postgres
   CREATE DATABASE calance_workflow;
   \q
   ```

4. **Configure environment variables**:
   ```bash
   # Copy example env file
   cp .env.example .env
   
   # Edit .env and add your credentials
   ```

5. **Create GitHub OAuth App**:
   - Go to [GitHub Developer Settings](https://github.com/settings/developers)
   - Click "New OAuth App"
   - Set **Authorization callback URL** to: `http://localhost:8080/api/auth/github/callback`
   - Copy the Client ID and Client Secret to your `.env` file

6. **Run the server**:
   ```bash
   go run cmd/server/main.go
   ```

The server will start on `http://localhost:8080`

## 📁 Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go              # Configuration management
│   ├── controllers/
│   │   └── auth_controller.go     # HTTP request handlers
│   ├── database/
│   │   └── database.go            # Database connection
│   ├── middleware/
│   │   └── auth_middleware.go     # JWT validation & CORS
│   ├── models/
│   │   └── user.go                # Database models
│   ├── repositories/
│   │   └── user_repository.go     # Data access layer
│   ├── router/
│   │   └── router.go              # Route configuration
│   ├── services/
│   │   └── github_oauth_service.go # Business logic
│   └── utils/
│       ├── jwt.go                 # JWT utilities
│       └── response.go            # API response helpers
├── .env                           # Environment variables (not committed)
├── .env.example                   # Environment template
├── .gitignore
├── go.mod
├── go.sum
├── GITHUB_OAUTH_GUIDE.md         # Frontend integration guide
└── README.md
```

## 🔑 Environment Variables

See `.env.example` for all available configuration options. Key variables:

```env
# GitHub OAuth
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret

# JWT
JWT_SECRET=your_secure_random_secret

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=calance_workflow

# Frontend
FRONTEND_URL=http://localhost:3000
ALLOWED_ORIGINS=http://localhost:3000
```

## 📡 API Endpoints

### Public Endpoints

- `GET /ping` - Health check
- `GET /api/auth/github` - Initiate GitHub OAuth login
- `GET /api/auth/github/callback` - OAuth callback handler

### Protected Endpoints (Require JWT)

- `GET /api/auth/me` - Get current user
- `POST /api/auth/logout` - Logout

## 🎨 Frontend Integration

See [GITHUB_OAUTH_GUIDE.md](./GITHUB_OAUTH_GUIDE.md) for complete frontend integration instructions with React, Vue, and vanilla JavaScript examples.

### Quick Start:

1. Redirect user to: `http://localhost:8080/api/auth/github`
2. Handle callback at your frontend route: `/auth/callback?token=<JWT>`
3. Store token and use in Authorization header: `Bearer <token>`

## 🧪 Testing

Test the health endpoint:
```bash
curl http://localhost:8080/ping
```

Test authenticated endpoint:
```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/auth/me
```

## 🚀 Deployment

1. Update environment variables for production
2. Set `GIN_MODE=release` in `.env`
3. Update GitHub OAuth callback URL to production domain
4. Use HTTPS in production
5. Set strong JWT secret
6. Configure production database

## 📚 Tech Stack

- **Framework**: Gin Web Framework
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: OAuth 2.0 (GitHub) + JWT
- **Language**: Go 1.25.4

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 📄 License

[Your License Here]

## 📞 Support

For issues and questions, please open a GitHub issue.
