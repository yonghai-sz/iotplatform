package mqtt

import (
	"fmt"
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

// DialOptions configures the MQTT client connection.
type DialOptions struct {
	Broker         string
	ClientID       string
	Username       string
	Password       string
	ConnectTimeout time.Duration
}

func connectClient(opts DialOptions, onConnect func(gomqtt.Client)) (gomqtt.Client, error) {
	if opts.ConnectTimeout <= 0 {
		opts.ConnectTimeout = 25 * time.Second
	}
	if opts.Broker == "" {
		return nil, fmt.Errorf("mqtt broker is empty")
	}

	pahoOpts := gomqtt.NewClientOptions()
	pahoOpts.SetAutoReconnect(true)

	pahoOpts.AddBroker(opts.Broker)
	pahoOpts.SetClientID(opts.ClientID)
	if opts.Username != "" {
		pahoOpts.SetUsername(opts.Username)
	}
	if opts.Password != "" {
		pahoOpts.SetPassword(opts.Password)
	}

	pahoOpts.SetOnConnectHandler(onConnect)

	client := gomqtt.NewClient(pahoOpts)
	token := client.Connect()
	ack := token.WaitTimeout(opts.ConnectTimeout)
	if token.Error() != nil {
		return nil, fmt.Errorf("connect mqtt %s: %w", opts.Broker, token.Error())
	}
	if !ack {
		return nil, fmt.Errorf("connect mqtt %s: timeout after %s", opts.Broker, opts.ConnectTimeout)
	}
	return client, nil
}
