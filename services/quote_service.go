package services

import (
	"auth-app/models"
	"auth-app/repositories"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type QuoteService interface {
	FetchAndSaveTodayQuote() (*models.Quote, error)
}

type quoteService struct {
	repo repositories.QuoteRepository
}

func NewQuoteService(repo repositories.QuoteRepository) QuoteService {
	return &quoteService{repo: repo}
}

func (s *quoteService) FetchAndSaveTodayQuote() (*models.Quote, error) {
	resp, err := http.Get("https://zenquotes.io/api/today")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch quote: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	var zenQuotes []struct {
		Q string `json:"q"`
		A string `json:"a"`
	}
	err = json.Unmarshal(body, &zenQuotes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	if len(zenQuotes) == 0 {
		return nil, errors.New("empty quote array returned from API")
	}

	quote := &models.Quote{
		Name:      zenQuotes[0].A,
		QuoteLine: zenQuotes[0].Q,
	}
	// Now save to database
	err = s.repo.Save(quote)
	if err != nil {
		return nil, fmt.Errorf("failed to save quote to database: %v", err)
	}

	return quote, nil
}
