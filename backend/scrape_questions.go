package main

import (
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func extractQuestions(htmlContent string) []string {
	var questions []string

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		fmt.Println("error: parse html")
		return questions
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h3" { // h3タグの文章を抽出
			var questionText string
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					questionText += c.Data
				}
			}
			questions = append(questions, strings.TrimSpace(questionText))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return questions
}

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
