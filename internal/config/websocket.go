package config

import (
	"net"
	"os"
)

const (
	websocketHostEnvName = "WEBSOCKET_HOST"
	websocketPortEnvName = "WEBSOCKET_PORT"
)

type WebSocketConfig interface {
	Address() string
}

type websocketConfig struct {
	host string
	port string
}

func NewWebSocketConfig() (WebSocketConfig, error) {
	host := os.Getenv(websocketHostEnvName)
	if len(host) == 0 {
		// По умолчанию используем localhost
		host = "localhost"
	}

	port := os.Getenv(websocketPortEnvName)
	if len(port) == 0 {
		// По умолчанию используем порт 8080
		port = "8090"
	}

	return &websocketConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *websocketConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
