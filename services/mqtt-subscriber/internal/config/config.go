package config

import (
	"time"

	"github.com/zeromicro/go-zero/core/service"
)

// Subscription is one MQTT subscription declared in YAML.
type Subscription struct {
	Topic string `json:",optional"`
	Qos   int    `json:",default=0"`
}

// MqttConf holds broker connection and optional subscription list.
type MqttConf struct {
	Broker         string         `json:",optional"`
	ClientID       string         `json:",optional"`
	Username       string         `json:",optional"`
	Password       string         `json:",optional"`
	ConnectTimeout time.Duration  `json:",default=25s"`
	DisconnectMs   uint           `json:",default=250"`
	Subscriptions  []Subscription `json:",optional"`
}

// Config is the mqtt-subscriber service configuration.
type Config struct {
	service.ServiceConf
	Mqtt MqttConf
}
