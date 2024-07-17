package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

func verifySession(w http.ResponseWriter, r *http.Request) {
	// トークンを検証
	fmt.Println("welcome page called...")

	err := validateIDToken(r)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		fmt.Println("user signed in successfully")
		fmt.Fprintf(w, "User signed in successfully")
		return
	}

	// トークンが無効な場合, リフレッシュトークンによる更新を行う
	fmt.Println("attempt to get new token by refresh token...")

	refreshTokenValue, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, fmt.Sprintf("リフレッシュトークンの取得に失敗しました: %v", err), http.StatusUnauthorized)
		return
	}

	// AWS設定の読み込み
	svc, err := createCognitoClient(r.Context())
	if err != nil {
		http.Error(w, "Unable to load SDK config", http.StatusInternalServerError)
		return
	}

	// 認証入力パラメータを設定
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeRefreshToken,
		ClientId: aws.String(clientId),
		AuthParameters: map[string]string{
			"REFRESH_TOKEN": refreshTokenValue.Value,
		},
	}

	// Cognitoにサインインリクエストを送信
	authResp, err := svc.InitiateAuth(r.Context(), authInput)
	if err != nil {
		fmt.Println("Authentication error:", err)
		http.Error(w, "Failed to refresh tokens", http.StatusInternalServerError)
		return
	}

	// JWTトークンの取得
	idToken := authResp.AuthenticationResult.IdToken
	accessToken := authResp.AuthenticationResult.AccessToken

	// トークンをクッキーに保存
	http.SetCookie(w, &http.Cookie{
		Name:     "idToken",
		Value:    *idToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    *accessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})

	// サインイン成功のレスポンスを返す
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User signed in successfully")
	fmt.Println("success to get new token by refresh token")
}
