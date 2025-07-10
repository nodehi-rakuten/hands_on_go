package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID                      int            `json:"id"`
	FirstName               string         `json:"first_name"`
	LastName                string         `json:"last_name"`
	Age                     int            `json:"age"`
	PhoneNumber             sql.NullString `json:"phone_number"`
	PhoneVerificationStatus bool           `json:"phone_verification_status"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
}

var db *sql.DB

func main() {
	dsn := "dockeruser:dockerpass@tcp(localhost:3306)/hands_on_go?parseTime=true"

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/get-user", getUser)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is running. Use /get-user?id={user_id} to fetch a user.\n")
}

func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	if _, err := strconv.Atoi(idStr); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	query := `SELECT id, first_name, last_name, age, phone_number,
                     phone_verification_status, created_at, updated_at
              FROM users WHERE id = ?`

	var user User
	err := db.QueryRow(query, idStr).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Age,
		&user.PhoneNumber,
		&user.PhoneVerificationStatus,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Printf("Database query/scan error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
