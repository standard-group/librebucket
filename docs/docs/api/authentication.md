# Authentication API

The Authentication API provides endpoints for user registration and login. Authentication is required for most Librebucket operations.

## Authentication Methods

Librebucket supports token-based authentication. After successful login or registration, you'll receive a token that must be included in subsequent API requests.

### Token Formats

Librebucket accepts tokens in three ways:

1. **Authorization Header (Recommended):**
   ```
   Authorization: Bearer <token>
   ```

2. **Custom Header:**
   ```
   X-Auth-Token: <token>
   ```

3. **Query Parameter:**
   ```
   ?token=<token>
   ```

## Endpoints

### Register User

Register a new user account.

**Endpoint:** `POST /api/v1/users/register`

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Parameters:**
- `username` (string, required): Username for the new account (3-50 characters)
- `password` (string, required): Password for the new account (minimum 8 characters)

**Response (201 Created):**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "alice",
    "is_admin": false,
    "created_at": "2024-01-01T12:00:00Z"
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "secure_password123"
  }'
```

**Errors:**
- `400 Bad Request` - Invalid JSON or missing fields
- `409 Conflict` - Username already exists

### Login User

Authenticate an existing user and receive a token.

**Endpoint:** `POST /api/v1/users/login`

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Parameters:**
- `username` (string, required): Existing username
- `password` (string, required): User's password

**Response (200 OK):**
```json
{
  "status": "success",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "alice",
    "is_admin": false,
    "last_login": "2024-01-01T12:00:00Z"
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:3000/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "alice",
    "password": "secure_password123"
  }'
```

**Errors:**
- `400 Bad Request` - Invalid JSON or missing fields
- `401 Unauthorized` - Invalid username or password

### Validate Token

Check if a token is valid and get user information.

**Endpoint:** `GET /api/v1/auth/validate`

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "status": "success",
  "user": {
    "id": 1,
    "username": "alice",
    "is_admin": false,
    "token_expires": "2024-01-08T12:00:00Z"
  }
}
```

**Example:**
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:3000/api/v1/auth/validate
```

**Errors:**
- `401 Unauthorized` - Invalid or expired token

## Token Management

### Token Lifetime

Tokens are currently long-lived and don't expire automatically. This may change in future versions to include configurable expiration times.

### Token Security

- Tokens are generated using cryptographically secure random bytes
- Store tokens securely and never commit them to version control
- Consider them as sensitive as passwords

### Revoking Tokens

Currently, there's no explicit token revocation endpoint. Tokens remain valid until:
- The user changes their password
- The user is deleted
- The server secret key changes

## Authentication Examples

### JavaScript/Node.js

```javascript
class LibrebucketAuth {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
    this.token = null;
  }

  async register(username, password) {
    const response = await fetch(`${this.baseUrl}/api/v1/users/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });

    const data = await response.json();
    if (data.status === 'success') {
      this.token = data.token;
      localStorage.setItem('librebucket_token', data.token);
    }
    return data;
  }

  async login(username, password) {
    const response = await fetch(`${this.baseUrl}/api/v1/users/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ username, password })
    });

    const data = await response.json();
    if (data.status === 'success') {
      this.token = data.token;
      localStorage.setItem('librebucket_token', data.token);
    }
    return data;
  }

  getAuthHeaders() {
    return {
      'Authorization': `Bearer ${this.token}`,
      'Content-Type': 'application/json'
    };
  }
}

// Usage
const auth = new LibrebucketAuth('http://localhost:3000');
await auth.login('alice', 'password123');
```

### Python

```python
import requests
import json

class LibrebucketAuth:
    def __init__(self, base_url):
        self.base_url = base_url
        self.token = None

    def register(self, username, password):
        url = f"{self.base_url}/api/v1/users/register"
        data = {"username": username, "password": password}
        
        response = requests.post(url, json=data)
        result = response.json()
        
        if result.get('status') == 'success':
            self.token = result['token']
            
        return result

    def login(self, username, password):
        url = f"{self.base_url}/api/v1/users/login"
        data = {"username": username, "password": password}
        
        response = requests.post(url, json=data)
        result = response.json()
        
        if result.get('status') == 'success':
            self.token = result['token']
            
        return result

    def get_auth_headers(self):
        return {
            'Authorization': f'Bearer {self.token}',
            'Content-Type': 'application/json'
        }

# Usage
auth = LibrebucketAuth('http://localhost:3000')
result = auth.login('alice', 'password123')
```

### Go

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type AuthClient struct {
    BaseURL string
    Token   string
}

type LoginRequest struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type AuthResponse struct {
    Status string `json:"status"`
    Token  string `json:"token"`
    User   struct {
        ID       int    `json:"id"`
        Username string `json:"username"`
        IsAdmin  bool   `json:"is_admin"`
    } `json:"user"`
}

func (c *AuthClient) Login(username, password string) (*AuthResponse, error) {
    loginReq := LoginRequest{Username: username, Password: password}
    jsonData, err := json.Marshal(loginReq)
    if err != nil {
        return nil, err
    }

    resp, err := http.Post(
        c.BaseURL+"/api/v1/users/login",
        "application/json",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var authResp AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        return nil, err
    }

    if authResp.Status == "success" {
        c.Token = authResp.Token
    }

    return &authResp, nil
}

func (c *AuthClient) AuthHeader() string {
    return "Bearer " + c.Token
}
```

## Security Best Practices

### Password Requirements

While Librebucket doesn't enforce complex password requirements by default, consider implementing these in your application:

- Minimum 8 characters
- Mix of uppercase, lowercase, numbers, and symbols
- No common dictionary words
- No username in password

### Token Storage

**Client-side Applications:**
- Use secure storage (Keychain on macOS, Credential Manager on Windows)
- For web apps, consider secure HTTP-only cookies instead of localStorage

**Server-side Applications:**
- Store in environment variables
- Use secure configuration management
- Never log tokens

### HTTPS Only

Always use HTTPS in production to protect tokens in transit:

```bash
# Bad - tokens visible in network traffic
curl -H "Authorization: Bearer TOKEN" http://example.com/api/v1/repos

# Good - tokens encrypted in transit
curl -H "Authorization: Bearer TOKEN" https://example.com/api/v1/repos
```

## Error Handling

### Common Authentication Errors

```json
// Missing authentication
{
  "status": "error",
  "message": "Authentication required",
  "code": "AUTH_REQUIRED"
}

// Invalid token
{
  "status": "error",
  "message": "Invalid authentication token",
  "code": "INVALID_TOKEN"
}

// Username already exists
{
  "status": "error",
  "message": "Username already exists",
  "code": "USERNAME_EXISTS"
}

// Invalid credentials
{
  "status": "error",
  "message": "Invalid username or password",
  "code": "INVALID_CREDENTIALS"
}
```

### Handling Authentication in Your App

```javascript
async function makeAuthenticatedRequest(url, options = {}) {
  const token = localStorage.getItem('librebucket_token');
  
  if (!token) {
    throw new Error('No authentication token found');
  }

  const response = await fetch(url, {
    ...options,
    headers: {
      ...options.headers,
      'Authorization': `Bearer ${token}`
    }
  });

  if (response.status === 401) {
    // Token expired or invalid, redirect to login
    localStorage.removeItem('librebucket_token');
    window.location.href = '/login';
    return;
  }

  return response.json();
}
```

## Next Steps

- [User Management API](users.md) - Managing user accounts
- [Repository API](repositories.md) - Creating and managing repositories
- [Deployment Security](../deployment/security.md) - Production security considerations
