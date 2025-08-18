package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"otus-project/internal/config"
	"otus-project/internal/model"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Exchange names
	FeedEventsExchange = "feed.events"

	// Queue names
	FeedMaterializationQueue = "feed.materialization"

	// Routing key patterns
	FeedEventRoutingKey = "feed.event.%s" // feed.event.{user_id}
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  config.RabbitMQConfig
}

// NewClient создает новый клиент RabbitMQ
func NewClient(cfg config.RabbitMQConfig) (*Client, error) {
	conn, err := amqp.Dial(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	client := &Client{
		conn:    conn,
		channel: ch,
		config:  cfg,
	}

	// Настраиваем exchange и очереди
	if err := client.setupExchangeAndQueues(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to setup exchange and queues: %w", err)
	}

	return client, nil
}

// setupExchangeAndQueues настраивает exchange и очереди
func (c *Client) setupExchangeAndQueues() error {
	// Объявляем exchange для событий ленты
	err := c.channel.ExchangeDeclare(
		FeedEventsExchange, // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Объявляем очередь для материализации ленты
	_, err = c.channel.QueueDeclare(
		FeedMaterializationQueue, // name
		true,                     // durable
		false,                    // delete when unused
		false,                    // exclusive
		false,                    // no-wait
		amqp.Table{
			"x-message-ttl": int32(24 * 60 * 60 * 1000), // 24 часа в миллисекундах
		}, // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Привязываем очередь к exchange с routing key для всех пользователей
	err = c.channel.QueueBind(
		FeedMaterializationQueue, // queue name
		"feed.event.*",           // routing key
		FeedEventsExchange,       // exchange
		false,                    // no-wait
		nil,                      // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	return nil
}

// PublishFeedEvent публикует событие ленты для конкретного пользователя
func (c *Client) PublishFeedEvent(ctx context.Context, userID string, event *model.FeedEvent) error {
	// Создаем routing key для конкретного пользователя
	routingKey := fmt.Sprintf(FeedEventRoutingKey, userID)

	// Сериализуем событие
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// Публикуем сообщение
	err = c.channel.PublishWithContext(ctx,
		FeedEventsExchange, // exchange
		routingKey,         // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published feed event for user %s: %s", userID, event.PostID)
	return nil
}

// PublishFeedUpdateTask публикует задачу обновления ленты
func (c *Client) PublishFeedUpdateTask(ctx context.Context, task *model.FeedUpdateTask) error {
	// Сериализуем задачу
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Публикуем сообщение в очередь материализации
	err = c.channel.PublishWithContext(ctx,
		"",                       // exchange (default)
		FeedMaterializationQueue, // routing key
		false,                    // mandatory
		false,                    // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
			Priority:     uint8(task.Priority),
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	log.Printf("Published feed update task for user %s: %s", task.UserID, task.PostID)
	return nil
}

// ConsumeFeedMaterializationTasks потребляет задачи материализации ленты
func (c *Client) ConsumeFeedMaterializationTasks(ctx context.Context, handler func(context.Context, *model.FeedUpdateTask) error) error {
	msgs, err := c.channel.Consume(
		FeedMaterializationQueue, // queue
		"",                       // consumer
		false,                    // auto-ack
		false,                    // exclusive
		false,                    // no-local
		false,                    // no-wait
		nil,                      // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var task model.FeedUpdateTask
				if err := json.Unmarshal(msg.Body, &task); err != nil {
					log.Printf("Error unmarshaling task: %v", err)
					msg.Nack(false, false)
					continue
				}

				if err := handler(ctx, &task); err != nil {
					log.Printf("Error processing task: %v", err)
					msg.Nack(false, true) // requeue
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// Close закрывает соединение с RabbitMQ
func (c *Client) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			log.Printf("Error closing connection: %v", err)
		}
	}
	return nil
}
