package main

import (
	"log"
	"mangades-backend/config"  // <-- DIUBAH
	"mangades-backend/handler" // <-- DIUBAH

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	r := gin.Default()

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

	mangaHandler := handler.NewMangaHandler()
	imageProxyHandler := handler.NewImageProxyHandler()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"service": "MangaDex Backend Proxy",
			"version": "1.0.0",
			"endpoints": []string{
				"GET  /api/manga/search?title={title}&limit={limit}&page={page}",
				"GET  /api/manga/{id}",
				"GET  /api/manga/{id}/chapters?limit={limit}&page={page}&lang={lang}",
				"GET  /api/chapter/{id}/pages?quality={data|data-saver}",
				"GET  /api/image?url={encoded_url}",
			},
		})
	})

	api := r.Group("/api")
	{
		api.GET("/manga/search", mangaHandler.SearchManga)
		api.GET("/manga/:id", mangaHandler.GetMangaDetail)
		api.GET("/manga/:id/chapters", mangaHandler.GetMangaFeed)
		api.GET("/chapter/:id/pages", mangaHandler.GetChapterPages)
		api.GET("/image", imageProxyHandler.ProxyImage)
	}

	port := config.AppConfig.Port
	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Printf("📡 Connected to MangaDex API: %s", config.AppConfig.MangaDexBaseURL)
	r.Run(":" + port)
}