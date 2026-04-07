package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/SWernerH/test_1/db"
	"github.com/SWernerH/test_1/models"
)

// GET /users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Email)
		users = append(users, u)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// POST /users
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var u models.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
		return
	}

	_, err = db.DB.Exec(
		"INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, NOW())",
		u.Username,
		u.Email,
		"test123",
	)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}