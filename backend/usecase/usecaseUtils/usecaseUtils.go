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
	studentExperience := fmt.Sprintf(
		"私の略歴（アルバイト、インターン、イベントなど）は%sです。" +
		"私のスキル・資格・研究内容は%sです。" +
		"自己PRは%sです。" +
		"将来の目標とキャリアプランは%sです。",
		profile.WorkExperience, profile.Skills, profile.SelfPR, profile.FutureGoals,
	)

	return fmt.Sprintf(
		"学生の質問に対し、簡潔かつ具体的に記述し、#や*,-などは使用せずに平文で回答してください。" +
		"回答は質問に書かれた文字数制限の約9割で記述してください。" +
		"学生:私の経歴に基づいて質問に回答してください。" +
		"%s。私の経歴は以下の通りです。" +
		"%s回答:",
		question, studentExperience,
	)
}
