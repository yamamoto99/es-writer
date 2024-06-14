package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// AIに送信するリクエストの構造体
type AiRequest struct {
	Contents []Content `json:"contents"`
}

// リクエスト内のコンテンツ部分
type Content struct {
	Parts []Part `json:"parts"`
}

// コンテンツ部分の中身の文章
type Part struct {
	Text string `json:"text"`
}

// AIからのレスポンスを受け取る
type AiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func sendToAi(ctx context.Context, question string) (string, error) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("error: GOOGLE_API_KEY environment variable not set")
	}

	// エンドポイントURLを設定
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=%s", apiKey)
	// 質問を含むリクエストボディをJSON形式に変換
	reqBody, err := json.Marshal(AiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: question},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	// url先に質問(reqBody)を送るオブジェクト作成 ctx=リクエストのサイクルを制御する、タイムアウトやキャンセルなど 
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// httpクライアントの初期化
	client := &http.Client{}
	// httpリクエスト(req)を送信
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// レスポンスの中身の読み取り
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 確認用
	// fmt.Printf("HTTP Status: %d\n", resp.StatusCode)
	// fmt.Printf("Response Body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: request %d: %s", resp.StatusCode, body)
	}

	// レスポンスをjson形式からAiREsponseの構造体の型に直す
	var geminiResp AiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	// 中身がある(正しく返却された)時
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("no answer found")
}

func generatePromptWithBio(bio string, questions []string) string {
	prompt := fmt.Sprintf("あなたの経歴は以下です。%s\n以下の質問に答えてください。\n", bio)
	for i, question := range questions {
		prompt += fmt.Sprintf("%d. %s\n", i+1, question)
	}
	return prompt
}

func main() {
	// HTMLファイルの読み込み
	filePath := "es_sample.html"
	htmlContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error: read file")
		return
	}

	// 質問の抽出
	questions := extractQuestions(string(htmlContent))

	// 経歴情報を定義
	bio := "大学一年生の頃に海外で英語を一年学び、その後、大学でプログラミングの勉強をし、今は個人開発などをしている。将来的にはエンジニアとしてさまざまな開発に携わりたい。"

	// 質問と経歴情報からプロンプトを生成
	prompt := generatePromptWithBio(bio, questions)

	// AIに質問を送信して回答を得る
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	answer, err := sendToAi(ctx, prompt)
	if err != nil {
		log.Fatalf("Error sending to AI: %v", err)
	}

	fmt.Printf("AI Response:\n%s\n", answer)
}
