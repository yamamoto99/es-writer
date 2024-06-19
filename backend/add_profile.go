package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type UserProfile struct {
    UserID     int    `json:"user_id"`
    Bio        string `json:"bio"`
    Experience string `json:"experience"`
    Skills     string `json:"skills"`
}

func saveProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int) // JWTからユーザーIDを取得
	var profile UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	profile.UserID = userID

	// デバッグ用ログ
	log.Printf("Received profile: %+v", profile)

	// プロフィール情報をデータベースに保存
	_, err := db.Exec("INSERT INTO user_profiles (user_id, bio, experience, skills) VALUES ($1, $2, $3, $4)", profile.UserID, profile.Bio, profile.Experience, profile.Skills)
	if err != nil {
		http.Error(w, "Failed to save profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// プロフィール取得用ハンドラー
func getProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var profile UserProfile
	err := db.QueryRow("SELECT bio, experience, skills FROM user_profiles WHERE user_id = $1", userID).Scan(&profile.Bio, &profile.Experience, &profile.Skills)
	if err != nil {
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(profile)
}

// authenticate関数
func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// クッキーからトークンを取得
		cookie, err := r.Cookie("idToken")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// トークンをパースおよび検証
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_jwt_secret_key"), nil
		})
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			// コンテキストにユーザーIDを設定して次のハンドラに渡す
			ctx := context.WithValue(r.Context(), "userID", claims.Username)
			next(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}
}


// func main() {
//     // 他のハンドラーの登録
//     http.HandleFunc("/saveProfile", authenticate(saveProfile))
//     http.HandleFunc("/getProfile", authenticate(getProfile))
//     log.Fatal(http.ListenAndServe(":8080", nil))
// }