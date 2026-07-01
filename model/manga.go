package model

// Response standar dari MangaDex API
type MangaDexResponse struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
}

// Manga
type Manga struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Attributes struct {
		Title        map[string]string `json:"title"`
		Description  map[string]string `json:"description"`
		Status       string            `json:"status"`
		Year         int               `json:"year"`
		ContentRating string           `json:"contentRating"`
		Tags         []struct {
			ID   string `json:"id"`
			Attributes struct {
				Name map[string]string `json:"name"`
			} `json:"attributes"`
		} `json:"tags"`
	} `json:"attributes"`
}

// Chapter
type Chapter struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Attributes struct {
		Volume       string `json:"volume"`
		Chapter      string `json:"chapter"`
		Title        string `json:"title"`
		TranslatedLanguage string `json:"translatedLanguage"`
		Pages        int    `json:"pages"`
	} `json:"attributes"`
	Relationships []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
	} `json:"relationships"`
}

// At-Home Server Response (untuk mendapatkan URL gambar)
type AtHomeResponse struct {
	Result  string `json:"result"`
	BaseURL string `json:"baseUrl"`
	Chapter struct {
		Hash      string   `json:"hash"`
		Data      []string `json:"data"`      // original quality
		DataSaver []string `json:"dataSaver"` // compressed quality
	} `json:"chapter"`
}

// Request DTO
type SearchMangaRequest struct {
	Title string `form:"title"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
}

type ChapterPagesRequest struct {
	ChapterID  string `uri:"id"`
	Quality    string `form:"quality"` // "data" or "data-saver"
}