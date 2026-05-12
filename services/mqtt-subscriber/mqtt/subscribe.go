package mqtt

import (
	"time"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

// SubscribeItem is one subscription to apply on (re)connect.
type SubscriptionItem struct {
	Topic      string
	Qos        byte
	Callback   gomqtt.MessageHandler
	RetryTimes int // 0 means retry indefinitely on failure.
}

func subscribe(client gomqtt.Client, item SubscriptionItem) {
	token := client.Subscribe(item.Topic, item.Qos, item.Callback)
	if err := token.Error(); err != nil {
		logx.Errorf("mqtt subscribe token error topic=%s: %v", item.Topic, err)
		return
	}
	ack := token.WaitTimeout(6 * time.Second)
	if !ack {
		logx.Errorf("mqtt subscribe ack timeout topic=%s", item.Topic)
	}
}
