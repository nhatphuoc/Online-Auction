package stream

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"search-service/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
)

type Consumer struct {
	client       *redis.Client
	streamKey    string
	groupName    string
	consumerName string
}

func NewConsumer(redisAddr, password string, db int, streamKey, groupName, consumerName string) (*Consumer, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Kết nối Redis thành công!")

	consumer := &Consumer{
		client:       client,
		streamKey:    streamKey,
		groupName:    groupName,
		consumerName: consumerName,
	}

	if err := consumer.createConsumerGroup(ctx); err != nil {
		log.Printf("Warning: Could not create consumer group: %v", err)
	}

	return consumer, nil
}

func (c *Consumer) createConsumerGroup(ctx context.Context) error {
	err := c.client.XGroupCreateMkStream(ctx, c.streamKey, c.groupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

func (c *Consumer) ReadMessages(ctx context.Context, handler func(*models.Event) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    c.groupName,
				Consumer: c.consumerName,
				Streams:  []string{c.streamKey, ">"},
				Count:    10,
				Block:    time.Second * 2,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}
				log.Printf("Error reading from stream: %v", err)
				time.Sleep(time.Second)
				continue
			}

			for _, stream := range streams {
				for _, message := range stream.Messages {
					event, err := c.parseEvent(message.Values)
					if err != nil {
						log.Printf("Error parsing event: %v", err)
						c.client.XAck(ctx, c.streamKey, c.groupName, message.ID)
						continue
					}

					if err := handler(event); err != nil {
						log.Printf("Error handling event: %v", err)
						continue
					}

					c.client.XAck(ctx, c.streamKey, c.groupName, message.ID)
				}
			}
		}
	}
}

func (c *Consumer) parseEvent(values map[string]interface{}) (*models.Event, error) {
	eventData, ok := values["event"].(string)
	if !ok {
		return nil, fmt.Errorf("missing event data")
	}

	var event models.Event
	if err := json.Unmarshal([]byte(eventData), &event); err != nil {
		return nil, fmt.Errorf("error unmarshaling event: %w", err)
	}

	return &event, nil
}

func (c *Consumer) Close() error {
	return c.client.Close()
}
