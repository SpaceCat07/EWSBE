package http

import (
	"EWSBE/internal/usecase"
	ws "EWSBE/internal/websocket"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	dataHandler *DataHandler
	authHandler *AuthHandler
	newsHandler *NewsHandler
	r           *gin.Engine
}

func NewHandler(dataUc *usecase.DataUsecase, authUc *usecase.AuthUsecase, newsUc *usecase.NewsUsecase, hub *ws.Hub) *Handler {
	r := gin.Default()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"} // Frontend URLs
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	dataHandler := NewDataHandler(dataUc, hub)
	authHandler := NewAuthHandler(authUc)
	newsHandler := NewNewsHandler(newsUc)

	h := &Handler{
		dataHandler: dataHandler,
		authHandler: authHandler,
		newsHandler: newsHandler,
		r:           r,
	}

	h.routes()

	return h
}

func (h *Handler) routes() {
	// WebSocket
	h.r.GET("/ws", h.dataHandler.HandleWebSocket)

	// API Group
	api := h.r.Group("/api")

	// Data Routes
	api.POST("/data", h.dataHandler.CreateData)
	api.GET("/data", h.dataHandler.GetAllData)
	api.GET("/data/latest", h.dataHandler.GetLatestData)
	api.GET("/data/history", h.dataHandler.GetDataHistory)
	api.GET("/data/insights", h.dataHandler.GetDataInsights)
	api.GET("/health", h.dataHandler.HealthCheck)

	// Auth Routes
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", h.authHandler.Register)
		authGroup.POST("/login", h.authHandler.Login)
	}

	// News Routes
	newsGroup := api.Group("/news")
	{
		newsGroup.GET("", h.newsHandler.GetAllNews)
		newsGroup.GET("/:slug", h.newsHandler.GetNewsBySlug)
	}

	// Protected routes
	authorized := api.Group("/news")
	authorized.Use(AuthMiddleware())
	{
		authorized.POST("", h.newsHandler.CreateNews)
		authorized.PUT("/:id", h.newsHandler.UpdateNews)
		authorized.DELETE("/:id", h.newsHandler.DeleteNews)
	}
}

func (h *Handler) Router() http.Handler {
	return h.r
}
