package main

import (
	"bytes"
	"context"
	"regexp"
	// "fmt"
	"log"

	"golang.org/x/net/html"

	// "net/http"
	// "os"
	"strings"
)

// HTMLから質問を抽出し、AIに送信して質問かどうかを確認する関数
func filterQuestions(ctx context.Context, htmlContent string) ([]string, error) {
	bodyContent := extractBodyContent(htmlContent)
	// fmt.Println("Body Content:", bodyContent)

	//htmlを全てAIに投げて質問を出力してもらう。各質問の終わりに"\n"を入れてもらい、それで区切って質問を配列に入れる
	questions, err := sendToAi(ctx, "以下のHTMLを解析し、textareaのある質問文のみを抽出し、質問文のみを出力してください。出力する際は全ての質問を一つに繋いでください。そして、それぞれの質問文の間には#*#を入れてください。" + bodyContent)
	// fmt.Println("Questions:", questions)
	if err != nil {
		log.Fatalf("Error sending to AI: %v", err)
	}
	// 質問文を一問ずつ配列に格納
	questionArray := strings.Split(questions, "#*#")
	// for _, question := range questionArray {
	// 	fmt.Println("Valid Question:", question)
	// }
	return questionArray, nil
}


// HTMLからbodyタグの中身を抽出する関数
func extractBodyContent(htmlContent string) string {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		log.Fatalf("Failed to parse HTML: %v", err)
	}

	var bodyNode *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			bodyNode = n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if bodyNode == nil {
		log.Fatalf("Failed to find body tag")
	}

	var buf bytes.Buffer
	html.Render(&buf, bodyNode)
	return buf.String()
}

// HTMLノードを文字列にレンダリングするヘルパー関数
// func renderNode(n *html.Node) string {
// 	var buf bytes.Buffer
// 	html.Render(&buf, n)
// 	return buf.String()
// }

func cleanHTMLContent(htmlContent string) string {
	// Five-serverによって挿入されたコードを削除する正規表現
	re := regexp.MustCompile(`(?s)<!-- Code injected by Five-server -->.*?<!--.*?-->`)
	cleanedHTML := re.ReplaceAllString(htmlContent, "")

	// スクリプトタグを削除する正規表現
	reScript := regexp.MustCompile(`(?s)<script.*?>.*?</script>`)
	cleanedHTML = reScript.ReplaceAllString(cleanedHTML, "")

	// スタイルタグを削除する正規表現
	reStyle := regexp.MustCompile(`(?s)<style.*?>.*?</style>`)
	cleanedHTML = reStyle.ReplaceAllString(cleanedHTML, "")

	return cleanedHTML
}

// func main() {
// 	filePath := "es_sample.html"
// 	htmlContent, err := os.ReadFile(filePath)
// 	if err != nil {
// 		fmt.Println("error: read file")
// 		return
// 	}

// 	ctx := context.Background()
// 	questions, err := filterQuestions(ctx, string(htmlContent))
// 	if err != nil {
// 		log.Fatalf("Error extracting and filtering questions: %v", err)
// 	}

// 	for _, question := range questions {
// 		fmt.Println("Valid Question:", question)
// 	}
// }
