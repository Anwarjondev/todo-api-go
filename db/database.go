package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)


var DB *sql.DB


func InitDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSsLMode := os.Getenv("DB_SSLMODE")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, dbSsLMode)
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect db", err)
	}

	createTableUsers := `CREATE TABLE IF NOT EXISTS users(
		id SERIAL primary key,
		username text unique not null, 
		password text not null,
		role text not null check(role in('admin', 'user'))
	);`
	_, err = DB.Exec(createTableUsers)
	if err != nil {
		log.Fatal("Failde to create users table", err)
	}
	createTableQuery := `CREATE TABLE IF NOT EXISTS todos(
		id SERIAL PRIMARY KEY,
		title text not null,
		completed BOOLEAN not null default false,
		user_id smallint not null,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`
	_, err = DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal("Failed to create table", err)
	}
	
	
	
}
