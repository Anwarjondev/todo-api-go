package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Anwarjondev/todo-api-go/db"
	"github.com/Anwarjondev/todo-api-go/models"
)

// GetTodos retrieves all todos or filters by status (completed/pending)
// @Summary Get Todos
// @Description Retrieve todos based on user role
// @Tags Todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 401 {string} string "Unauthorized"
// @Router /todos [get]
func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id")
	roleValue := r.Context().Value("role").(string)

	if userID == nil {
		http.Error(w, "Unauthorized: Missing user ID", http.StatusUnauthorized)
		return
	}
	id, ok := userID.(int)
	if !ok {
		http.Error(w, "Unauthorized: Invalid user ID", http.StatusUnauthorized)
		return
	}
	var query string
	var rows *sql.Rows
	var err error
	completedParam := r.URL.Query().Get("completed")
	if roleValue == "admin" {
		query = "SELECT * FROM todos"
	} else {
		query = "select * from todos where user_id = $1"
	}
	if completedParam == "true" {
		query += "and completed = true"
	} else if completedParam == "false" {
		query += "and completed = false"
		rows, err = db.DB.Query(query, id)
	}
	if roleValue == "admin" {
		rows, err = db.DB.Query(query)
	} else {
		rows, err = db.DB.Query(query, userID)
	}
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.UserId); err != nil {
			http.Error(w, "Error scanning row:", http.StatusInternalServerError)
			return
		}
		todos = append(todos, todo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// CreateTodo creates a new todo
// @Summary Create Todo
// @Description Create a new todo (only authenticated users)
// @Tags Todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param todo body models.TodoModel true "Todo data"
// @Success 201 {object} models.TodoModel
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /todos/create [post]
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var todo models.TodoModel
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if todo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	stmt, err := db.DB.Prepare("insert into todos(title, user_id) values($1, $2)")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(todo.Title, userID)
	if err != nil {
		http.Error(w, "Error creating with todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodo updates an existing todo
// @Summary Update Todo
// @Description Update a todo (users can update only their own todos)
// @Tags Todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param todo body models.UpdateTodoModel true "Updated todo data"
// @Param id query int true "Todo ID"
// @Success 200 {object} models.UpdateTodoModel
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Server error"
// @Router /todos/update [put]
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "Unauthorized: User ID not found", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}

	var todo models.UpdateTodoModel

	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var exists bool
	err := db.DB.QueryRow("select exists(select 1 from todos where id = $1 and user_id = $2)", id, userID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Todo not found or you do not have access", http.StatusForbidden)
		return
	}

	stmt, err := db.DB.Prepare("Update todos set title = $1, completed = $2 where id = $3 and user_id = $4")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Databse error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(todo.Title, todo.Completed, id, userID)
	if err != nil {
		http.Error(w, "Error with updating todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes a todo
// @Summary Delete Todo
// @Description Delete a todo (users can delete only their own todos, admins can delete any)
// @Tags Todos
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id query int true "Todo ID"
// @Success 200 {string} string "Todo deleted"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 500 {string} string "Server error"
// @Router /todos/delete [delete]
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Id is required", http.StatusBadRequest)
		return
	}
	var exists bool
	err := db.DB.QueryRow("Select Exists(select 1 from todos where id = $1 and user_id = $2)", id, userID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Todo not found or you do not have access", http.StatusForbidden)
		return
	}
	stmt, err := db.DB.Prepare("Delete from todos where id = $1 and user_id = $2")
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, userID)
	if err != nil {
		http.Error(w, "Error with deleting todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// DeleteAllTodos godoc
// @Summary Delete Any Todo (Admin Only)
// @Description Admin can delete any todo
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id query int true "Todo ID"
// @Success 200 {string} string "Todo deleted"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden (Admins only)"
// @Failure 500 {string} string "Server error"
// @Router /admin/todos [delete]
func DeleteAllTodos(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)
	id := r.URL.Query().Get("id")
	if role != "admin" {
		http.Error(w, "Forbidden: only Admin can delete all todos", http.StatusForbidden)
		return
	}

	_, err := db.DB.Exec("delete from todos where id = $1", id)
	if err != nil {
		http.Error(w, "Failed to delete al todos", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "All todos successfully deleted"})
}

// Get all users
// @Summary Get all users (Admin Only)
// @Description Admin can get all users
// @Tags Admin
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {string} string "Get Users"
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden (Admins only)"
// @Failure 500 {string} string "Server error"
// @Router /admin/getallusers [get]
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role").(string)
	if role != "admin" {
		http.Error(w, "Only admin can see all users", http.StatusForbidden)
		return
	}
	var users []models.AllUser
	rows, err := db.DB.Query("select id, username, role from users")
	if err != nil {
		http.Error(w, "Error with fetching all users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var username string
		var role string
		err := rows.Scan(&id, &username, &role)
		if err != nil {
			http.Error(w, "Error with fetching users", http.StatusInternalServerError)
			return
		}
		user := models.AllUser{
			ID:       id,
			Username: username,
			Role:     role,
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Error reading rows", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}
