package http

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DataHandler struct {
	dataUc *usecase.DataUsecase
	r      *gin.Engine
}

func NewDataHandler(dataUc *usecase.DataUsecase) *DataHandler {
	r := gin.Default()
	h := &DataHandler{dataUc: dataUc, r: r}

	h.routes()

	return h
}

func (h *DataHandler) routes() {
	h.r.POST("/data", h.CreateData)
}

func (h *DataHandler) Router() http.Handler {
	return h.r
}

func (h *DataHandler) CreateData(c *gin.Context) {
	var d entity.SensorData
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.dataUc.Create(&d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, d)
}
