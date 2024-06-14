package main

import (
	"net/http"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/MicahParks/keyfunc/v3"
)

// トークンを検証し、パースされたトークンを返す
func parseTokenFromCookie(r *http.Request, cookieName string) (*jwt.Token, error) {
	// クッキーからトークンを取得
	c, err := r.Cookie(cookieName)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, fmt.Errorf("%vが見つかりません: %v", cookieName, err)
		}
		return nil, fmt.Errorf("不正なリクエスト：クッキーの取得エラー: %v", err)
	}

	tknStr := c.Value

	// JWKセットを指定されたURLから取得
	jwks, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("指定されたURLのリソースからJWKセットを作成できませんでした: %v", err)
	}

	// トークンをパース
	token, err := jwt.Parse(tknStr, jwks.Keyfunc)
	if err != nil {
		return nil, fmt.Errorf("JWTのパースに失敗しました: %v", err)
	}

	// トークンが有効かどうかをチェック
	if !token.Valid {
		return nil, fmt.Errorf("認証失敗：無効なトークン")
	}

	return token, nil
}

// トークンを検証し、有効であればtrueを返す
func validateIDToken(r *http.Request) error {
	// クッキーからトークンを取得
	_, err := parseTokenFromCookie(r, "accessToken")
	if err != nil {
		return err
	}
	return nil
}

// トークンからユーザー名を取得
func getUserName(r *http.Request) (string, error) {
	// クッキーからトークンを取得
	token, err := parseTokenFromCookie(r, "idToken")
	if err != nil {
		return "",  err
	}

	// クレームからユーザー名を取得
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if username, ok := claims["cognito:username"].(string); ok {
			return username, nil
		}
		return "", fmt.Errorf("クレームからユーザー名を取得できませんでした: claims: %v", claims)
	}
	return "", fmt.Errorf("トークンが無効です: token: %v", token)
}

// ウェルカムメッセージを表示
func welcome(w http.ResponseWriter, r *http.Request) {
	// トークンを検証
	err := validateIDToken(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("IDトークンの検証に失敗しました: %v", err), http.StatusUnauthorized)
		return
	}

	// ユーザー名を取得
	username, err := getUserName(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("ユーザー名の取得に失敗しました: %v", err), http.StatusUnauthorized)
		return
	}

	// ウェルカムメッセージを送信
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ようこそ、%sさん", username)
}
