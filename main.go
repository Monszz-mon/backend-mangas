package main

import (
	"log"
	"mangades-backend/config"
	"mangades-backend/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	config.LoadConfig()

	// Setup Gin
	r := gin.Default()

	// CORS middleware - izinkan semua origin untuk development
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Inisialisasi handler
	mangaHandler := handler.NewMangaHandler()
	imageProxyHandler := handler.NewImageProxyHandler()

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "MangaDex Backend Proxy",
			"version": "1.0.0",
			"endpoints": []string{
				"GET  /api/manga/search?title={title}&limit={limit}&page={page}",
				"GET  /api/manga/{id}",
				"GET  /api/manga/{id}/chapters?limit={limit}&page={page}&lang={lang}",
				"GET  /api/chapter/{id}/pages?quality={data|data-saver}",
				"GET  /api/image?url={encoded_url}  (proxy gambar, jangan panggil langsung dari frontend)",
			},
		})
	})

	// API routes
	api := r.Group("/api")
	{
		api.GET("/manga/search", mangaHandler.SearchManga)
		api.GET("/manga/:id", mangaHandler.GetMangaDetail)
		api.GET("/manga/:id/chapters", mangaHandler.GetMangaFeed)
		api.GET("/chapter/:id/pages", mangaHandler.GetChapterPages)

		// Proxy image - WAJIB untuk menghindari hotlinking
		api.GET("/image", imageProxyHandler.ProxyImage)
	}

	// Start server
	port := config.AppConfig.Port
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Printf("📡 Connected to MangaDex API: %s", config.AppConfig.MangaDexBaseURL)
	r.Run(":" + port)
}