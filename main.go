// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description "JWT Authorization header using the Bearer scheme. Example: 'Bearer {token}'"

// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"

	"github.com/Anwarjondev/todo-api-go/db"
	"github.com/Anwarjondev/todo-api-go/routes"
	_ "github.com/Anwarjondev/todo-api-go/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)
// @title Todo List API with Authentication
// @version 1.0
// @description This is a simple API for managing todo lists with authentication.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main(){

	
	db.InitDB()

	mux := http.NewServeMux()
	routes.SetupRoutes(mux)
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	log.Println("Server is running port :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}