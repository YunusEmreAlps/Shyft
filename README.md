# Shyft

This service that allows users to manage their shift periods and create their shift lists. This service allows users to select which days and times they will be on duty in a given date range and create their shifts. In addition, users can view the shift roster for their shift period, and makeshift change requests, and have the option to delete or update existing shifts. Through notifications, users can be notified of updates related to their shift periods. The Shift Service provides a user-friendly platform that helps users to schedule shifts in an organized manner and manage their roster.

## Table of Contents

- [Shyft](#shyft)
  - [Table of Contents](#table-of-contents)
  - [Tech Stack](#tech-stack)
  - [Installation](#installation)
    - [Requirements](#requirements)
    - [Quick Start](#quick-start)
  - [Project Structure](#project-structure)
  - [Swagger Documentation](#swagger-documentation)
  - [Contact](#contact)

## Tech Stack

- **Backend**: Go, Gin Framework
- **Database**: PostgreSQL, GORM
- **Cache**: Redis
- **Monitoring**: Grafana, Prometheus, Jaeger
- **Docs**: Swagger
- **Container**: Docker

## Installation

### Requirements

- Docker & Docker Compose
- Go 1.21+ (for local development)

### Quick Start

```bash
# Move to your workspace
cd your-workspace

# Clone this project into your workspace
git clone https://github.com/YunusEmreAlps/shyft.git

# Move to the project root directory
cd shyft

# Edit .env with your configuration
cp .env.example .env

# Run docker containers
docker-compose up -d

# To start only the Go application (in a separate terminal):
go run .
```

```bash
- http://localhost:9097 (API Server)
- http://localhost:9097/shyft/swagger/index.html (Swagger UI)
- http://localhost:16686 (Jaeger UI)
- http://localhost:9090 (Prometheus)
- http://localhost:3000 (Grafana)
```

> Grafana: Default username and password: **admin / admin**

## Project Structure

```bash
shyft/
|-- config/         # Configuration files
|-- docs/           # Swagger documentation files
|-- internal/       # Private application code
|   |-- handlers/   # HTTP handlers
|   |-- models/     # Database models
|   |-- repository/ # Database repository layer
|-- pkg/            # Public application code
|-- Dockerfile    # Dockerfile for building the application image
|-- go.mod
|-- go.sum
|-- main.go
|-- README.md
```

## Swagger Documentation

The Swagger documentation for the shift scheduler project can be accessed by following these steps:

1. After running the project, open a web browser and navigate to `http://localhost:9097/shyft/swagger/index.html` to access the Swagger documentation
2. When you add a new API endpoint to the project, you need to run the `swag init` command to generate the Swagger documentation for the new endpoint. This command will update the `docs` directory with the new Swagger documentation.

## Contact

For any questions or support, please contact us at:

- Linkedin at [Yunus Emre Alpu](https://www.linkedin.com/in/yunus-emre-alpu-5b1496151/)
