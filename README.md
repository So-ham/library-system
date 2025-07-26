# Library Management System API

A RESTful API for a fictional public library built with Go, following clean architecture principles.

## Features

- List all books in the library
- Create, Read, Update, and Delete operations for books
- PostgreSQL database for data persistence
- Docker and Docker Compose setup for easy deployment

## Project Structure

The project follows clean architecture principles with the following layers:

- **Entities**: Core business objects (Book)
- **Models**: Data access layer
- **Services**: Business logic layer
- **Handlers**: HTTP request/response handling
- **Web/Rest**: Routing and API endpoints

## API Endpoints

- `GET /api/books` - List all books
- `GET /api/books/{id}` - Get a specific book
- `POST /api/books` - Create a new book
- `PUT /api/books/{id}` - Update a book
- `DELETE /api/books/{id}` - Delete a book

## Running the Application

### Prerequisites

- Docker and Docker Compose installed

### Using Docker Compose

1. Clone the repository
2. Navigate to the project directory
3. Run the application:

```bash
docker-compose up
```

The API will be available at http://localhost:8080

### Without Docker

1. Make sure you have PostgreSQL installed and running
2. Update the `.env` file with your database connection string
3. Run the application:

```bash
go run cmd/server.go
```

## Example API Usage

### Create a Book

```bash
curl -X POST http://localhost:8080/api/books \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "isbn": "9780743273565",
    "publisher": "Scribner",
    "publish_date": "2004-09-30T00:00:00Z",
    "description": "A classic novel about the American Dream",
    "copies": 5
  }'
```

### Get All Books

```bash
curl -X GET http://localhost:8080/api/books
```

### Get a Book by ID

```bash
curl -X GET http://localhost:8080/api/books/{id}
```

### Update a Book

```bash
curl -X PUT http://localhost:8080/api/books/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "isbn": "9780743273565",
    "publisher": "Scribner",
    "publish_date": "2004-09-30T00:00:00Z",
    "description": "Updated description",
    "copies": 10
  }'
```

### Delete a Book

```bash
curl -X DELETE http://localhost:8080/api/books/{id}
```
