package main

import (
	"EWSBE/internal/config"
	"EWSBE/internal/db"
	deliver "EWSBE/internal/delivery"
	"EWSBE/internal/entity"
	"EWSBE/internal/model"
	"EWSBE/internal/mqtt"
	"EWSBE/internal/usecase"
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
	cfg := config.LoadConfig()

	// open DB
	gormDB, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// init Cloudinary
	config.InitCloudinary()
	// auto migrate
	if err := gormDB.AutoMigrate(&entity.SensorData{}, &entity.User{}, &entity.News{}); err != nil {
		log.Fatalf("automigrate: %v", err)
	}

	// wiring repo -> usecase -> handler (Gin)
	dataRepo := model.NewDataRepo(gormDB)
	userRepo := model.NewUserRepo(gormDB)
	newsRepo := model.NewNewsRepo(gormDB)

	dataUc := usecase.NewDataUsecase(dataRepo)
	authUc := usecase.NewAuthUsecase(userRepo)
	newsUc := usecase.NewNewsUsecase(newsRepo)

	handler := deliver.NewHandler(dataUc, authUc, newsUc)

	// mqtt init
	broker := os.Getenv("MQTT_BROKER")
	clientID := os.Getenv("MQTT_CLIENT_ID")
	topic := os.Getenv("MQTT_TOPIC")

	mqttClient, err := mqtt.Connect(broker, clientID)
	if err != nil {
		log.Printf("mqtt connect error: %v", err)
	} else {
		// unsubscribe / disconnect handled on shutdown
		if err := mqtt.SubscribeSensorTopic(mqttClient, topic, 0, dataUc); err != nil {
			log.Printf("mqtt subscribe error: %v", err)
		} else {
			log.Printf("mqtt subscribed to %s", topic)
		}
	}

	// server
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
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	// graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	// disconnect mqtt gracefully
	if mqttClient != nil && mqttClient.IsConnected() {
		mqttClient.Disconnect(250)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
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
