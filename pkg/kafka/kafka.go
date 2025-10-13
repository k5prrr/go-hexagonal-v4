package kafka

// go get -u github.com/confluentinc/confluent-kafka-go/kafka && go get -u github.com/confluentinc/confluent-kafka-go/v2/kafka
// https://www.youtube.com/watch?v=7Hm2RsH8bS8
// https://www.youtube.com/watch?v=KSej3yivuPY
/*
import (
	"fmt"
	//"strings"
	//"time"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout = 5000 // ms
	consumerSessionTimeout = 5000 // ms
)

type Producer struct {
	producer *kafka.Producer
}
type Consumer struct {
	consumer *kafka.Consumer
}
// Адреса через запятую
func NewProducer(address string) (*Producer, error) {
	// https://github.com/confluentinc/librdkafka/blob/master/CONFIGURATION.md
	config := &kafka.ConfigMap{
		"bootstrap.servers": address, //strings.Join(address, ","),
	}
	p, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf("Error new producer: %w", err)
	}
	return &Producer{producer: p}, nil
}

func (p *Producer) Send(topic, key, text string) error {
	var messageKey *[]byte
	if key != "" {
		keyBytes := []byte(key)
		messageKey = &keyBytes
	}

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:          &topic,
			TopicPartition: kafka.PartitionAny,
		},
		Value: []byte(text),
		Key:   messageKey,
	}
	kafkaChan := make(chan kafka.Event)
	defer close(kafkaChan)

	err := p.producer.Produce(kafkaMessage, kafkaChan)
	if err != nil {
		return fmt.Errorf("Error Send: %w", err)
	}

	event := <-kafkaChan
	switch eventResult := event.(type) {
		case *kafka.Message:
			return nil
		case *kafka.Error:
			return eventResult
		default:
			return fmt.Errorf("Error kafkaChan")
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}

func NewConsumer(address, topics, consumerGroup string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": address,
		"group.id": consumerGroup,
		"session.timeout.ms": consumerSessionTimeout,
	}
	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("Error new consumer: %w", err)
	}

	err = c.SubscribeTopics([]string{topics}, nil)
	if err != nil {
		return nil, fmt.Errorf("Error subscribing to topics: %w", err)
	}
	Сделай подписку на топики и в комментариях покажи пример как их указывать
	return &Consumer{consumer: c}, nil
}
func (с *Consumer) Start() {
	for {
		kafkaMessage, err := c.consumer,ReadMessage(-1) //(time.Second * 10) // без таймаута -1
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}
		if kafkaMessage == nil {
			continue
		}
		// Обработка сообщения
		fmt.Printf("Received message: %s", kafkaMessage.Value)
	}
}

*/
