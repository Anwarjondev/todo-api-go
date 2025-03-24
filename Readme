# Todo List API

This is a simple Todo List API built using **Go** and the **net/http** package. The API supports user authentication, CRUD operations for todos, and API documentation using Swagger.

## Features

- User authentication (JWT-based)

- CRUD operations for todos

- Role-based access control (admin/user)

- PostgreSQL database integration

- API documentation using Swagger

## Project Structure

```
todo-api/
│── main.go
│── handlers/
│   ├── auth.go
│   ├── todos.go
│── middleware/
│   ├── auth.go
│── models/
│   ├── user.go
│   ├── todo.go
│── routes/
│   ├── routes.go
│── database/
│   ├── db.go
│── docs/
│   ├── swagger.json
│── go.mod
│── README.md
```

## Setup & Installation

### Prerequisites

- Install [Go](https://go.dev/doc/install)

- Install [PostgreSQL](https://www.postgresql.org/download/)

### Clone the Repository
```sh
git clone https://github.com/Anwarjondev/todo-api.git
cd todo-api
```

### Install Dependencies
```sh
go mod tidy
```

### Set Up Environment Variables

Create a `.env` file (optional but recommended):
```sh
echo "DB_HOST=localhost" > .env
echo "DB_PORT=5432" >> .env
echo "DB_USER=your_user" >> .env
echo "DB_PASSWORD=your_password" >> .env
echo "DB_NAME=your_database" >> .env
echo "DB_SSLMODE=disable" >> .env
echo "JWT_SECRET=your_secret_key" >> .env
```

Or set it manually:
```sh
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_user
export DB_PASSWORD=your_password
export DB_NAME=your_database
export DB_SSLMODE=disable
export JWT_SECRET=your_secret_key
```

### Run the Server
```sh
go run main.go
```

The server should be running at **`http://localhost:8080`**.

## API Endpoints
| Method    | Endpoint       | Description                 |
|-----------|----------------|-----------------------------|
| `POST`    | `/register`    | Register New User           |
| `POST`    | `/login`       | Authenticate user           |
| `GET`     | `/todos`       | List all todos              |
| `POST`    | `/create`      | Create new todo             |
| `PUT`     | `/update`      | Edit todo                   |
| `DELETE`  | `/delete`      | Delete todo only the own    |
| `DELETE`  | `/admin/todos` | Delete all any user todos   |

## Database Schema

### Users Table
```sh
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
```

### Todos Table
```sh
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title TEXT NOT NULL,
    completed BOOLEAN DEFAULT FALSE
);
```

### API Documentation

Swagger documentation is available at:
```
http://localhost:8080/swagger/index.html
```
### License

MIT License

### Author

Developed by **Anvarjon**

