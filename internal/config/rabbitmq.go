package config

import (
	"fmt"
	"os"
	"strconv"
)

type RabbitMQConfig interface {
	Host() string
	Port() int
	Username() string
	Password() string
	VHost() string
	DSN() string
}

type rabbitMQConfig struct {
	host     string
	port     int
	username string
	password string
	vhost    string
}

func NewRabbitMQConfig() (RabbitMQConfig, error) {
	portStr := os.Getenv("RABBITMQ_PORT")
	if portStr == "" {
		portStr = "5672"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RABBITMQ_PORT: %w", err)
	}

	host := os.Getenv("RABBITMQ_HOST")
	if host == "" {
		host = "localhost"
	}

	username := os.Getenv("RABBITMQ_USERNAME")
	if username == "" {
		username = "guest"
	}

	password := os.Getenv("RABBITMQ_PASSWORD")
	if password == "" {
		password = "guest"
	}

	vhost := os.Getenv("RABBITMQ_VHOST")
	if vhost == "" {
		vhost = "/"
	}

	return &rabbitMQConfig{
		host:     host,
		port:     port,
		username: username,
		password: password,
		vhost:    vhost,
	}, nil
}

func (c *rabbitMQConfig) Host() string {
	return c.host
}

func (c *rabbitMQConfig) Port() int {
	return c.port
}

func (c *rabbitMQConfig) Username() string {
	return c.username
}

func (c *rabbitMQConfig) Password() string {
	return c.password
}

func (c *rabbitMQConfig) VHost() string {
	return c.vhost
}

func (c *rabbitMQConfig) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d%s", c.username, c.password, c.host, c.port, c.vhost)
}
