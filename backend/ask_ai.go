package main

import (
    // "bytes"
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "sync"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

type HtmlRequest struct {
	Html string `json:"html"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeRequest struct {
	AnthropicVersion string    `json:"anthropic_version"`
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens"`
	Temperature      float64   `json:"temperature,omitempty"`
}

type ClaudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

const modelId = "anthropic.claude-3-haiku-20240307-v1:0"

func sendToAi(ctx context.Context, question string) (string, error) {
    // AWSの認証情報を取得
    region := "us-west-2"
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        config.WithCredentialsProvider(
            credentials.NewStaticCredentialsProvider(
                os.Getenv("AWS_ACCESS_KEY_ID"),
                os.Getenv("AWS_SECRET_ACCESS_KEY"),
                os.Getenv("AWS_SESSION_TOKEN"),
            ),
        ),
    )
    fmt.Println(os.Getenv("AWS_ACCESS_KEY_ID"),os.Getenv("AWS_SECRET_ACCESS_KEY"))
    
    if err != nil {
        return "", fmt.Errorf("failed to load AWS config: %w", err)
    }

    // bedrockにリクエストを送るためのクライアント作成
    client := bedrockruntime.NewFromConfig(cfg)

    // メッセージの作成
    content := "Human: " + question + "\n\nAssistant:"

    messages := []Message{
        {
            Role:    "user",
            Content: content,
        },
    }

    // リクエストボディを作成
    reqBody, err := json.Marshal(ClaudeRequest{
        Messages:         messages,
        AnthropicVersion: "bedrock-2023-05-31",
        MaxTokens:        1000,
        Temperature:       0.2,
    })
    if err != nil {
        return "", fmt.Errorf("failed to marshal request body: %w", err)
    }

    //　質問を投げかける
    output, err := client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        reqBody,
	})
    if err != nil {
        return "", fmt.Errorf("failed to invoke model: %w", err)
    }

	fmt.Printf("Response Body: %s\n", string(output.Body))

    // レスポンスをパース
    var response ClaudeResponse
    if err := json.Unmarshal(output.Body, &response); err != nil {
        return "", fmt.Errorf("failed to unmarshal response: %w", err)
    }

	// レスポンスを返す
    if len(response.Content) > 0 {
		return response.Content[0].Text, nil
	}

	return "", fmt.Errorf("no answer found")
}

func generatePromptWithBio(bio, question string) string {
    return fmt.Sprintf("あなたの経歴は%sです。以下の質問に答えてください。簡潔かつ具体的に記述し、#や*,-などは使用せずに平文で解答部分のみを出力してください。\n%s", bio, question)
}

func processQuestionsWithAI(w http.ResponseWriter, r *http.Request) {
    // CORSヘッダーを追加
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    // OPTIONSリクエストに対する処理
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }

    // HTMLの読み込み
    var req HtmlRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
     log.Printf("Error decoding request body: %v", err)
     http.Error(w, "Bad request", http.StatusBadRequest)
     return
    }

    // // 不要な部分を取り除く
    cleanHtml := cleanHTMLContent(req.Html)
    log.Printf("Cleaned HTML: %s", cleanHtml)

    // HTMLファイルの読み込み
    // filePath := "es_sample.html"
    // htmlContent, err := os.ReadFile(filePath)
    // if err != nil {
    //     fmt.Println("error: read file")
    //     return
    // }

    // 質問の抽出
    questions := extractQuestions(string(cleanHtml))
    if len(questions) == 0 {
        log.Printf("No questions found in the HTML content")
        http.Error(w, "No questions found", http.StatusBadRequest)
        return
    }
    //TOOD htmlを投げて質問に答えさせる
    for i:=0; i < len(questions); i++{
        fmt.Println(questions[i])
    }

    // 経歴情報を定義
    bio := "大学一年生の頃に海外で英語を一年学び、その後、大学でプログラミングの勉強をし、今は個人開発などをしている。webアプリケーションも作成した。(https://github.com/yamamoto99/es-writer)将来的にはエンジニアとしてさまざまな開発に携わりたい。普段は42Tokyoに通っており、CやGoを学んでいる。"

    // 並列処理のためのWaitGroupを作成
    var wg sync.WaitGroup

    // コンテキストを設定
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    type Answer struct {
        Question string `json:"question"`
        Answer   string `json:"answer"`
    }

    answers := make([]Answer, len(questions))

    // 時間計測開始
    startTime := time.Now()

    // 質問ごとにゴルーチンを作成して並列処理を実行
    for i, question := range questions {
        wg.Add(1)
        go func(i int, q string) {
            defer wg.Done()
            prompt := generatePromptWithBio(bio, q)
            answer, err := sendToAi(ctx, prompt)
            if err != nil {
                log.Printf("Error sending to AI: %v", err)
                return
            }
            answers[i] = Answer{Question: q, Answer: answer}
        }(i, question)
    }

    // 全てのゴルーチンが終了するのを待機
    wg.Wait()

    // 時間計測(確認用)
    elapsedTime := time.Since(startTime)
    fmt.Printf("Total processing time: %s\n", elapsedTime)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(answers)
}

func main() {
    http.HandleFunc("/getAnswers", processQuestionsWithAI)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

