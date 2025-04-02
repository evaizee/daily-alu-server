# DailyAlu API Documentation

## Overview
This document provides comprehensive documentation for the DailyAlu API, a RESTful API for managing baby activities, children profiles, and user accounts. The API follows RESTful principles and uses JWT for authentication.

## Base URL
```
{{base_url}}/api/v1
```

## Authentication
Most endpoints require authentication using JWT tokens. The token should be included in the Authorization header as a Bearer token:
```
Authorization: Bearer {{access_token}}
```

Additionally, all requests require an API key in the header:
```
X-API-Key: {{api_key}}
```

## Standard Response Format

### Success Response
All successful responses follow this format:
```json
{
  "success": true,
  "message": "Human-readable success message",
  "data": { ... } // Response data or null
}
```

### Paginated Response
Endpoints that return multiple items use pagination:
```json
{
  "success": true,
  "message": "Human-readable success message",
  "data": [ ... ], // Array of items
  "pagination": {
    "total": 25,
    "current_page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

### Error Response
All error responses follow this format:
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field1": ["Error details for field1"],
      "field2": ["Error details for field2"]
    }
  }
}
```

## Common Error Codes
- `BAD_REQUEST`: Invalid request parameters or body
- `UNAUTHORIZED`: Authentication required or invalid credentials
- `FORBIDDEN`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `VALIDATION_ERROR`: Request validation failed
- `INTERNAL_SERVER_ERROR`: Server error

---

## User Management

### Register
Creates a new user account.

- **URL**: `/auth/register`
- **Method**: `POST`
- **Auth Required**: No (API key only)
- **Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "confirm_password": "password123",
  "name": "John Doe"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Your account is almost ready! To unlock all of our features, please verify your email address.",
  "data": {
    "id": "user-id",
    "email": "user@example.com",
    "name": "John Doe",
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Login
Authenticates a user and returns access and refresh tokens.

- **URL**: `/auth/login`
- **Method**: `POST`
- **Auth Required**: No (API key only)
- **Request Body**:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user-id",
      "email": "user@example.com",
      "name": "John Doe",
      "last_login": "2025-03-28T07:43:04Z",
      "created_at": "2025-03-28T07:43:04Z",
      "updated_at": "2025-03-28T07:43:04Z"
    }
  }
}
```

### Verify Email
Verifies a user's email address using the verification token.

- **URL**: `/auth/verify-email/:token` or `/auth/verify-email?token=verification-token`
- **Method**: `GET`
- **Auth Required**: No (API key only)
- **Response**:
```json
{
  "success": true,
  "message": "Email verified successfully",
  "data": null
}
```

### Refresh Token
Refreshes the access token using a valid refresh token.

- **URL**: `/auth/refresh-token`
- **Method**: `POST`
- **Auth Required**: Yes (JWT + API key)
- **Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Forgot Password
Initiates the password recovery process.

- **URL**: `/auth/forgot-password`
- **Method**: `POST`
- **Auth Required**: No (API key only)
- **Request Body**:
```json
{
  "email": "user@example.com"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "If your email is registered with us, you will receive password reset instructions shortly",
  "data": null
}
```

### Reset Password
Resets a user's password using a reset token.

- **URL**: `/auth/reset-password`
- **Method**: `POST`
- **Auth Required**: No (API key only)
- **Request Body**:
```json
{
  "token": "reset-token",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Your password has been reset successfully",
  "data": null
}
```

### Get User
Retrieves a user's profile information.

- **URL**: `/users/:id`
- **Method**: `GET`
- **Auth Required**: Yes (JWT + API key)
- **Response**:
```json
{
  "success": true,
  "message": "User retrieved successfully",
  "data": {
    "id": "user-id",
    "email": "user@example.com",
    "name": "John Doe",
    "last_login": "2025-03-28T07:43:04Z",
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Update User
Updates a user's profile information.

- **URL**: `/users/:id`
- **Method**: `PUT`
- **Auth Required**: Yes (JWT + API key)
- **Request Body**:
```json
{
  "email": "updated@example.com",
  "name": "Updated Name"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": "user-id",
    "email": "updated@example.com",
    "name": "Updated Name",
    "last_login": "2025-03-28T07:43:04Z",
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Update Password
Updates a user's password.

- **URL**: `/users/:id/password`
- **Method**: `PATCH`
- **Auth Required**: Yes (JWT + API key)
- **Request Body**:
```json
{
  "old_password": "password123",
  "new_password": "newpassword123",
  "confirm_password": "newpassword123"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Password updated successfully",
  "data": null
}
```

### Delete User (Admin Only)
Deletes a user account.

- **URL**: `/users/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes (JWT + API key + Admin role)
- **Response**:
```json
{
  "success": true,
  "message": "User deleted successfully",
  "data": null
}
```

---

## Activities

### Create Activity
Creates a new activity record.

- **URL**: `/activities`
- **Method**: `POST`
- **Auth Required**: Yes (JWT + API key)
- **Request Body**:
```json
{
  "child_id": 1,
  "type": "feeding",
  "details": {
    "amount": 120,
    "unit": "ml",
    "notes": "Formula milk"
  },
  "happens_at": "2025-03-28T07:43:04Z"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Activity created successfully",
  "data": {
    "id": 1,
    "user_id": "user-id",
    "child_id": 1,
    "type": "feeding",
    "details": {
      "amount": 120,
      "unit": "ml",
      "notes": "Formula milk"
    },
    "happens_at": "2025-03-28T07:43:04Z",
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Get Activity
Retrieves a specific activity by ID.

- **URL**: `/activities/:id`
- **Method**: `GET`
- **Auth Required**: Yes (JWT + API key)
- **Response**:
```json
{
  "success": true,
  "message": "Activity retrieved successfully",
  "data": {
    "id": 1,
    "user_id": "user-id",
    "child_id": 1,
    "type": "feeding",
    "details": {
      "amount": 120,
      "unit": "ml",
      "notes": "Formula milk"
    },
    "happens_at": "2025-03-28T07:43:04Z",
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Update Activity
Updates an existing activity.

- **URL**: `/activities/:id`
- **Method**: `PUT`
- **Auth Required**: Yes (JWT + API key)
- **Request Body**:
```json
{
  "child_id": 1,
  "details": {
    "amount": 150,
    "unit": "ml",
    "notes": "Formula milk with cereal"
  },
  "happens_at": "2025-03-28T08:00:00Z"
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Activity updated successfully",
  "data": {
    "id": 1,
    "user_id": "user-id",
    "child_id": 1,
    "type": "feeding",
    "details": {
      "amount": 150,
      "unit": "ml",
      "notes": "Formula milk with cereal"
    },
    "happens_at": "2025-03-28T08:00:00Z",
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T08:00:00Z"
  }
}
```

### Delete Activity
Deletes an activity.

- **URL**: `/activities/:id`
- **Method**: `DELETE`
- **Auth Required**: Yes (JWT + API key)
- **Response**:
```json
{
  "success": true,
  "message": "Activity deleted successfully",
  "data": null
}
```

### Search Activities
Searches for activities with various filters and pagination.

- **URL**: `/activities/search`
- **Method**: `GET`
- **Auth Required**: Yes (JWT + API key)
- **Query Parameters**:
  - `type`: Activity type (e.g., feeding, sleep, diaper)
  - `start_date`: Start date for filtering (ISO 8601 format)
  - `end_date`: End date for filtering (ISO 8601 format)
  - `details`: JSON string with details to filter by
  - `page`: Page number (default: 1)
  - `page_size`: Number of items per page (default: 10, max: 100)
- **Response**:
```json
{
  "success": true,
  "message": "Activities retrieved successfully",
  "data": [
    {
      "id": 1,
      "user_id": "user-id",
      "child_id": 1,
      "type": "feeding",
      "details": {
        "amount": 150,
        "unit": "ml",
        "notes": "Formula milk"
      },
      "happens_at": "2025-03-28T08:00:00Z",
      "created_at": "2025-03-28T07:43:04Z",
      "updated_at": "2025-03-28T08:00:00Z"
    },
    {
      "id": 2,
      "user_id": "user-id",
      "child_id": 1,
      "type": "diaper",
      "details": {
        "type": "wet",
        "notes": "Normal"
      },
      "happens_at": "2025-03-28T09:15:00Z",
      "created_at": "2025-03-28T09:15:00Z",
      "updated_at": "2025-03-28T09:15:00Z"
    }
  ],
  "pagination": {
    "total": 25,
    "current_page": 1,
    "page_size": 10,
    "total_pages": 3
  }
}
```

---

## Children

### Create Child
Creates a new child record.

- **URL**: `/api/children`
- **Method**: `POST`
- **Auth Required**: Yes (API key)
- **Request Body**:
```json
{
  "name": "Baby Smith",
  "details": {
    "birth_date": "2024-12-25",
    "gender": "male",
    "weight": 3.5,
    "height": 50
  }
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Child created successfully",
  "data": {
    "id": 1,
    "user_id": "user-id",
    "name": "Baby Smith",
    "details": {
      "birth_date": "2024-12-25",
      "gender": "male",
      "weight": 3.5,
      "height": 50
    },
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Get Child
Retrieves a specific child by ID.

- **URL**: `/api/children/:id`
- **Method**: `GET`
- **Auth Required**: Yes (API key)
- **Response**:
```json
{
  "success": true,
  "message": "Child retrieved successfully",
  "data": {
    "id": 1,
    "user_id": "user-id",
    "name": "Baby Smith",
    "details": {
      "birth_date": "2024-12-25",
      "gender": "male",
      "weight": 3.5,
      "height": 50
    },
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T07:43:04Z"
  }
}
```

### Get Children
Retrieves all children for the authenticated user with pagination.

- **URL**: `/api/children`
- **Method**: `GET`
- **Auth Required**: Yes (API key)
- **Query Parameters**:
  - `page`: Page number (default: 1)
  - `page_size`: Number of items per page (default: 10, max: 100)
- **Response**:
```json
{
  "success": true,
  "message": "Children retrieved successfully",
  "data": [
    {
      "id": 1,
      "user_id": "user-id",
      "name": "Baby Smith",
      "details": {
        "birth_date": "2024-12-25",
        "gender": "male",
        "weight": 3.5,
        "height": 50
      },
      "created_at": "2025-03-28T07:43:04Z",
      "updated_at": "2025-03-28T07:43:04Z"
    },
    {
      "id": 2,
      "user_id": "user-id",
      "name": "Baby Jane",
      "details": {
        "birth_date": "2025-01-15",
        "gender": "female",
        "weight": 3.2,
        "height": 48
      },
      "created_at": "2025-03-28T08:00:00Z",
      "updated_at": "2025-03-28T08:00:00Z"
    }
  ],
  "pagination": {
    "total": 2,
    "current_page": 1,
    "page_size": 10,
    "total_pages": 1
  }
}
```

### Update Child
Updates an existing child record.

- **URL**: `/api/children/:id`
- **Method**: `PUT`
- **Auth Required**: Yes (API key)
- **Request Body**:
```json
{
  "name": "Baby Smith Jr.",
  "details": {
    "birth_date": "2024-12-25",
    "gender": "male",
    "weight": 4.0,
    "height": 52
  }
}
```
- **Response**:
```json
{
  "success": true,
  "message": "Child updated successfully",
  "data": {
    "id": 1,
    "user_id": "user-id",
    "name": "Baby Smith Jr.",
    "details": {
      "birth_date": "2024-12-25",
      "gender": "male",
      "weight": 4.0,
      "height": 52
    },
    "created_at": "2025-03-28T07:43:04Z",
    "updated_at": "2025-03-28T09:00:00Z"
  }
}
```

## Postman Collection Setup

To use this API with Postman:

1. Create a new Postman Collection named "DailyAlu API"
2. Set up environment variables:
   - `base_url`: Your server URL (e.g., http://localhost:8080)
   - `access_token`: JWT token received after login
   - `refresh_token`: Refresh token received after login
   - `api_key`: API key for accessing the API (default test key: dk_test_12345)

3. Create request folders:
   - Authentication
   - Users
   - Activities
   - Children

4. Add the requests to their respective folders
5. For authenticated requests, add this to the Authorization tab:
   - Type: Bearer Token
   - Token: {{access_token}}

6. For all requests, add this header:
   - X-API-Key: {{api_key}}

7. After login, use the "Tests" tab in Postman to automatically set the tokens:
```javascript
var jsonData = pm.response.json();
if (jsonData.success && jsonData.data) {
    pm.environment.set("access_token", jsonData.data.access_token);
    pm.environment.set("refresh_token", jsonData.data.refresh_token);
}
```
