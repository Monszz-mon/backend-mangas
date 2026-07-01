package handler

import (
	"fmt"
	"mangades-backend/model"   // <-- DIUBAH
	"mangades-backend/service" // <-- DIUBAH
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MangaHandler struct {
	service *service.MangaDexService
}

func NewMangaHandler() *MangaHandler {
	return &MangaHandler{
		service: service.NewMangaDexService(),
	}
}

func sendSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success",
		"data":    data,
	})
}

func sendError(c *gin.Context, status int, errMsg string) {
	c.JSON(status, gin.H{
		"status": status,
		"error":  errMsg,
	})
}

func (h *MangaHandler) SearchManga(c *gin.Context) {
	var req model.SearchMangaRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		sendError(c, http.StatusBadRequest, "Invalid query parameters")
		return
	}

	if req.Title == "" {
		sendError(c, http.StatusBadRequest, "Query parameter 'title' is required")
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	if req.Page < 1 {
		req.Page = 1
	}

	result, err := h.service.SearchManga(c.Request.Context(), req.Title, req.Limit, req.Page)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

func (h *MangaHandler) GetMangaDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		sendError(c, http.StatusBadRequest, "Manga ID is required")
		return
	}

	result, err := h.service.GetMangaByID(c.Request.Context(), id)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if result.Data == nil {
		sendError(c, http.StatusNotFound, "Manga not found")
		return
	}

	sendSuccess(c, result.Data)
}

func (h *MangaHandler) GetMangaFeed(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		sendError(c, http.StatusBadRequest, "Manga ID is required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	lang := c.DefaultQuery("lang", "en")

	if limit > 100 {
		limit = 100
	}
	if page < 1 {
		page = 1
	}

	result, err := h.service.GetMangaFeed(c.Request.Context(), id, limit, page, lang)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

func (h *MangaHandler) GetChapterPages(c *gin.Context) {
	chapterID := c.Param("id")
	if chapterID == "" {
		sendError(c, http.StatusBadRequest, "Chapter ID is required")
		return
	}

	quality := c.DefaultQuery("quality", "data-saver")

	atHome, err := h.service.GetAtHomeServer(c.Request.Context(), chapterID)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	baseURL := atHome.BaseURL
	hash := atHome.Chapter.Hash

	var imageFiles []string
	if quality == "data-saver" {
		imageFiles = atHome.Chapter.DataSaver
	} else {
		imageFiles = atHome.Chapter.Data
	}

	imageURLs := make([]string, len(imageFiles))
	for i, filename := range imageFiles {
		imageURLs[i] = fmt.Sprintf("/api/image?url=%s/%s/%s/%s", baseURL, quality, hash, filename)
	}

	response := gin.H{
		"chapterId": chapterID,
		"baseUrl":   baseURL,
		"hash":      hash,
		"quality":   quality,
		"pages":     imageURLs,
	}

	sendSuccess(c, response)
}