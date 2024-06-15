package main

import (
	"net/http"
	"fmt"
)

// ウェルカムメッセージを表示
func welcome(w http.ResponseWriter, r *http.Request) {
	// トークンを検証
	err := validateIDToken(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("IDトークンの検証に失敗しました: %v", err), http.StatusUnauthorized)
		return
	}

	// ユーザー名を取得
	username, err := getValueFromToken(r, "cognito:username")
	if err != nil {
		http.Error(w, fmt.Sprintf("ユーザー名の取得に失敗しました: %v", err), http.StatusUnauthorized)
		return
	}

	// ウェルカムメッセージを送信
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ようこそ、%sさん", username)
}
