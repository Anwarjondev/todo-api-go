package routes

import (
	"fmt"
	"net/http"

	"github.com/Anwarjondev/todo-api-go/handlers"
	"github.com/Anwarjondev/todo-api-go/middleware"
)

func SetupRoutes(mux *http.ServeMux) {
	// Default route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "🚀 Welcome to the Todo API!\n\nUse Postman or Swagger to access endpoints.")
	})

	// Public routes
	mux.HandleFunc("POST /register", handlers.Register)
	mux.HandleFunc("POST /login", handlers.Login)

	// Protected user routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /todos", handlers.GetTodos)
	protectedMux.HandleFunc("POST /todos/create", handlers.CreateTodo)
	protectedMux.HandleFunc("PUT /todos/update", handlers.UpdateTodo)
	protectedMux.HandleFunc("DELETE /todos/delete", handlers.DeleteTodo)

	// Admin routes
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("DELETE /admin/todos", handlers.DeleteAllTodos)
	adminMux.HandleFunc("GET /admin/getallusers", handlers.GetAllUsers)

	// Apply middlewares
	mux.Handle("/todos", middleware.AuthMiddleware(protectedMux))
	mux.Handle("/todos/", middleware.AuthMiddleware(protectedMux))
	mux.Handle("/admin/", middleware.AuthMiddleware(middleware.AdminMiddleware(adminMux)))
}
