package main

import (
	"fmt"
	"os"
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
		if n.Type == html.ElementNode && n.Data == "h3" { //とりあえずh3タグの文章引っ張る
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					questions = append(questions, strings.TrimSpace(c.Data))
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return questions
}

func main() {
	filePath := "es_sample.html"
	htmlContent, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("error: read file")
		return
	}

	questions := extractQuestions(string(htmlContent))
	for _, question := range questions {
		fmt.Println(question)
	}
}
