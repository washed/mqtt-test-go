package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	mqtttestgo "github.com/washed/mqtt-test-go"
)

func main() {
	zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"
	log.Logger = log.Output(
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339Nano},
	)

	log.Error().Msg("Starting publish!")

	mqttOpts := mqtttestgo.GetMQTTOpts()

	mqttClient := MQTT.NewClient(mqttOpts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Error().
			Err(token.Error()).
			Msg("Error connecting to MQTT!")
		return
	}
	log.Info().Msg("connected")

	defer mqttClient.Disconnect(250)

	topic := "test-topic"
	tick := time.Tick(time.Second)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	for {
		select {
		case now := <-tick:
			payload := fmt.Sprintf("test-message: %s", now)
			token := mqttClient.Publish(topic, byte(0), false, payload)
			token.Wait()
			log.Info().Str("payload", payload).Msg("published message")
		case <-sig:
			log.Info().Msg("Exiting")
		}
	}
}
