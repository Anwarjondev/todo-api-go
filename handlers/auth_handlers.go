package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Anwarjondev/todo-api-go/db"
	"github.com/Anwarjondev/todo-api-go/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var jwtkey []byte

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: .env file not found, falling back to OS env variables")
	}

	
	jwtkey = []byte(os.Getenv("JWT_KEY"))
}

type Claims struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register a new user
// @Summary Register User
// @Description Register a new user (default role: user)
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.UserModel true "User Registration Data"
// @Success 201 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Server error"
// @Router /register [post]
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if user.Role == "" {
		user.Role = "user"
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error with hashing password", http.StatusInternalServerError)
		return
	}
	_, err = db.DB.Exec("Insert into users(username, password, role) values($1, $2, $3)", user.Username, string(hashedPassword), user.Role)
	if err != nil {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User create registered successfuly"})
}

// Login and generate JWT
// @Summary User Login
// @Description Login user and receive JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param user body models.UserModel true "User Credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Router /login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	var storedPassword string
	var userId int
	var userRole string
	err = db.DB.QueryRow("Select id, password, role from users where username = $1", user.Username).Scan(&userId, &storedPassword, &userRole)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	expritionTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		UserId: userId,
		Role:   userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expritionTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
