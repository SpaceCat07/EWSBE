package http

import (
	"EWSBE/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	dataHandler *DataHandler
	authHandler *AuthHandler
	newsHandler *NewsHandler
	r           *gin.Engine
}

func NewHandler(dataUc *usecase.DataUsecase, authUc *usecase.AuthUsecase, newsUc *usecase.NewsUsecase) *Handler {
	r := gin.Default()

	dataHandler := NewDataHandler(dataUc)
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
	// Data routes
	h.r.POST("/data", h.dataHandler.CreateData)

	// Auth routes
	h.r.POST("/register", h.authHandler.Register)
	h.r.POST("/login", h.authHandler.Login)

	// News routes
	newsGroup := h.r.Group("/news")
	newsGroup.GET("", h.newsHandler.GetAllNews)
	newsGroup.GET("/:slug", h.newsHandler.GetNewsBySlug)

	// Protected routes
	authorized := h.r.Group("/news")
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
