package service

import (
	"context"
	"encoding/json"
	"fmt"
	"mangades-backend/config"
	"mangades-backend/model"
	"mangades-backend/utils"
	"net/http"
	"net/url"
	"time"
)

type MangaDexService struct {
	client      *http.Client
	baseURL     string
	rateLimiter *utils.RateLimiter
	userAgent   string
}

func NewMangaDexService() *MangaDexService {
	return &MangaDexService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:     config.AppConfig.MangaDexBaseURL,
		rateLimiter: utils.NewRateLimiter(config.AppConfig.MangaDexRateLimit),
		userAgent:   "MangaDex-Backend/1.0 (https://github.com/your-repo)",
	}
}

// doRequest adalah method internal untuk request ke MangaDex API dengan rate limiting
func (s *MangaDexService) doRequest(ctx context.Context, endpoint string, queryParams map[string]string) (*http.Response, error) {
	// Terapkan rate limiting
	s.rateLimiter.Wait()

	// Bangun URL
	fullURL := fmt.Sprintf("%s%s", s.baseURL, endpoint)
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}

	// Tambahkan query params
	q := parsedURL.Query()
	for key, value := range queryParams {
		if value != "" {
			q.Set(key, value)
		}
	}
	parsedURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "application/json")

	return s.client.Do(req)
}

// SearchManga - mencari manga berdasarkan judul
func (s *MangaDexService) SearchManga(ctx context.Context, title string, limit, page int) (*model.MangaDexResponse, error) {
	if limit == 0 {
		limit = 20
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	params := map[string]string{
		"title":   title,
		"limit":   fmt.Sprintf("%d", limit),
		"offset":  fmt.Sprintf("%d", offset),
		"order[relevance]": "desc",
		"includes[]": "cover_art",
	}

	resp, err := s.doRequest(ctx, "/manga", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MangaDex API error: %s", resp.Status)
	}

	var result model.MangaDexResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMangaByID - mendapatkan detail manga
func (s *MangaDexService) GetMangaByID(ctx context.Context, id string) (*model.MangaDexResponse, error) {
	params := map[string]string{
		"includes[]": "cover_art,author,artist",
	}

	resp, err := s.doRequest(ctx, fmt.Sprintf("/manga/%s", id), params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MangaDex API error: %s", resp.Status)
	}

	var result model.MangaDexResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetMangaFeed - mendapatkan daftar chapter dari manga
func (s *MangaDexService) GetMangaFeed(ctx context.Context, mangaID string, limit, page int, translatedLanguage string) (*model.MangaDexResponse, error) {
	if limit == 0 {
		limit = 20
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	params := map[string]string{
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
		"order[chapter]": "desc",
		"order[volume]":  "desc",
		"includes[]": "scanlation_group",
	}

	if translatedLanguage != "" {
		params["translatedLanguage[]"] = translatedLanguage
	} else {
		params["translatedLanguage[]"] = "en" // default English
	}

	resp, err := s.doRequest(ctx, fmt.Sprintf("/manga/%s/feed", mangaID), params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MangaDex API error: %s", resp.Status)
	}

	var result model.MangaDexResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAtHomeServer - mendapatkan base URL dan daftar file gambar untuk chapter
func (s *MangaDexService) GetAtHomeServer(ctx context.Context, chapterID string) (*model.AtHomeResponse, error) {
	resp, err := s.doRequest(ctx, fmt.Sprintf("/at-home/server/%s", chapterID), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("MangaDex API error: %s", resp.Status)
	}

	var result model.AtHomeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}