# Address API

A REST API for managing addresses, built with Go and chi router. Designed for deployment on AWS Lambda + API Gateway.

## Features

- Full CRUD operations (Create, Read, List, Update, Delete)
- UUID-based address identification
- Timestamp tracking (created_at, updated_at)
- Comprehensive validation (all fields required except state)
- Zip/postal code format validation
- Thread-safe in-memory storage
- Synchronous responses
- AWS Lambda integration

## Schema

Address model includes the following fields:

- `id` (string): UUID, auto-generated
- `street` (string): Street address (required)
- `city` (string): City name (required)
- `state` (string): State/province (optional)
- `zip` (string): Zip/postal code (required, format: 3-10 alphanumeric characters with optional hyphens)
- `country` (string): Country name (required)
- `created_at` (timestamp): Creation timestamp (UTC)
- `updated_at` (timestamp): Last update timestamp (UTC)

## API Endpoints

### Health Check
```
GET /health
```

### Create Address
```
POST /addresses
Content-Type: application/json

{
  "street": "123 Main St",
  "city": "New York",
  "state": "NY",
  "zip": "10001",
  "country": "USA"
}

Response: 201 Created
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "street": "123 Main St",
  "city": "New York",
  "state": "NY",
  "zip": "10001",
  "country": "USA",
  "created_at": "2026-07-11T12:00:00Z",
  "updated_at": "2026-07-11T12:00:00Z"
}
```

### Get Address
```
GET /addresses/{id}

Response: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "street": "123 Main St",
  "city": "New York",
  "state": "NY",
  "zip": "10001",
  "country": "USA",
  "created_at": "2026-07-11T12:00:00Z",
  "updated_at": "2026-07-11T12:00:00Z"
}
```

### List All Addresses
```
GET /addresses

Response: 200 OK
{
  "addresses": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "street": "123 Main St",
      "city": "New York",
      "state": "NY",
      "zip": "10001",
      "country": "USA",
      "created_at": "2026-07-11T12:00:00Z",
      "updated_at": "2026-07-11T12:00:00Z"
    }
  ],
  "count": 1
}
```

### Update Address
```
PUT /addresses/{id}
Content-Type: application/json

{
  "street": "456 Oak Ave",
  "city": "Boston",
  "zip": "02101"
}

Response: 200 OK
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "street": "456 Oak Ave",
  "city": "Boston",
  "state": "NY",
  "zip": "02101",
  "country": "USA",
  "created_at": "2026-07-11T12:00:00Z",
  "updated_at": "2026-07-11T12:00:15Z"
}
```

### Delete Address
```
DELETE /addresses/{id}

Response: 204 No Content
```

## Error Response Format

All errors follow this format:

```json
{
  "error": "description of the error"
}
```

Common HTTP status codes:
- `201 Created`: Address successfully created
- `200 OK`: Successful operation (GET, PUT)
- `204 No Content`: Successful deletion
- `400 Bad Request`: Validation error or invalid request format
- `404 Not Found`: Address not found or operation failed

## Installation

### Prerequisites

- Go 1.21 or later
- AWS CLI (for Lambda deployment)

### Build

```bash
go build -o bootstrap .
```

The output binary is named `bootstrap` as required by AWS Lambda custom runtime.

## Testing

Run all tests:
```bash
go test -v
```

Run tests with coverage:
```bash
go test -cover
```

Tests include:
- Unit tests for Address validation
- Unit tests for in-memory store (CRUD operations)
- Integration tests for HTTP handlers
- Concurrent operation tests

## Local Development

Run the API locally:

```bash
go run .
```

The API will start on `http://localhost:8080`

Test the API:
```bash
curl -X POST http://localhost:8080/addresses \
  -H "Content-Type: application/json" \
  -d '{
    "street": "123 Main St",
    "city": "New York",
    "state": "NY",
    "zip": "10001",
    "country": "USA"
  }'
```

## Deployment to AWS Lambda + API Gateway

### Step 1: Build for Lambda

```bash
GOOS=linux GOARCH=arm64 go build -o bootstrap .
zip deployment.zip bootstrap
```

### Step 2: Create Lambda Function

```bash
aws lambda create-function \
  --function-name address-api \
  --role arn:aws:iam::YOUR_ACCOUNT_ID:role/lambda-execution-role \
  --zip-file fileb://deployment.zip \
  --runtime provided.al2 \
  --timeout 30 \
  --memory-size 256 \
  --handler bootstrap
```

Or update an existing function:

```bash
aws lambda update-function-code \
  --function-name address-api \
  --zip-file fileb://deployment.zip
```

### Step 3: Create API Gateway

Create a REST API with:
- `POST /addresses` → Lambda function
- `GET /addresses` → Lambda function
- `GET /addresses/{id}` → Lambda function
- `PUT /addresses/{id}` → Lambda function
- `DELETE /addresses/{id}` → Lambda function
- `GET /health` → Lambda function

### Step 4: Deploy API Gateway

```bash
aws apigateway create-deployment \
  --rest-api-id YOUR_API_ID \
  --stage-name prod
```

## Using LocalStack for Testing

LocalStack provides local AWS service mocking for testing:

```bash
# Start LocalStack
docker run -d -p 4566:4566 localstack/localstack

# Deploy to LocalStack
aws --endpoint-url=http://localhost:4566 lambda create-function \
  --function-name address-api \
  --runtime provided.al2 \
  --role arn:aws:iam::000000000000:role/lambda-role \
  --handler bootstrap \
  --zip-file fileb://deployment.zip
```

## Architecture

The API follows a layered architecture:

- **Handler Layer** (`handler.go`): HTTP request/response handling
- **Store Layer** (`store.go`): In-memory persistence with thread safety
- **Model Layer** (`address.go`): Address data model and validation
- **Main** (`main.go`): Router setup and Lambda integration

## Notes

- All timestamps are in UTC
- IDs are UUIDs (version 4)
- State field is optional for international addresses
- The in-memory store does not persist data across Lambda invocations
- For persistent storage, integrate with DynamoDB or another database service

## Dependencies

- `github.com/go-chi/chi/v5`: Lightweight HTTP router
- `github.com/google/uuid`: UUID generation
- `github.com/aws/aws-lambda-go`: Lambda runtime
- `github.com/awslabs/aws-lambda-go-api-proxy`: API Gateway proxy adapter