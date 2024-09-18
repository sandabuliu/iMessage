package queue

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"iMessage/client/redis"
	"time"
)

const (
	KafkaQueueType = "kafka"
	RedisQueueType = "redis"
)

var Client Queue

// Queue 消息队列接口
type Queue interface {
	Produce(ctx context.Context, message string) error
	Consume(ctx context.Context) (string, error)
}

func InitQueue(queueType, addr, passwd, keyOrTopic string) error {
	switch queueType {
	case KafkaQueueType:
		Client = NewKafkaQueue([]string{addr}, keyOrTopic)
		return nil
	case RedisQueueType:
		Client = NewRedisQueue(addr, passwd, keyOrTopic)
		return nil
	default:
		return fmt.Errorf("unsupported queue type: %s", queueType)
	}
}

// RedisQueue Redis消息队列
type RedisQueue struct {
	key string
}

func NewRedisQueue(addr, password, key string) *RedisQueue {
	redis.InitRedisClient(addr, password, 0)
	return &RedisQueue{key: key}
}

func (r *RedisQueue) Produce(ctx context.Context, message string) error {
	return redis.LPush(ctx, r.key, message)
}

func (r *RedisQueue) Consume(ctx context.Context) (string, error) {
	result, err := redis.BRPop(ctx, r.key, 1*time.Second)
	if err != nil {
		return "", err
	}
	return result[1], nil
}

// KafkaQueue Kafka消息队列
type KafkaQueue struct {
	writer *kafka.Writer
	reader *kafka.Reader
}

func NewKafkaQueue(brokers []string, topic string) *KafkaQueue {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		MaxWait: 1 * time.Second,
		GroupID: "Group-id",
	})
	return &KafkaQueue{writer: writer, reader: reader}
}

func (k *KafkaQueue) Produce(ctx context.Context, message string) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Value: []byte(message),
	})
}

func (k *KafkaQueue) Consume(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	msg, err := k.reader.ReadMessage(ctx)

	if err != nil {
		return "", err
	}

	//if err := k.reader.CommitMessages(ctx, msg); err != nil {
	//	fmt.Println("Failed to commit messages:", err)
	//}
	return string(msg.Value), nil
}
