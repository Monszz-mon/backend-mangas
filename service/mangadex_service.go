package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
			Timeout: 60 * time.Second,
		},
		baseURL:     config.AppConfig.MangaDexBaseURL,
		rateLimiter: utils.NewRateLimiter(config.AppConfig.MangaDexRateLimit),
		userAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

func (s *MangaDexService) doRequest(ctx context.Context, endpoint string, queryParams map[string]string) (*http.Response, error) {
	s.rateLimiter.Wait()

	fullURL := fmt.Sprintf("%s%s", s.baseURL, endpoint)
	parsedURL, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}

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
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")

	return s.client.Do(req)
}

func (s *MangaDexService) doRequestWithRetry(ctx context.Context, endpoint string, queryParams map[string]string, retries int) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i < retries; i++ {
		resp, err = s.doRequest(ctx, endpoint, queryParams)
		if err == nil {
			return resp, nil
		}

		// Jika error karena koneksi, tunggu dan coba lagi
		if err.Error() == "EOF" || err.Error() == "read: connection reset" {
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			continue
		}
		break
	}
	return nil, err
}

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
        "includes[]": "cover_art",
        // Hapus order[relevance] - kadang bikin error
    }

    // Tambahkan logging untuk debugging (opsional)
    fmt.Printf("Searching manga with params: %+v\n", params)

    resp, err := s.doRequestWithRetry(ctx, "/manga", params, 3)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        // Baca response body untuk tahu detail error
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("MangaDex API error: %s - %s", resp.Status, string(body))
    }

    var result model.MangaDexResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

func (s *MangaDexService) GetMangaByID(ctx context.Context, id string) (*model.MangaDexResponse, error) {
	params := map[string]string{
		"includes[]": "cover_art,author,artist",
	}

	resp, err := s.doRequestWithRetry(ctx, fmt.Sprintf("/manga/%s", id), params, 3)
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

func (s *MangaDexService) GetMangaFeed(ctx context.Context, mangaID string, limit, page int, translatedLanguage string) (*model.MangaDexResponse, error) {
	if limit == 0 {
		limit = 20
	}
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	params := map[string]string{
		"limit":          fmt.Sprintf("%d", limit),
		"offset":         fmt.Sprintf("%d", offset),
		"order[chapter]": "desc",
		"order[volume]":  "desc",
		"includes[]":     "scanlation_group",
	}

	if translatedLanguage != "" {
		params["translatedLanguage[]"] = translatedLanguage
	} else {
		params["translatedLanguage[]"] = "en"
	}

	resp, err := s.doRequestWithRetry(ctx, fmt.Sprintf("/manga/%s/feed", mangaID), params, 3)
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

func (s *MangaDexService) GetAtHomeServer(ctx context.Context, chapterID string) (*model.AtHomeResponse, error) {
	resp, err := s.doRequestWithRetry(ctx, fmt.Sprintf("/at-home/server/%s", chapterID), nil, 3)
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