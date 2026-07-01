package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ImageProxyHandler struct {
	client *http.Client
}

func NewImageProxyHandler() *ImageProxyHandler {
	return &ImageProxyHandler{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (h *ImageProxyHandler) ProxyImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' parameter"})
		return
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", imageURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	req.Header.Set("User-Agent", "MangaDex-Backend/1.0")
	req.Header.Set("Referer", "https://mangadex.org")

	resp, err := h.client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch image"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Image server returned error"})
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")

	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream image"})
		return
	}
}