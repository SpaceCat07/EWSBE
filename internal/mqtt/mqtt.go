package mqtt

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/usecase"
	ws "EWSBE/internal/websocket"
	"encoding/json"
	"errors"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

func Connect(broker, clientID string) (paho.Client, error) {
	opts := paho.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)
	opts.SetConnectTimeout(30 * time.Second)

	client := paho.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func SubscribeSensorTopic(client paho.Client, topic string, qos byte, uc *usecase.DataUsecase, hub *ws.Hub) error {
	if client == nil || !client.IsConnected() {
		return errors.New("mqtt client not connected")
	}

	token := client.Subscribe(topic, qos, func(_ paho.Client, msg paho.Message) {
		// parse MQTT payload
		var mqttPayload entity.MQTTSensorPayload
		if err := json.Unmarshal(msg.Payload(), &mqttPayload); err != nil {
			log.Printf("mqtt: unmarshal payload error: %v", err)
			return
		}

		// convert to SensorData
		sensorData := mqttPayload.ToSensorData()

		// save to database
		if err := uc.Create(sensorData); err != nil {
			log.Printf("mqtt: failed to save sensor data: %v", err)
			return
		}
		log.Printf("mqtt: saved sensor data from topic %s (temp: %.1f°C, hum: %.1f%%)", 
			msg.Topic(), sensorData.Temperature, sensorData.Humidity)

		// broadcast to WebSocket clients in real-time
		if hub != nil {
			if err := hub.BroadcastSensorData(sensorData); err != nil {
				log.Printf("mqtt: failed to broadcast sensor data: %v", err)
			} else {
				log.Printf("mqtt: broadcasted to %d WebSocket clients", hub.ClientCount())
			}
		}
	})
	token.Wait()
	return token.Error()
}