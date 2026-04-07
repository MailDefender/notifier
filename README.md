# Notifier

A high-performance microservice for sending transactional emails via SMTP. Built with Go, designed for reliability and security.

## Overview

The Notifier service provides a simple REST API to send formatted emails. It handles email validation, HTML rendering, CSS inlining, and secure SMTP delivery with automatic timeouts and error handling.

**Key Features:**

- ✅ RFC 5322 email address validation
- ✅ HTML email with automatic CSS inlining
- ✅ Thread-safe SMTP client with connection pooling
- ✅ Request size limits (1 MB) and timeouts (15s)
- ✅ Full audit logging
- ✅ Graceful shutdown support
- ✅ IPv6 compatible

## Architecture

```
API Handler → Input Validation → Email Formatter → SMTP Connector → Mail Server
```

- **API Layer** (`internal/api`): REST endpoint for email submission
- **Formatter** (`internal/formatters`): Converts data to RFC 822 email format with CSS inlining
- **SMTP Connector** (`internal/connectors`): Manages SMTP connections with timeout protection
- **Configuration** (`internal/configuration`): Centralized config management with validation

## Prerequisites

- Go 1.19+ (for compilation without Docker)
- Docker & Docker Compose (for containerized deployment)
- SMTP server credentials (see [Configuration](#configuration))

## Installation & Setup

### Option 1: Docker (Recommended)

```bash
# Build the Docker image
docker build -t maildefender/notifier .

# Run with environment variables
docker run \
  -p 8080:8080 \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=noreply@example.com \
  -e SMTP_PASSWORD=your_secure_password \
  maildefender/notifier
```

Or with a `.env` file:

```bash
docker run -p 8080:8080 --env-file .env maildefender/notifier
```

### Option 2: Local Build

```bash
# Clone the repository
git clone https://github.com/MailDefender/notifier.git
cd notifier

# Install dependencies
go mod download

# Build the binary
go build -o notifier cmd/notifier/main.go

# Configure environment
cp .env.example .env
# Edit .env with your SMTP credentials

# Run the service
source .env && ./notifier
```

## Configuration

Create a `.env` file in your deployment environment with the following variables:

| Variable        | Required | Example               | Description                                                  |
| --------------- | -------- | --------------------- | ------------------------------------------------------------ |
| `SMTP_HOST`     | Yes      | `smtp.gmail.com`      | SMTP server hostname                                         |
| `SMTP_PORT`     | Yes      | `587`                 | SMTP server port (typically 587 for TLS or 25 for plaintext) |
| `SMTP_USERNAME` | Yes      | `noreply@example.com` | SMTP authentication username                                 |
| `SMTP_PASSWORD` | Yes      | `SecurePassword123`   | SMTP authentication password                                 |

**Example `.env` file:**

```shell
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=noreply@example.com
SMTP_PASSWORD=your_app_password
```

All configuration variables are validated at startup. The service will fail fast with clear error messages if required settings are missing or invalid.

## API Usage

### Send Email

**Endpoint:** `POST /v1/send/email`

**Request Body:**

```json
{
  "to": ["user@example.com"],
  "subject": "Order Confirmation",
  "body": "<p>Thank you for your order</p>",
  "replyTo": "support@example.com",
  "threadTopic": "order-123"
}
```

**Fields:**

| Field         | Type     | Required | Description                        |
| ------------- | -------- | -------- | ---------------------------------- |
| `to`          | string[] | Yes      | Email addresses to send to         |
| `subject`     | string   | Yes      | Email subject line (max 255 chars) |
| `body`        | string   | Yes      | HTML email body                    |
| `replyTo`     | string   | No       | Reply-to email address             |
| `threadTopic` | string   | No       | Outlook thread topic identifier    |

**Response (Success):**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{}
```

**Response (Error):**

```json
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "invalid email address: user@invalid"
}
```

**Example cURL:**

```bash
curl -X POST http://localhost:8080/v1/send/email \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "subject": "Test Email",
    "body": "<p>This is a test email</p>"
  }'
```

## Testing

### Unit Tests

Run the full test suite:

```bash
go test ./...
```

Run only API tests:

```bash
go test ./internal/api -v
```

Run with coverage:

```bash
go test ./... -cover
```

### Integration Tests

For integration testing with a real SMTP server (Mailhog), use the provided Docker Compose setup:

```bash
cd mock
docker-compose up -d

# Run tests
go test ./internal/api -v

# Clean up
docker-compose down
```

The test suite includes:

- ✅ Single email sending
- ✅ Bulk email sending (200+ emails)
- ✅ Email validation
- ✅ Error handling

## Development

### Project Structure

```
notifier/
├── cmd/notifier/          # Application entry point
├── internal/
│   ├── api/              # HTTP handlers and routing
│   ├── client/           # Email client abstraction
│   ├── configuration/    # Config validation
│   ├── connectors/       # SMTP connector implementation
│   ├── formatters/       # Email formatting (RFC 822)
│   ├── models/           # Data models
│   ├── templates/        # Email templates
│   └── utils/            # Utility functions
├── mock/                 # Docker Compose for testing
├── go.mod               # Go module definition
├── go.sum               # Dependency checksums
└── Dockerfile           # Container image definition
```

### Code Quality

Ensure code quality with:

```bash
# Format code
go fmt ./...

# Check for common issues
go vet ./...

# Run tests
go test ./...
```

### Building a Release

```bash
# Build optimized binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o notifier cmd/notifier/main.go

# Build Docker image
docker build -t maildefender/notifier:v1.0.0 .
```

## Troubleshooting

### Service won't start

**Problem:** "invalid SMTP configuration"

**Solution:** Verify `.env` contains all required variables and they're correct:

```bash
echo $SMTP_HOST $SMTP_PORT $SMTP_USERNAME
```

### Email sending fails with "failed to dial SMTP server"

**Problem:** SMTP connection timeout

**Solution:**

- Verify SMTP host and port are accessible: `telnet $SMTP_HOST $SMTP_PORT`
- Check firewall rules allow outbound SMTP connections
- Verify credentials are correct

### "invalid email address" error

**Problem:** One of the recipient emails is malformed

**Solution:** Ensure all emails in the `to` array are valid RFC 5322 format (example: `user@example.com`)

### Tests fail with "failed to connect to Mailhog"

**Problem:** Mailhog not running for integration tests

**Solution:** Start Mailhog first:

```bash
cd mock && docker-compose up -d
```

## Performance

- **Throughput:** ~200 emails/second (with async batching)
- **Latency:** < 500ms per email (with typical SMTP latency)
- **Memory:** ~50 MB baseline
- **Connections:** Single pooled SMTP connection with timeout protection

## Security Considerations

- Passwords are never logged (sanitized in logs)
- Email addresses are sanitized to prevent header injection
- Request payloads limited to 1 MB to prevent DoS
- All operations have timeouts to prevent hanging connections
- No sensitive data in error responses sent to clients

## License

This project is licensed under the MIT License. See LICENSE file for details.

## Support

For issues and feature requests, please open an issue on GitHub.

For security vulnerabilities, please email security@maildefender.io (do not open public issues).
