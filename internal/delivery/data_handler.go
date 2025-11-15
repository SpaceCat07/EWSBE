package http

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	uc *usecase.DataUsecase
	r  *gin.Engine
}

func NewHandler(uc *usecase.DataUsecase) *Handler {
	r := gin.Default()
	h := &Handler{uc: uc, r: r}

	h.routes()

	return h
}

func (h *Handler) routes() {
	h.r.POST("/data", h.CreateData)
}

func (h *Handler) Router() http.Handler {
	return h.r
}

func (h *Handler) CreateData(c *gin.Context) {
	var d entity.SensorData
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.uc.Create(&d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, d)
}
