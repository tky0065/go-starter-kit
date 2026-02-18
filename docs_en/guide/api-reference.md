# API Reference

<div class="navigation">
  <a href="index.md"><i class="material-icons">arrow_back</i> Guide Index</a>
</div>

---

## Available Endpoints

#### Health Check

```
GET /health
```

**Response** (200):
```json
{
  "status": "ok"
}
```

---

### Authentication

#### Register

```
POST /api/v1/auth/register
Content-Type: application/json
```

**Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Validation**:
- email: required, valid email, max 255 chars
- password: required, min 8 chars, max 72 chars

**Response** (201):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `409 Conflict`: Email already exists

**Curl example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

---

#### Login

```
POST /api/v1/auth/login
Content-Type: application/json
```

**Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials

**Curl example**:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

---

#### Refresh Token

```
POST /api/v1/auth/refresh
Content-Type: application/json
```

**Body**:
```json
{
  "refresh_token": "eyJhbGc..."
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "access_token": "eyJhbGc...",
    "token_type": "Bearer",
    "expires_in": 900
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid or expired refresh token

**Curl example**:
```bash
REFRESH_TOKEN="<refresh_token_from_login>"
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}"
```

---

### Users (Protected)

All user endpoints require a valid JWT token.

#### List Users

```
GET /api/v1/users
Authorization: Bearer <access_token>
```

**Response** (200):
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "email": "user1@example.com",
      "created_at": "2026-01-09T10:00:00Z",
      "updated_at": "2026-01-09T10:00:00Z"
    },
    {
      "id": 2,
      "email": "user2@example.com",
      "created_at": "2026-01-09T11:00:00Z",
      "updated_at": "2026-01-09T11:00:00Z"
    }
  ]
}
```

**Errors**:
- `401 Unauthorized`: Missing or invalid token

**Curl example**:
```bash
TOKEN="<access_token>"
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

---

#### Get User by ID

```
GET /api/v1/users/:id
Authorization: Bearer <access_token>
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "created_at": "2026-01-09T10:00:00Z",
    "updated_at": "2026-01-09T10:00:00Z"
  }
}
```

**Errors**:
- `401 Unauthorized`: Invalid token
- `404 Not Found`: User not found

**Curl example**:
```bash
TOKEN="<access_token>"
curl -X GET http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $TOKEN"
```

---

#### Update User

```
PUT /api/v1/users/:id
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Body**:
```json
{
  "email": "newemail@example.com"
}
```

**Response** (200):
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "email": "newemail@example.com",
    "created_at": "2026-01-09T10:00:00Z",
    "updated_at": "2026-01-10T15:30:00Z"
  }
}
```

**Errors**:
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid token
- `404 Not Found`: User not found
- `409 Conflict`: Email already exists

---

#### Delete User

```
DELETE /api/v1/users/:id
Authorization: Bearer <access_token>
```

**Response** (200):
```json
{
  "status": "success",
  "message": "User deleted successfully"
}
```

**Errors**:
- `401 Unauthorized`: Invalid token
- `404 Not Found`: User not found

**Note**: Uses soft delete (DeletedAt), data remains in the DB.

---

### Complete API Workflow

```bash
# 1. Register
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}')

# Extract access_token
ACCESS_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo $REGISTER_RESP | jq -r '.data.refresh_token')

# 2. List users (with token)
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $ACCESS_TOKEN"

# 3. Update user
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"updated@example.com"}'

# 4. When access token expires (15min), use refresh token
NEW_ACCESS=$(curl -s -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"$REFRESH_TOKEN\"}" | jq -r '.data.access_token')

# 5. Continue with new token
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $NEW_ACCESS"
```

---


---

## Navigation

**Previous**: [Practical Examples](examples.md)  
**Next**: [Testing](testing.md)  
**Index**: [Guide Index](index.md)
