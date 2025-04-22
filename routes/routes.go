package routes

import (
	"net/http"

	"github.com/Anwarjondev/todo-api-go/handlers"
	"github.com/Anwarjondev/todo-api-go/middleware"
)

func SetupRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /register", handlers.Register)
	mux.HandleFunc("POST /login", handlers.Login)

	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /todos", handlers.GetTodos)
	protectedMux.HandleFunc("POST /todos/create", handlers.CreateTodo)
	protectedMux.HandleFunc("PUT /todos/update", handlers.UpdateTodo)
	protectedMux.HandleFunc("DELETE /todos/delete", handlers.DeleteTodo)

	adminmux := http.NewServeMux()
	adminmux.HandleFunc("DELETE /admin/todos", handlers.DeleteAllTodos)
	adminmux.HandleFunc("GET /admin/getallusers", handlers.GetAllUsers)

	mux.Handle("/", middleware.AuthMiddleware(protectedMux))
	mux.Handle("/admin/", middleware.AuthMiddleware(middleware.AdminMiddleware(adminmux)))
	
}
