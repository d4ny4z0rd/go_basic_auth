package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"new_user_auth_prac/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func InitiateDBConnection(database *sql.DB) {
	db = database
}

var jwtSecretKey = os.Getenv("JWT_SECRET") 
var jwtSecretByte  = []byte(jwtSecretKey)

type Claims struct {
	UserId	string	`json:"userId"`
	jwt.StandardClaims	
}

func Login (w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	
	var user model.User

	//Reading the user body
	body, err := io.ReadAll(r.Body)
	if err!=nil {
		http.Error(w, "Failed to read the request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err!=nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if user.Password == "" || user.Email == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	//Querying the database based on email
	var storedUser model.User
	query := `SELECT userid, firstname, lastname, email, password FROM users WHERE email=$1`
	err = db.QueryRow(query, user.Email).Scan(&storedUser.UserId, &storedUser.FirstName, &storedUser.LastName, &storedUser.Email, &storedUser.Password)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if err!=nil {
		http.Error(w, "Error fetching user from the database", http.StatusInternalServerError)
		return
	}

	//Compare hashed password from db with provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password))
	if err!=nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	//Create JWT token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserId : storedUser.UserId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer: "auth_app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Sign the token with the secret key
	fmt.Println(jwtSecretByte)
	tokenString, err := token.SignedString(jwtSecretByte)
	if err!=nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"token": tokenString}
	responseJSON, err := json.Marshal(response)
	if err!=nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func Signup (w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello from the Signup route")

	w.Header().Set("Content-Type","application/json")

	var user model.User

	body, err := io.ReadAll(r.Body)
	if err!=nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &user)
	if err!=nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err!=nil {
		http.Error(w, "Error hashing the password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	user.UserId = uuid.NewString()

	insertUserQuery := `INSERT INTO USERS (firstname, lastname, password, email, userid) VALUES ($1, $2, $3, $4, $5)`

	_, err = db.Exec(insertUserQuery, user.FirstName, user.LastName, user.Password, user.Email, user.UserId)	
	if err!=nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	user.Password = ""

	response, err := json.Marshal(user)
	if err!=nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}
