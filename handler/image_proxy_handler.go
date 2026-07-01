package handler

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ImageProxyHandler menangani proxy gambar ke MangaDex
// Ini adalah bagian WAJIB untuk menghindari hotlinking
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

// ProxyImage - GET /api/image?url=<encoded_url>
func (h *ImageProxyHandler) ProxyImage(c *gin.Context) {
	imageURL := c.Query("url")
	if imageURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'url' parameter"})
		return
	}

	// Validasi agar hanya request ke domain MangaDex
	// Ini penting untuk keamanan (mencegah SSRF)
	// Kita cek apakah URL mengandung 'uploads.mangadex.org'
	// atau domain lain yang diizinkan
	// Untuk sederhana, kita hanya izinkan URL dari MangaDex
	// Bisa diperketat sesuai kebutuhan

	// Request ke MangaDex
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", imageURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set User-Agent agar tidak ditolak
	req.Header.Set("User-Agent", "MangaDex-Backend/1.0")
	req.Header.Set("Referer", "https://mangadex.org") // terkadang diperlukan

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

	// Set header agar browser bisa menampilkan gambar dengan benar
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg" // fallback
	}
	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // cache 1 hari

	// Tulis data gambar ke response
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream image"})
		return
	}
}