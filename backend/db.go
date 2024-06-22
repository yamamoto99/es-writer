package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type saveProfileRequest struct {
	Bio        string `json:"bio"`
	Experience string `json:"experience"`
	Projects   string `json:"projects"`
}

func saveProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// トークンからユーザーIDを取得
	userID, err := getValueFromToken(r, "sub")
	if err != nil {
		http.Error(w, fmt.Sprintf("トークンからユーザーIDの取得に失敗しました: %v", err), http.StatusUnauthorized)
		return
	}

	var req saveProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("UPDATE users SET bio=$1, experience=$2, projects=$3 WHERE id=$4")
	if err != nil {
		log.Println("Database prepare statement error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(req.Bio, req.Experience, req.Projects, userID)
	if err != nil {
		log.Println("Database update error:", err)
		http.Error(w, "Database update failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Profile updated successfully for user ID %s", userID)
}
