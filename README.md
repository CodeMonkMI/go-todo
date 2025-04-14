# Todo Application with MongoDB

## Overview
This is a simple Todo application built with Go and MongoDB. The application provides a RESTful API for managing todo items.

## Prerequisites
- Docker and Docker Compose
- Go (for development)

## Getting Started

### Running MongoDB with Docker Compose

1. Start the MongoDB container:
   ```bash
   docker-compose up -d
   ```

2. This will start MongoDB on `localhost:27017` with the database `demo_todo` initialized.

3. To stop the MongoDB container:
   ```bash
   docker-compose down
   ```

### Running the Application

1. Build and run the Go application:
   ```bash
   go build
   ./todo
   ```

2. The application will be available at `http://localhost:4000`

## API Endpoints

- `GET /` - Check if server is running
- `GET /todo` - Fetch all todos
- `POST /todo` - Create a new todo
- `PUT /todo/{id}` - Update a todo by ID

## Docker Compose Configuration

The `docker-compose.yml` file includes:
- MongoDB latest version
- Persistent data storage using Docker volumes
- Network configuration for service communication
- Port mapping to access MongoDB from the host machine