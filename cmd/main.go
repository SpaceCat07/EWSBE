package main

import (
	"EWSBE/internal/config"
	"EWSBE/internal/db"
	deliver "EWSBE/internal/delivery"
	"EWSBE/internal/entity"
	"EWSBE/internal/model"
	"EWSBE/internal/mqtt"
	"EWSBE/internal/usecase"
	ws "EWSBE/internal/websocket"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %s (using defaults)", err)
	}

	cfg := config.LoadConfig()

	// initialize Cloudinary
	config.InitCloudinary()

	// open DB
	gormDB, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	// auto migrate
	if err := gormDB.AutoMigrate(&entity.SensorData{}, &entity.User{}, &entity.News{}); err != nil {
		log.Fatalf("automigrate: %v", err)
	}
	log.Println("Database migration completed")

	// Initialize WebSocket hub
	hub := ws.NewHub()
	go hub.Run()
	log.Println("WebSocket hub started")

	// wiring repo -> usecase -> handler (GIN)
	dataRepo := model.NewDataRepo(gormDB)
	dataUc := usecase.NewDataUsecase(dataRepo)

	// auth components
	userRepo := model.NewUserRepo(gormDB)
	authUc := usecase.NewAuthUsecase(userRepo)

	// news components
	newsRepo := model.NewNewsRepo(gormDB)
	newsUc := usecase.NewNewsUsecase(newsRepo)

	// unified handler
	handler := deliver.NewHandler(dataUc, authUc, newsUc, hub)

	// mqtt init
	broker := os.Getenv("MQTT_BROKER")
	clientID := os.Getenv("MQTT_CLIENT_ID")
	topic := os.Getenv("MQTT_TOPIC")

	if broker == "" {
		log.Println("Warning: MQTT_BROKER not set, skipping MQTT connection")
	} else {
		mqttClient, err := mqtt.Connect(broker, clientID)
		if err != nil {
			log.Printf("mqtt connect error: %v (continuing without MQTT)", err)
		} else {
			// subscribe to sensor topic with WebSocket hub for broadcasting
			if err := mqtt.SubscribeSensorTopic(mqttClient, topic, 0, dataUc, hub); err != nil {
				log.Printf("mqtt subscribe error: %v", err)
			} else {
				log.Printf("mqtt subscribed to topic: %s", topic)
			}

			// graceful MQTT disconnect on shutdown
			defer func() {
				if mqttClient != nil && mqttClient.IsConnected() {
					mqttClient.Disconnect(250)
					log.Println("MQTT client disconnected")
				}
			}()
		}
	}

	// http server
	addr := normalizeAddr(cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler.Router(),
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	go func() {
		log.Printf("listening on %s", addr)
		log.Printf("webSocket endpoint: ws://localhost%s/ws", addr)
		log.Printf("REST API endpoint: http://localhost%s/api", addr)
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	} else {
		log.Println("Server gracefully stopped")
	}
}

func normalizeAddr(port string) string {
	if port == "" {
		return ":8080"
	}
	if port[0] == ':' {
		return port
	}
	return ":" + port
}
