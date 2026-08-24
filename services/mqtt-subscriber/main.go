package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"iot-zero/pkg/mqtt"
	"iot-zero/services/mqtt-subscriber/internal/config"
	"iot-zero/services/mqtt-subscriber/internal/logic"
	"iot-zero/services/mqtt-subscriber/internal/svc"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"

	"time"
)

// SubscribeItem is one subscription to apply on (re)connect.
type SubscriptionItem struct {
	Topic    string
	Qos      byte
	Callback gomqtt.MessageHandler
}

func buildSubscribeItems(ctx *svc.ServiceContext) []SubscriptionItem {
	var items []SubscriptionItem
	for _, s := range ctx.Config.Mqtt.Subscriptions {
		if s.Topic == "" {
			continue
		}
		qos := byte(0)
		if s.Qos >= 0 && s.Qos <= 2 {
			qos = byte(s.Qos)
		}
		topic := s.Topic

		items = append(items, SubscriptionItem{
			Topic:    topic,
			Qos:      qos,
			Callback: messageCallback(),
		})
	}
	return items
}

func messageCallback() gomqtt.MessageHandler {
	return func(_ gomqtt.Client, msg gomqtt.Message) {
		go func() {
			var err error
			switch {
			case strings.HasPrefix(msg.Topic(), "example_up/"):
				err = logic.HandleExampleUp(context.Background(), msg.Topic(), msg.Payload())
			default:
				logx.Infof("topic=%s payload=%s", msg.Topic(), string(msg.Payload()))
			}
			if err != nil {
				logx.Errorf("handle message topic=%s: %v", msg.Topic(), err)
			}
		}()
	}
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

var configFile = flag.String("f", "etc/mqtt-subscriber.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.MustSetUp()

	dial := mqtt.DialOptions{
		Broker:         c.Mqtt.Broker,
		ClientID:       c.Mqtt.ClientID,
		Username:       c.Mqtt.Username,
		Password:       c.Mqtt.Password,
		ConnectTimeout: c.Mqtt.ConnectTimeout,
	}

	svcCtx := svc.NewServiceContext(c)
	items := buildSubscribeItems(svcCtx)
	onConnect := func(c gomqtt.Client) {
		for _, item := range items {
			subscribe(c, item)
		}
	}

	mgr, err := mqtt.NewManager(dial, onConnect, c.Mqtt.DisconnectMs)
	if err != nil {
		logx.Must(err)
	}

	proc.AddShutdownListener(func() { mgr.Disconnect() })

	fmt.Printf("mqtt-subscriber connected to %s, client_id=%s, subscriptions=%d\n", c.Mqtt.Broker, c.Mqtt.ClientID, len(items))

	<-proc.Done()
}
