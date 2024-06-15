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

/** トークンから指定されたユーザーデータを取得
 * 取得したい値を引数に設定、"tokenはidToken"を使用
 * id -> "sub"
 * username -> "cognito:username"
 * email -> "email"
 */
func getValueFromToken(r *http.Request, cookieName, claimKey string) (string, error) {
	// クッキーからトークンを取得
	token, err := parseTokenFromCookie(r, cookieName)
	if err != nil {
		return "", err
	}

	// クレームから指定されたキーの値を取得
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if value, ok := claims[claimKey].(string); ok {
			return value, nil
		}
		return "", fmt.Errorf("クレームからキーの値を取得できませんでした: claims: %v, key: %s", claims, claimKey)
	}
	return "", fmt.Errorf("トークンが無効です: token: %v", token)
}
