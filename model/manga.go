package model

type MangaDexResponse struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
}

type SearchMangaRequest struct {
	Title string `form:"title"`
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
}

type AtHomeResponse struct {
	Result  string `json:"result"`
	BaseURL string `json:"baseUrl"`
	Chapter struct {
		Hash      string   `json:"hash"`
		Data      []string `json:"data"`
		DataSaver []string `json:"dataSaver"`
	} `json:"chapter"`
}