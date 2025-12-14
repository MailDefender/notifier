# **notifier**

**This microservice is used to send notifications, mainly via emails**

> Note: this documentation is currently being drafted and will be completed in a future version.

## 📦 Prerequisites

- **Golang**
- **Docker** (optional)
- **[Environment Variables](#-configuration)**

## 🚀 Installation

### With Docker (Recommended)

```bash
docker build -t maildefender/notifier .
docker run -p 8080:8080 --env-file .env maildefender/notifier
```

### Without Docker

1. Clone the repository

```bash
# Clone the repoisitory
git clone https://github.com/MailDefender/notifier.git
cd notifier

# Install dependencies
go mod download

# Build
go build -o notifier

# Run
source .env
./notifier
```

## 🏃‍♂️ Usage

This app exposes APIs, so please refer to the Swagger to get more details about its usage.

## 🧪 Tests

Tests will be added soon.

## 🛠 Configuration

Create a .env file in the project root with the following variables:

```shell
SMTP_HOST=example.com
SMTP_PORT=587
SMTP_USERNAME=hello@me.com
SMTP_PASSWORD=P4ssw0rd
```

## 📜 License

This project is licensed under MIT.
