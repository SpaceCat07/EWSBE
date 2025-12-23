package http

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/usecase"
	ws "EWSBE/internal/websocket"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type DataHandler struct {
	dataUc *usecase.DataUsecase
	hub    *ws.Hub
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// allow all origins
		return true
	},
}

func NewDataHandler(dataUc *usecase.DataUsecase, hub *ws.Hub) *DataHandler {
	return &DataHandler{dataUc: dataUc, hub: hub}
}

func (h *DataHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"wsClients": h.hub.ClientCount(),
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

// HandleWebSocket upgrades HTTP connection to WebSocket
func (h *DataHandler) HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	client := ws.NewClient(h.hub, conn)
	h.hub.Register(client)

	// start client read/write pumps
	go client.WritePump()
	go client.ReadPump()
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

	// broadcast to WebSocket clients
	if h.hub != nil {
		h.hub.BroadcastSensorData(&d)
	}

	c.JSON(http.StatusCreated, d)
}

func (h *DataHandler) GetAllData(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100
	}

	data, err := h.dataUc.GetDataByLimit(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *DataHandler) GetLatestData(c *gin.Context) {
	data, err := h.dataUc.GetLatestData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No data available"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *DataHandler) GetDataHistory(c *gin.Context) {
	startStr := c.Query("start")
	endStr := c.Query("end")
	interval := c.DefaultQuery("interval", "raw")

	var start, end time.Time
	var err error

	if startStr == "" {
		start = time.Now().Add(-24 * time.Hour)
	} else {
		start, err = time.Parse(time.RFC3339, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format (use RFC3339)"})
			return
		}
	}

	if endStr == "" {
		end = time.Now()
	} else {
		end, err = time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end time format (use RFC3339)"})
			return
		}
	}

	if interval != "raw" && interval != "" {
		data, err := h.dataUc.GetAggregatedData(interval, start, end)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"start":    start.Format(time.RFC3339),
			"end":      end.Format(time.RFC3339),
			"interval": interval,
			"count":    len(data),
			"data":     data,
		})
		return
	}

	data, err := h.dataUc.GetDataByTimeRange(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"start": start.Format(time.RFC3339),
		"end":   end.Format(time.RFC3339),
		"count": len(data),
		"data":  data,
	})
}

func (h *DataHandler) GetDataInsights(c *gin.Context) {
	insights, err := h.dataUc.GetDataInsights()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, insights)
}
