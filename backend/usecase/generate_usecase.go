package usecase

import (
	"es-app/model"
	"es-app/repository"
	"es-app/usecase/usecaseUtils"
	"os"
	"strings"
	"sync"

	"github.com/labstack/echo/v4"
)

type IGenerateUsecase interface {
	GenerateAnswers(c echo.Context, input string) ([]model.Answer, error)
}

type generateUsecase struct {
	GenerateRepo repository.IGenerateRepository
	userRepo     repository.IUserRepository
}

func NewGenerateUsecase(GenerateRepo repository.IGenerateRepository, userRepo repository.IUserRepository) IGenerateUsecase {
	return &generateUsecase{
		GenerateRepo: GenerateRepo,
		userRepo:     userRepo,
	}
}

func (gu *generateUsecase) GenerateAnswers(c echo.Context, input string) ([]model.Answer, error) {
	user, err := gu.userRepo.GetUser(c, c.Get("user_id").(string))
	if err != nil {
		return []model.Answer{}, err
	}

	var profile model.UserProfile
	profile.Bio = user.Bio
	profile.Experience = user.Experience
	profile.Projects = user.Projects

	cleanHtml := usecaseUtils.CleanHTMLContent(input)
	bodyContent, err := usecaseUtils.ExtractBodyContent(cleanHtml)
	if err != nil {
		return []model.Answer{}, err
	}

	sendMessage := `以下のHTMLを解析し、textareaのある質問文のみを抽出し、質問文のみを出力してください。出力する際は全ての質問を一つに繋いでください。そして、それぞれの質問文の間には#*#を入れてください。`
	questionsScrapeQuery := sendMessage + bodyContent
	url := `https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=` + os.Getenv("GOOGLE_API_KEY")
	res, err := gu.GenerateRepo.SendAIRequest(c, questionsScrapeQuery, url)
	if err != nil {
		return []model.Answer{}, err
	}

	questionList := strings.Split(res, "#*#")
	if len(questionList) == 0 {
		return []model.Answer{}, nil
	}

	var wg sync.WaitGroup
	answers := make([]model.Answer, len(questionList))
	for i, question := range questionList {
		wg.Add(1)
		go func(i int, question string) {
			defer wg.Done()
			prompt := usecaseUtils.GeneratePrompt(profile, question)
			answer, err := gu.GenerateRepo.SendAIRequest(c, prompt, url)
			if err != nil {
				return
			}
			answers[i] = model.Answer{Question: question, Answer: answer}
		}(i, question)
	}

	wg.Wait()
	return answers, nil
}
