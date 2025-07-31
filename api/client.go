package api

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/types"
	"net"

	"go.uber.org/zap"
)

// Client encapsulates the logic for sending requests.
type Client struct {
	Hostname string
	Port     string
	Logger   *zap.Logger
}

// NewClient initializes a new Client.
func NewClient(hostname, port string, logger *zap.Logger) *Client {
	return &Client{
		Hostname: hostname,
		Port:     port,
		Logger:   logger,
	}
}

// SendRequest sends a serialized payload to the server.
func (c *Client) SendRequest(payload types.Payload) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", c.Hostname, c.Port))
	if err != nil {
		c.Logger.Error("Failed to connect to server", zap.Error(err))
		return err
	}
	defer conn.Close()

	data, err := payload.Serialize()
	if err != nil {
		c.Logger.Error("Failed to serialize payload", zap.Error(err))
		return err
	}

	_, err = conn.Write(data)
	if err != nil {
		c.Logger.Error("Failed to send data to server", zap.Error(err))
		return err
	}

	c.Logger.Info("Request sent successfully")
	return nil
}
