package mqtt

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/usecase"
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

func SubscribeSensorTopic(client paho.Client, topic string, qos byte, uc *usecase.DataUsecase) error {
    if client == nil || !client.IsConnected() {
        return errors.New("mqtt client not connected")
    }

    token := client.Subscribe(topic, qos, func(_ paho.Client, msg paho.Message) {
        var d entity.SensorData
        if err := json.Unmarshal(msg.Payload(), &d); err != nil {
            log.Printf("mqtt: unmarshal payload error: %v", err)
            return
        }
        if err := uc.Create(&d); err != nil {
            log.Printf("mqtt: failed to save sensor data: %v", err)
            return
        }
        log.Printf("mqtt: saved sensor data from topic %s", msg.Topic())
    })
    token.Wait()
    return token.Error()
}