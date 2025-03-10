package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware configures CORS for the API
func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	})
}

// SetupRouter configures the API routes
func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(CorsMiddleware())

	r.POST("/api/rag", createRag)
	r.GET("/api/rag", listRags)
	r.GET("/api/rag/:name", getRag)
	r.DELETE("/api/rag/:name", deleteRag)
	r.POST("/api/query/:name", queryRag)
	r.POST("/api/upload", handleFileUpload)

	return r
}

// Handler implementations would follow here...
