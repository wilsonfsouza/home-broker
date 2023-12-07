package kafka

import ckafka "github.com/confluentinc/confluent-kafka-go/kafka"

type Producer struct {
	ConfigMap *ckafka.ConfigMap
}

func NewKafkaProducer(configMap *ckafka.ConfigMap) *Producer {
	return &Producer{
		ConfigMap: configMap,
	}
}

func (p *Producer) Publish(message interface{}, key []byte, topic string) error {
	producer, error := ckafka.NewProducer(p.ConfigMap)

	if error != nil {
		return error
	}

	kafkaMessage := &ckafka.Message{
		TopicPartition: ckafka.TopicPartition{
			Topic:     &topic,
			Partition: ckafka.PartitionAny,
		},
		Key:   key,
		Value: message.([]byte),
	}

	error = producer.Produce(kafkaMessage, nil)

	if error != nil {
		return error
	}

	return nil
}
