package main

import (
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

	log.Error().Msg("Starting subscribe!")

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

	topic := "#"

	callback := func(client MQTT.Client, message MQTT.Message) {
		log.Info().
			Str("topic", message.Topic()).
			Str("payload", string(message.Payload())).
			Msg("received message")
	}

	if token := mqttClient.Subscribe(topic, byte(0), callback); token.Wait() &&
		token.Error() != nil {
		log.Error().
			Str("topic", topic).
			Err(token.Error()).
			Msg("Error subscribing!")
		return
	}

	log.Info().
		Str("topic", topic).
		Msg("subscribed")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	<-sig
	log.Info().Msg("Exiting")
}
