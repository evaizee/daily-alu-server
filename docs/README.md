# DailyAlu API Documentation

This directory contains comprehensive documentation for the DailyAlu API, including a detailed API reference and a ready-to-use Postman collection.

## Contents

- [API Documentation](api_documentation.md) - Complete reference for all API endpoints, request/response formats, and authentication requirements
- [Postman Collection](dailyalu_postman_collection.json) - Ready-to-import Postman collection for testing the API

## Getting Started with the Postman Collection

1. **Import the Collection**:
   - Open Postman
   - Click "Import" in the top left
   - Select the `dailyalu_postman_collection.json` file

2. **Set Up Environment Variables**:
   - Create a new environment in Postman
   - Add the following variables:
     - `base_url`: Your server URL (e.g., http://localhost:8080)
     - `api_key`: Default test key is "dk_test_12345"
     - `access_token`: Will be automatically set after login
     - `refresh_token`: Will be automatically set after login

3. **Test the API**:
   - Start with the "Register" and "Login" endpoints in the Authentication folder
   - After successful login, the access and refresh tokens will be automatically set
   - Explore other endpoints using the authenticated requests

## API Structure

The API is organized into the following categories:

1. **Authentication** - User registration, login, and token management
2. **Users** - User profile management
3. **Activities** - Baby activity tracking (feeding, sleep, diaper changes, etc.)
4. **Children** - Child profile management

## Authentication

Most endpoints require authentication using:

1. **API Key** - Include in all requests via the `X-API-Key` header
2. **JWT Token** - Include in authenticated requests via the `Authorization: Bearer <token>` header

The Postman collection is configured to automatically handle authentication after login.

## Response Format

All API responses follow a standardized format:

- Success responses include `success: true` with a message and data
- Error responses include `success: false` with error details
- Paginated responses include pagination metadata

## Development Notes

- The API uses Go 1.22 with Fiber v2.52.0
- Responses follow the standardized response package format
- Input validation is performed using validator v10
