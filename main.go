package main

import (
	"fmt"
	"time"

	"github.com/micro-community/micro/v3/service/logger"

	"micro-nats/nats"

	"github.com/micro-community/micro/v3/service/broker"
)

var (
	topic = "go.micro.topic.testTopic"
)

// Example of a subscription which receives all the messages
func sub(na broker.Broker) {
	_, err := na.Subscribe(topic, func(m *broker.Message) error {
		fmt.Println("[sub] received message:", string(m.Body), "header", m.Header)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func pub(na broker.Broker) {
	tick := time.NewTicker(time.Second)
	i := 0
	for range tick.C {
		msg := &broker.Message{
			Header: map[string]string{
				"id": fmt.Sprintf("%d", i),
			},
			Body: []byte(fmt.Sprintf("%d: %s", i, time.Now().String())),
		}
		if err := na.Publish(topic, msg); err != nil {
			logger.Errorf("[pub] failed: %v", err)
		} else {
			fmt.Println("[pub] pubbed message:", string(msg.Body))
		}
		i++
	}
}

func main() {
	na := nats.NewBroker(nats.ClusterID("test-cluster"), nats.ClientID("smart-cloud-gateway"))
	if err := na.Init(); err != nil {
		logger.Fatalf("Broker Init error: %v", err)
	}

	if err := na.Connect(); err != nil {
		logger.Fatalf("Broker Connect error: %v", err)
	}

	go sub(na)
	pub(na)
	select {}
}
