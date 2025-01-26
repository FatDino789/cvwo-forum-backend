# Forum Backend Setup Guide

## Prerequisites

- **Go** (v1.20 or higher)
- **PostgreSQL** (v14.0 or higher)
- **Docker**

---

## Installation

### Clone Repository

```bash
git clone https://github.com/FatDino789/cvwo-forum-backend.git
cd forum-backend
```

### Load Environment Variables

Load the .env file with the below mentioned fields:

```
DB_HOST=your_host
DB_PORT=your_port
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
JWT_SECRET=your_jwt_secret
```

### Install Dependencies

```bash
# Install Go modules
go mod tidy
```

### Running the Application with Docker

#### Prerequisites

- Ensure Docker is installed on your machine. [Install Docker](https://docs.docker.com/get-docker/).
- Docker Compose is included with most Docker installations. Confirm by running `docker-compose --version`.

#### Docker Setup

The project includes a `docker-compose.yml` file that simplifies setting up both the backend and PostgreSQL database. (to be placed at the root of your directory)

#### Steps to Run the Application with Docker

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/FatDino789/cvwo-forum-backend.git
   cd forum-backend
   ```

2. **Start Docker Containers**:

   Run the following command to build and start the containers:

   ```bash
   docker-compose up --build
   ```

   This command will:

   - Build the backend service from the `Dockerfile`.
   - Set up a PostgreSQL database with the configuration specified in `docker-compose.yml`.

3. **Verify the Setup**:

   - The backend server will be accessible at `http://localhost:8080`.
   - PostgreSQL will be running and can be accessed on `localhost:5432` (or the port specified in `docker-compose.yml`).

4. **Stopping the Containers**:

   To stop and remove the containers:

   ```bash
   docker-compose down
   ```

   This will stop all running containers and remove any associated resources (except for volumes).

#### Notes:

- The `postgres-data` volume persists PostgreSQL data, so your database will retain its state even after restarting the containers.
- Ensure the `.env` file is present before running the Docker setup.

---

### Restoring the Database

To restore the database from a dump:

```bash
psql -U your_username -h your_host -p your_port -d your_database < exchange-forum-dump.sql
```

---

## Technology Stack

- **Go (Golang)**: Main programming language.
- **PostgreSQL**: Relational database.
- **Chi Router**: Lightweight router for building REST APIs.
- **Docker**: Hosting the PostgreSQL Database.
- **JWT**: Authentication.

---

## Project Structure

```
forum-backend/
    ├── cmd/
    │   └── api/
    │       ├── config/            # Configuration files (database, environment)
    │       ├── handlers.go        # API route handlers
    │       ├── middleware.go      # Middleware functions (e.g., authentication)
    │       ├── routes.go          # Route definitions
    │       └── main.go            # Entry point for the backend server
    ├── sql/
    │   └── create_tables.sql      # Database schema and table definitions
    ├── postgres-data/             # Volume for PostgreSQL data (Docker)
    ├── .env                       # Environment variables
    ├── docker-compose.yml         # Docker configuration
    ├── go.mod                     # Go module dependencies
    ├── go.sum                     # Go module checksum
    └── README.md                  # Documentation
```

---
