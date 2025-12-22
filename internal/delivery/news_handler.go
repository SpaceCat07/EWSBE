package http

import (
	"EWSBE/internal/config"
	"EWSBE/internal/usecase"
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	title := c.PostForm("title")
	content := c.PostForm("content")
	file, header, err := c.Request.FormFile("banner_photo")

	var bannerPhoto *string
	if err == nil && header != nil {
		// Upload to Cloudinary
		publicID := "news_" + uuid.New().String() + "_" + time.Now().Format("20060102150405")
		url, uploadErr := config.UploadImage(context.Background(), file, publicID)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
			return
		}
		bannerPhoto = &url
	}

	if title == "" || content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title and content are required"})
		return
	}

	news, err := h.newsUc.CreateNews(title, content, bannerPhoto, userID.(uint))
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

	title := c.PostForm("title")
	content := c.PostForm("content")
	file, header, err := c.Request.FormFile("banner_photo")

	var bannerPhoto *string
	if err == nil && header != nil {
		// Upload to Cloudinary
		publicID := "news_" + uuid.New().String() + "_" + time.Now().Format("20060102150405")
		url, uploadErr := config.UploadImage(context.Background(), file, publicID)
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
			return
		}
		bannerPhoto = &url
	}

	news, err := h.newsUc.UpdateNews(uint(newsID), title, content, bannerPhoto, userID.(uint))
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
