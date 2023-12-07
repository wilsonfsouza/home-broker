package main

import (
	"encoding/json"
	"fmt"
	"sync"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/wilsonfsouza/home-broker/go/internal/infra/kafka"
	"github.com/wilsonfsouza/home-broker/go/internal/market/dto"
	"github.com/wilsonfsouza/home-broker/go/internal/market/entity"
	"github.com/wilsonfsouza/home-broker/go/internal/market/transformer"
)

func main() {
	ordersInput := make(chan *entity.Order)
	ordersOutput := make(chan *entity.Order)
	waitGroup := &sync.WaitGroup{}
	defer waitGroup.Wait()

	kafkaMessageChannel := make(chan *ckafka.Message)

	configMap := &ckafka.ConfigMap{
		"bootstrap.servers": "host.docker.internal:9094",
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	}

	producer := kafka.NewKafkaProducer(configMap)
	kafka := kafka.NewKafkaConsumer(configMap, []string{"input"})

	go kafka.Consume(kafkaMessageChannel) // T2

	book := entity.NewBook(ordersInput, ordersOutput, waitGroup)
	go book.Trade() // T3

	go func() {
		for kafkaMessage := range kafkaMessageChannel {
			waitGroup.Add(1)
			fmt.Println(string(kafkaMessage.Value))
			tradeInput := dto.TradeInput{}
			error := json.Unmarshal(kafkaMessage.Value, &tradeInput)

			if error != nil {
				panic(error)
			}

			order := transformer.TransformOrderInput(tradeInput)
			ordersInput <- order
		}
	}() // T4

	for res := range ordersOutput {
		orderOutput := transformer.TransformOrderOutput(res)
		orderOutputJson, error := json.MarshalIndent(orderOutput, "", "   ")
		fmt.Println(orderOutput)
		if error != nil {
			fmt.Println(error)
		}

		producer.Publish(orderOutputJson, []byte("orders"), "output")
	}
}
