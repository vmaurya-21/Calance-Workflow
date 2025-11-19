# GitHub OAuth Integration Guide

## Backend Setup Complete ✅

The GitHub OAuth system has been fully implemented in your Go backend. Here's how to set it up and integrate with your frontend.

---

## 🔧 Backend Configuration

### 1. Set Up GitHub OAuth App

1. Go to [GitHub Developer Settings](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the details:
   - **Application name**: Calance Workflow
   - **Homepage URL**: `http://localhost:3000` (or your frontend URL)
   - **Authorization callback URL**: `http://localhost:8080/api/auth/github/callback`
4. Click "Register application"
5. Copy the **Client ID** and generate a **Client Secret**

### 2. Configure Environment Variables

Update your `.env` file with the GitHub credentials:

```env
# GitHub OAuth Configuration
GITHUB_CLIENT_ID=your_actual_github_client_id_here
GITHUB_CLIENT_SECRET=your_actual_github_client_secret_here
GITHUB_REDIRECT_URL=http://localhost:8080/api/auth/github/callback

# JWT Configuration (generate a secure random string)
JWT_SECRET=your_secure_random_jwt_secret_here

# Frontend Configuration
FRONTEND_URL=http://localhost:3000
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Set Up PostgreSQL Database

Install PostgreSQL and create a database:

```bash
# Using psql
psql -U postgres
CREATE DATABASE calance_workflow;
```

Update database credentials in `.env` if needed:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_postgres_password
DB_NAME=calance_workflow
```

### 4. Run the Backend

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` and automatically:
- Connect to PostgreSQL
- Create the users table
- Set up all OAuth routes

---

## 📡 API Endpoints

### Public Endpoints

#### 1. **Health Check**
```http
GET /ping
```

**Response:**
```json
{
  "message": "pong",
  "status": "healthy"
}
```

#### 2. **GitHub Login (Initiate OAuth)**
```http
GET /api/auth/github
```

**Description:** Redirects user to GitHub OAuth authorization page.

**Frontend Usage:**
```javascript
// Redirect user to backend OAuth endpoint
window.location.href = 'http://localhost:8080/api/auth/github';
```

#### 3. **GitHub Callback**
```http
GET /api/auth/github/callback?code=xxx&state=xxx
```

**Description:** GitHub redirects here after user authorization. Backend processes the OAuth flow and redirects to frontend with JWT token.

**Redirect to Frontend:**
```
http://localhost:3000/auth/callback?token=<JWT_TOKEN>
```

### Protected Endpoints (Require Authentication)

#### 4. **Get Current User**
```http
GET /api/auth/me
Authorization: Bearer <JWT_TOKEN>
```

**Response:**
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "uuid-string",
    "github_id": 12345678,
    "username": "johndoe",
    "email": "john@example.com",
    "avatar_url": "https://avatars.githubusercontent.com/u/12345678",
    "name": "John Doe",
    "bio": "Software Developer",
    "location": "San Francisco, CA",
    "company": "Tech Corp",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### 5. **Logout**
```http
POST /api/auth/logout
Authorization: Bearer <JWT_TOKEN>
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

---

## 🎨 Frontend Integration

### React Example

```javascript
import { useState, useEffect } from 'react';
import axios from 'axios';

const API_BASE_URL = 'http://localhost:8080/api';

// 1. Login Component
function Login() {
  const handleGitHubLogin = () => {
    window.location.href = 'http://localhost:8080/api/auth/github';
  };

  return (
    <button onClick={handleGitHubLogin}>
      Login with GitHub
    </button>
  );
}

// 2. Callback Handler (route: /auth/callback)
function AuthCallback() {
  useEffect(() => {
    // Get token from URL query params
    const params = new URLSearchParams(window.location.search);
    const token = params.get('token');
    
    if (token) {
      // Store token in localStorage
      localStorage.setItem('authToken', token);
      
      // Redirect to dashboard or home
      window.location.href = '/dashboard';
    } else {
      // Handle error
      window.location.href = '/login?error=auth_failed';
    }
  }, []);

  return <div>Processing authentication...</div>;
}

// 3. Axios Instance with Auth Token
const api = axios.create({
  baseURL: API_BASE_URL,
});

// Add token to all requests
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('authToken');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// 4. Fetch Current User
async function getCurrentUser() {
  try {
    const response = await api.get('/auth/me');
    return response.data.data; // User object
  } catch (error) {
    console.error('Failed to fetch user:', error);
    // Token might be expired, redirect to login
    if (error.response?.status === 401) {
      localStorage.removeItem('authToken');
      window.location.href = '/login';
    }
    throw error;
  }
}

// 5. Logout
async function logout() {
  try {
    await api.post('/auth/logout');
  } catch (error) {
    console.error('Logout error:', error);
  } finally {
    localStorage.removeItem('authToken');
    window.location.href = '/login';
  }
}

// 6. Protected Route Component
function Dashboard() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getCurrentUser()
      .then(setUser)
      .catch(() => {
        // Redirect to login on error
        window.location.href = '/login';
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h1>Welcome, {user.name}!</h1>
      <img src={user.avatar_url} alt={user.username} />
      <button onClick={logout}>Logout</button>
    </div>
  );
}
```

### Vue Example

```javascript
// Login.vue
<template>
  <button @click="loginWithGitHub">Login with GitHub</button>
</template>

<script>
export default {
  methods: {
    loginWithGitHub() {
      window.location.href = 'http://localhost:8080/api/auth/github';
    }
  }
}
</script>

// AuthCallback.vue
<template>
  <div>Processing authentication...</div>
</template>

<script>
export default {
  mounted() {
    const token = this.$route.query.token;
    if (token) {
      localStorage.setItem('authToken', token);
      this.$router.push('/dashboard');
    } else {
      this.$router.push('/login');
    }
  }
}
</script>

// axios setup (plugins/axios.js)
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api'
});

api.interceptors.request.use(config => {
  const token = localStorage.getItem('authToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
```

### Vanilla JavaScript Example

```javascript
// Login
function loginWithGitHub() {
  window.location.href = 'http://localhost:8080/api/auth/github';
}

// Callback Handler (on /auth/callback page)
const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('token');

if (token) {
  localStorage.setItem('authToken', token);
  window.location.href = '/dashboard';
}

// Make authenticated request
async function fetchUser() {
  const token = localStorage.getItem('authToken');
  
  const response = await fetch('http://localhost:8080/api/auth/me', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  if (!response.ok) {
    throw new Error('Unauthorized');
  }
  
  const data = await response.json();
  return data.data; // User object
}

// Logout
async function logout() {
  const token = localStorage.getItem('authToken');
  
  await fetch('http://localhost:8080/api/auth/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  localStorage.removeItem('authToken');
  window.location.href = '/login';
}
```

---

## 🔒 Security Best Practices

1. **HTTPS in Production**: Always use HTTPS for production deployments
2. **Secure JWT Secret**: Use a strong, random JWT secret (at least 32 characters)
3. **Token Expiration**: JWT tokens expire after 24 hours (configurable via `JWT_EXPIRATION_HOURS`)
4. **CORS**: Configure `ALLOWED_ORIGINS` to only include your frontend domain
5. **Environment Variables**: Never commit `.env` file to version control
6. **Token Storage**: Consider using httpOnly cookies instead of localStorage for better security

---

## 🔄 OAuth Flow Diagram

```
┌─────────┐         ┌──────────┐         ┌────────┐         ┌─────────┐
│Frontend │         │ Backend  │         │ GitHub │         │Database │
└────┬────┘         └────┬─────┘         └───┬────┘         └────┬────┘
     │                   │                   │                   │
     │ 1. Click Login    │                   │                   │
     ├──────────────────>│                   │                   │
     │                   │                   │                   │
     │ 2. Redirect to    │                   │                   │
     │    GitHub OAuth   │                   │                   │
     │<──────────────────┤                   │                   │
     │                   │                   │                   │
     │ 3. User Authorizes│                   │                   │
     ├───────────────────┼──────────────────>│                   │
     │                   │                   │                   │
     │                   │ 4. Callback with  │                   │
     │                   │    auth code      │                   │
     │                   │<──────────────────┤                   │
     │                   │                   │                   │
     │                   │ 5. Exchange code  │                   │
     │                   │    for token      │                   │
     │                   ├──────────────────>│                   │
     │                   │                   │                   │
     │                   │ 6. Access token + │                   │
     │                   │    user data      │                   │
     │                   │<──────────────────┤                   │
     │                   │                   │                   │
     │                   │ 7. Create/Update User                 │
     │                   ├───────────────────┼──────────────────>│
     │                   │                   │                   │
     │                   │ 8. User saved     │                   │
     │                   │<──────────────────┼───────────────────┤
     │                   │                   │                   │
     │ 9. Redirect with  │                   │                   │
     │    JWT token      │                   │                   │
     │<──────────────────┤                   │                   │
     │                   │                   │                   │
     │ 10. Store token   │                   │                   │
     │     in localStorage│                  │                   │
     │                   │                   │                   │
     │ 11. Make authenticated requests       │                   │
     │    with Bearer token                  │                   │
     ├──────────────────>│                   │                   │
```

---

## 🧪 Testing the Integration

1. **Start Backend**:
   ```bash
   go run cmd/server/main.go
   ```

2. **Test Health Endpoint**:
   ```bash
   curl http://localhost:8080/ping
   ```

3. **Test OAuth Flow**:
   - Visit `http://localhost:8080/api/auth/github` in browser
   - You'll be redirected to GitHub
   - After authorization, you'll be redirected to your frontend with a token

4. **Test Protected Endpoint**:
   ```bash
   curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
        http://localhost:8080/api/auth/me
   ```

---

## 🚀 Next Steps

1. ✅ Backend OAuth implementation complete
2. 📝 Create frontend routes: `/login`, `/auth/callback`, `/dashboard`
3. 🎨 Implement UI components for login and user profile
4. 🔐 Add protected routes in frontend
5. 🚀 Deploy to production (update URLs in `.env` and GitHub OAuth app)

---

## 📞 Troubleshooting

### "GITHUB_CLIENT_ID is required"
- Ensure `.env` file exists and contains valid GitHub OAuth credentials

### "Database connection failed"
- Check PostgreSQL is running: `pg_isready`
- Verify database credentials in `.env`
- Create database if it doesn't exist

### CORS Errors
- Add your frontend URL to `ALLOWED_ORIGINS` in `.env`
- Restart the backend after changing `.env`

### "Invalid token" / "Token has expired"
- JWT tokens expire after 24 hours
- User needs to re-authenticate
- Clear `authToken` from localStorage and redirect to login

---

## 📚 Additional Resources

- [GitHub OAuth Documentation](https://docs.github.com/en/developers/apps/building-oauth-apps)
- [JWT Best Practices](https://tools.ietf.org/html/rfc8725)
- [GORM Documentation](https://gorm.io/docs/)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)
