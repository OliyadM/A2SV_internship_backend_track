# Task Manager API Documentation

## Base URL
`http://localhost:8080`

## Authentication
This API uses JWT (JSON Web Token) for authentication. Include the token in the `Authorization` header for protected routes.

### 1. User Authentication
#### Register a User
**Endpoint:** `POST /register`

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response (Success):**
```json
{
  "message": "User created"
}
```
📌 *The first registered user is assigned admin privileges automatically.*

#### Login (Get JWT Token)
**Endpoint:** `POST /login`

**Request:**
```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Response:**
```json
{
  "token": "eyJhbGci...<your-jwt-token>"
}
```

📌 *Use the token in the Authorization header for authenticated requests:*
```
Authorization: Bearer <token>
```

## 2. Task Management
### Create a Task *(Admin Only)*
**Endpoint:** `POST /tasks`

**Request:**
```json
{
  "title": "Fix Bug",
  "description": "Critical login issue",
  "due_date": "2025-03-30T12:00:00Z",
  "status": "urgent"
}
```

**Response:**
```json
{
  "message": "Task created",
  "id": "65a1b2c3d4e5f6g7h8i9j0"
}
```

### Get All Tasks
**Endpoint:** `GET /tasks`

**Response:**
```json
{
  "tasks": [
    {
      "_id": "65a1b2c3d4e5f6g7h8i9j0",
      "title": "Fix Bug",
      "description": "...",
      "due_date": "2025-03-30T12:00:00Z",
      "status": "urgent"
    }
  ]
}
```

### Get a Task by ID
**Endpoint:** `GET /tasks/:id`

**Response:**
```json
{
  "_id": "65a1b2c3d4e5f6g7h8i9j0",
  "title": "Fix Bug",
  "description": "...",
  "due_date": "2025-03-30T12:00:00Z",
  "status": "urgent"
}
```

### Update a Task *(Admin Only)*
**Endpoint:** `PUT /tasks/:id`

**Request:**
```json
{
  "status": "completed"
}
```

**Response:**
```json
{
  "message": "Task updated"
}
```

### Delete a Task *(Admin Only)*
**Endpoint:** `DELETE /tasks/:id`

**Response:**
```json
{
  "message": "Task removed"
}
```

## 3. Admin Actions
### Promote a User to Admin
**Endpoint:** `POST /promote`

**Request:**
```json
{
  "username": "regular_user"
}
```

**Response:**
```json
{
  "message": "User promoted to admin"
}
```

## Error Responses
| Code | Error | Description |
|------|-------|-------------|
| 400  | Bad Request | Invalid JSON/data provided |
| 401  | Unauthorized | Missing or invalid JWT token |
| 403  | Forbidden | Insufficient privileges (non-admin access) |
| 404  | Not Found | Task or user not found |
| 500  | Server Error | Internal server issue |

## Task Status Options
- `pending` (default)
- `in_progress`
- `completed`
- `urgent`

📌 **Notes:**
- Replace `<token>` with the JWT obtained from `/login`.
- All dates follow ISO 8601 format (`YYYY-MM-DDTHH:MM:SSZ`).
- Admin privileges are required for:
  - Task creation, modification, and deletion.
  - User promotion.

