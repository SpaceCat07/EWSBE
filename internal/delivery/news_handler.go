package http

import (
	"EWSBE/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsUc *usecase.NewsUsecase
}

func NewNewsHandler(newsUc *usecase.NewsUsecase) *NewsHandler {
	return &NewsHandler{newsUc: newsUc}
}

func (h *NewsHandler) GetAllNews(c *gin.Context) {
	news, err := h.newsUc.GetAllNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) GetNewsBySlug(c *gin.Context) {
	slug := c.Param("slug")
	news, err := h.newsUc.GetNewsBySlug(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
		return
	}

	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) CreateNews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Title       string  `json:"title" binding:"required"`
		Content     string  `json:"content" binding:"required"`
		BannerPhoto *string `json:"banner_photo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	news, err := h.newsUc.CreateNews(req.Title, req.Content, req.BannerPhoto, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, news)
}

func (h *NewsHandler) UpdateNews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	newsID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req struct {
		Title       string  `json:"title"`
		Content     string  `json:"content"`
		BannerPhoto *string `json:"banner_photo"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	news, err := h.newsUc.UpdateNews(uint(newsID), req.Title, req.Content, req.BannerPhoto, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, news)
}

func (h *NewsHandler) DeleteNews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	newsID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.newsUc.DeleteNews(uint(newsID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "news deleted"})
}
