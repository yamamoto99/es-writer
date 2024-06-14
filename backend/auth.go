package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/golang-jwt/jwt/v5"
)

// JWTのクレームの構造体を定義
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// リクエストボディの構造体を定義
var input struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func loadAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, config.WithRegion(cognitoRegion))
}

func signup(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は許可しない
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// クライアントからのリクエストボディをデコードして、ユーザー情報を取得
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 必須フィールドが空でないか確認
	if input.Username == "" || input.Password == "" || input.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// AWSセッションを作成
	cfg, err := loadAWSConfig(r.Context())
	if err != nil {
		http.Error(w, "Failed to load AWS configuration", http.StatusInternalServerError)
		return
	}

	// Cognitoサービスクライアントを作成
	svc := cognitoidentityprovider.NewFromConfig(cfg)

	// サインアップリクエストの入力を設定
	signUpInput := &cognitoidentityprovider.SignUpInput{
		ClientId: aws.String(clientId),
		Username: aws.String(input.Username),
		Password: aws.String(input.Password),
		UserAttributes: []types.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(input.Email),
			},
		},
	}

	// Cognitoにサインアップリクエストを送信
	signUpOutput, err := svc.SignUp(r.Context(), signUpInput)
	if err != nil {
		http.Error(w, "Signup failed", http.StatusInternalServerError)
		log.Println("Signup error:", err)
		return
	}

	// ユーザーIDを取得
	userID := *signUpOutput.UserSub
	createdAt := time.Now()

	// サインアップ成功後にユーザー情報をデータベースに挿入
	stmt, err := db.Prepare("INSERT INTO users (id, username, email, created_at) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Println("Database prepare statement error:", err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, input.Username, input.Email, createdAt)
	if err != nil {
		http.Error(w, "Database insert failed", http.StatusInternalServerError)
		log.Println("Database insert error:", err)
		return
	}

	// サインアップ成功のレスポンスを返す
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User signed up successfully")
}

func initiateAuth(ctx context.Context, username, password, clientID string) (*cognitoidentityprovider.InitiateAuthOutput, error) {
	// AWSセッションを作成
	cfg, err := loadAWSConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Cognitoサービスクライアントを作成
	svc := cognitoidentityprovider.NewFromConfig(cfg)

	// 認証入力パラメータを設定
	authInput := &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: types.AuthFlowTypeUserPasswordAuth,
		ClientId: aws.String(clientID),
		AuthParameters: map[string]string{
			"USERNAME": username,
			"PASSWORD": password,
		},
	}

	return svc.InitiateAuth(ctx, authInput)
}

func signin(w http.ResponseWriter, r *http.Request) {
	// POSTメソッド以外は許可しない
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// クライアントからのリクエストボディをデコードして、認証情報を取得
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Cognitoにサインインリクエストを送信
	authResp, err := initiateAuth(r.Context(), input.Username, input.Password, clientId)
	if err != nil {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		log.Println("Authentication error:", err)
		return
	}

	// JWTトークンの取得
	idToken := authResp.AuthenticationResult.IdToken
	accessToken := authResp.AuthenticationResult.AccessToken
	refreshToken := authResp.AuthenticationResult.RefreshToken

	// トークンをクッキーに保存
	http.SetCookie(w, &http.Cookie{
		Name:     "idToken",
		Value:    *idToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    *accessToken,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    *refreshToken,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	// サインイン成功のレスポンスを返す
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User signed in successfully")
}
