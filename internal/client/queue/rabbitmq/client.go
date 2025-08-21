package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"otus-project/internal/config"
	"otus-project/internal/model"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	// Exchange names
	FeedEventsExchange = "feed.events"

	// Queue names
	FeedMaterializationQueue = "feed.materialization"
	FeedWebsocketQueuePrefix = "feed.websocket."

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

	// НЕ привязываем очередь материализации к exchange, так как задачи публикуются в default exchange
	// Очередь материализации работает независимо от feed.events exchange

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
	log.Printf("DEBUG: Publishing task - UserID: '%s', PostID: '%s'", task.UserID, task.PostID)

	// Сериализуем задачу
	body, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	log.Printf("DEBUG: Serialized task JSON: %s", string(body))

	// Создаем уникальный ID сообщения для предотвращения дубликатов
	messageID := fmt.Sprintf("%s:%s:%d", task.UserID, task.PostID, time.Now().UnixNano())

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
			MessageId:    messageID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish task: %w", err)
	}

	log.Printf("Published feed update task for user %s: %s (message ID: %s)", task.UserID, task.PostID, messageID)
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
				log.Printf("DEBUG: Received message from queue, delivery tag: %d, message ID: %s", msg.DeliveryTag, msg.MessageId)

				var task model.FeedUpdateTask
				if err := json.Unmarshal(msg.Body, &task); err != nil {
					log.Printf("Error unmarshaling task: %v", err)
					log.Printf("DEBUG: Raw message body: %s", string(msg.Body))
					msg.Nack(false, false)
					continue
				}

				log.Printf("DEBUG: Unmarshaled task - UserID: '%s', PostID: '%s', Priority: %d", task.UserID, task.PostID, task.Priority)

				// Проверяем валидность задачи перед обработкой
				if task.UserID == "" || task.PostID == "" {
					log.Printf("ERROR: Invalid task received - UserID: '%s', PostID: '%s', rejecting message", task.UserID, task.PostID)
					msg.Nack(false, false) // Не переотправляем невалидные сообщения
					continue
				}

				if err := handler(ctx, &task); err != nil {
					log.Printf("Error processing task: %v", err)
					msg.Nack(false, true) // requeue
				} else {
					log.Printf("DEBUG: Successfully processed task, acknowledging message (ID: %s)", msg.MessageId)
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}

// ConsumeFeedEvents потребляет события ленты и передает userID из routing key
func (c *Client) ConsumeFeedEvents(ctx context.Context, handler func(context.Context, string, *model.FeedEvent) error) error {
	// Объявляем временную очередь для этого потребителя
	q, err := c.channel.QueueDeclare(
		"",    // name - server-named
		true,  // durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare ws queue: %w", err)
	}

	// Подписка на все feed.event.*
	if err := c.channel.QueueBind(q.Name, "feed.event.*", FeedEventsExchange, false, nil); err != nil {
		return fmt.Errorf("failed to bind ws queue: %w", err)
	}

	msgs, err := c.channel.Consume(q.Name, "", true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to start consuming ws events: %w", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var ev model.FeedEvent
				if err := json.Unmarshal(msg.Body, &ev); err != nil {
					log.Printf("Error unmarshaling feed event: %v", err)
					continue
				}
				// routing key вида feed.event.{user_id}
				rk := msg.RoutingKey
				parts := strings.Split(rk, ".")
				if len(parts) < 3 {
					continue
				}
				userID := parts[2]
				if err := handler(ctx, userID, &ev); err != nil {
					log.Printf("Error handling ws event for user %s: %v", userID, err)
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
