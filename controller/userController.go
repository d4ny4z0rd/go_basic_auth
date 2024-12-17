package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"new_user_auth_prac/model"

	"github.com/gorilla/mux"
)


func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	
	vars := mux.Vars(r)
	userId := vars["userId"]
	fmt.Println(userId)
	if userId == "" {
		http.Error(w, "UserId is required", http.StatusBadRequest)
		return
	}

	query := `SELECT userid, firstname, lastname, email from users WHERE userid=$1`
	row := db.QueryRow(query, userId)

	var user model.User
	err := row.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email)
	if err!=nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error retreiving user data", http.StatusInternalServerError)
		}
		return
	}

	response, err := json.Marshal(user)
	if err!=nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	query := `SELECT userid, firstname, lastname, email from users`
	
	rows, err := db.Query(query)
	if err!=nil {
		http.Error(w, "Error fetching users from the database", http.StatusInternalServerError)
		log.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	var users []model.User

	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email)
		if err!=nil {
			http.Error(w, "Error scanning user data", http.StatusInternalServerError)
			log.Println("Error scanning row:", err)
			return
		}
		users = append(users, user)
	}

	if err:= rows.Err(); err!=nil {
		http.Error(w, "Error iterating over the rows", http.StatusInternalServerError)
		log.Println("Row iteration error:", err)
		return
	}

	response, err := json.Marshal(users)
	if err!=nil {
		http.Error(w, "Error encoding response to JSON", http.StatusInternalServerError)
		log.Println("Error marshaling users:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	vars := mux.Vars(r)
	userId := vars["userId"]
	if userId == "" {
		http.Error(w, "UserId is required", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err!=nil {
		http.Error(w, "Failed to read the request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var updates map[string]interface{}
	err = json.Unmarshal(body, &updates)
	if err!=nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query := "UPDATE users SET "
	params := []interface{}{}
	paramIndex := 1
	for key, value := range updates {
		query += fmt.Sprintf("%s = $%d, ", key, paramIndex)
		params = append(params, value)
		paramIndex++
	}

	query = query[:len(query)-2]
	query += " WHERE userid = $" + fmt.Sprintf("%d", paramIndex)
	params = append(params, userId)

	_, err = db.Exec(query, params...)
	if err!=nil {
		http.Error(w, "Failed to update the user", http.StatusInternalServerError)
		log.Println("Error updating the user:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"User updated successfully"}`))
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")

	vars := mux.Vars(r)
	userId := vars["userId"]

	deleteQuery := "DELETE FROM users WHERE userid = $1"

	_, err := db.Exec(deleteQuery, userId)
	if err!=nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message" : "user deleted successfully"}`))
}
