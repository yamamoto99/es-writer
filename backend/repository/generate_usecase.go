package repository

import (
	"bytes"
	"encoding/json"
	"es-app/model"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type IGenerateRepository interface {
	SendAIRequest(c echo.Context, inputQuery string, inputURL string) (string, error)
}

type generateRepository struct {
}

func NewGenerateRepository() IGenerateRepository {
	return &generateRepository{}
}

func (r *generateRepository) SendAIRequest(c echo.Context, inputQuery string, inputURL string) (string, error) {
	reqBody, err := json.Marshal(model.AiRequest{
		Contents: []model.Content{
			{
				Parts: []model.Part{
					{Text: inputQuery},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(c.Request().Context(), "POST", inputURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	var geminiResp model.AiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return "", err
	}

	return geminiResp.Candidates[0].Content.Parts[0].Text, nil
}
