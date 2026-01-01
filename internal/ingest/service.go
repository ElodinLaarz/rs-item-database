package ingest

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"rs-item-database/pb"
)

// Service handles fetching data from the RS3 API
type Service struct {
	client      *http.Client
	rateLimiter *time.Ticker
	lastReqTime time.Time
}

// NewService creates a new Ingest Service with a strict rate limiter
func NewService() *Service {
	return &Service{
		client:      &http.Client{Timeout: 10 * time.Second},
		rateLimiter: time.NewTicker(5 * time.Second), // 5s cooldown as per guidelines
	}
}

// FetchItem fetches a single item from the RS3 API.
// It blocks until the rate limiter allows the request.
func (s *Service) FetchItem(id int) (*pb.Item, error) {
	// Wait for rate limiter
	<-s.rateLimiter.C

	url := fmt.Sprintf("https://services.runescape.com/m=itemdb_rs/api/catalogue/detail.json?item=%d", id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (RS Item Database; Local Project)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	return Transform(body)
}

// Shutdown stops the service and its tickers
func (s *Service) Shutdown() {
	s.rateLimiter.Stop()
}
