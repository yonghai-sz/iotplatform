package main

import (
	"context"
	"flag"
	"fmt"

	"iotplatform/services/mqtt-subscriber/internal/config"
	"iotplatform/services/mqtt-subscriber/internal/logic"
	"iotplatform/services/mqtt-subscriber/internal/svc"
	"iotplatform/services/mqtt-subscriber/mqtt"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/proc"
)

var configFile = flag.String("f", "etc/mqtt-subscriber.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.MustSetUp()

	svcCtx := svc.NewServiceContext(c)
	items := buildSubscribeItems(svcCtx)

	dial := mqtt.DialOptions{
		Broker:         c.Mqtt.Broker,
		ClientID:       c.Mqtt.ClientID,
		Username:       c.Mqtt.Username,
		Password:       c.Mqtt.Password,
		ConnectTimeout: c.Mqtt.ConnectTimeout,
	}

	mgr, err := mqtt.NewManager(dial, items, c.Mqtt.DisconnectMs)
	if err != nil {
		logx.Must(err)
	}

	proc.AddShutdownListener(func() { mgr.Disconnect() })

	fmt.Printf("mqtt-subscriber connected to %s, client_id=%s, subscriptions=%d\n", c.Mqtt.Broker, c.Mqtt.ClientID, len(items))

	<-proc.Done()
}

func buildSubscribeItems(ctx *svc.ServiceContext) []mqtt.SubscriptionItem {
	var items []mqtt.SubscriptionItem
	for _, s := range ctx.Config.Mqtt.Subscriptions {
		if s.Topic == "" {
			continue
		}
		qos := byte(0)
		if s.Qos >= 0 && s.Qos <= 2 {
			qos = byte(s.Qos)
		}
		topic := s.Topic
		items = append(items, mqtt.SubscriptionItem{
			Topic:      topic,
			Qos:        qos,
			Callback:   messageCallback(),
			RetryTimes: 0,
		})
	}
	return items
}

func messageCallback() gomqtt.MessageHandler {
	return func(_ gomqtt.Client, msg gomqtt.Message) {
		go func() {
			if err := logic.HandleEventUp(context.Background(), msg.Topic(), msg.Payload()); err != nil {
				logx.Errorf("handle message topic=%s: %v", msg.Topic(), err)
			}
		}()
	}
}
