package usecaseUtils

import (
	"bytes"
	"es-app/model"
	"fmt"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

func CleanHTMLContent(input string) string {
	re := regexp.MustCompile(`(?s)<!-- Code injected by Five-server -->.*?<!--.*?-->`)
	cleanedHTML := re.ReplaceAllString(input, "")

	reScript := regexp.MustCompile(`(?s)<script.*?>.*?</script>`)
	cleanedHTML = reScript.ReplaceAllString(cleanedHTML, "")

	reStyle := regexp.MustCompile(`(?s)<style.*?>.*?</style>`)
	cleanedHTML = reStyle.ReplaceAllString(cleanedHTML, "")

	return cleanedHTML
}

func ExtractBodyContent(cleanedHTML string) (string, error) {
	doc, err := html.Parse(strings.NewReader(cleanedHTML))
	if err != nil {
		return "", err
	}

	bodyNode := findBodyNode(doc)
	if bodyNode == nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, bodyNode); err != nil {
		return "", err
	}

	return buf.String(), err
}

func findBodyNode(n *html.Node) *html.Node {
	if n.Type == html.ElementNode && n.Data == "body" {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if bodyNode := findBodyNode(c); bodyNode != nil {
			return bodyNode
		}
	}
	return nil
}

func GeneratePrompt(profile model.UserProfile, question string) string {
	combinedBio := fmt.Sprintf("%sです。今までの経験は%sです。これまでに作ってきた作品は%s", profile.Bio, profile.Experience, profile.Projects)
	return fmt.Sprintf("あなたの経歴は%sです。以下の質問に答えてください。簡潔かつ具体的に記述し、#や*,-などは使用せずに平文で解答部分のみを出力してください。\n%s", combinedBio, question)
}
