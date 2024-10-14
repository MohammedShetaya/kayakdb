package e2e

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/utils/log"
	"net"
	"os"
	"testing"
	"time"
)

func startTestServer(t *testing.T) *api.Server {
	// Set environment variables for the server
	os.Setenv("KAYAK_HOSTNAME", "localhost")
	os.Setenv("KAYAK_PORT", "8080")

	// Initialize the logger
	logger := log.InitLogger()
	defer func() {
		_ = logger.Sync()
	}()

	var server *api.Server
	// Start the server in a separate goroutine
	go func() {
		server = api.NewServer(logger)
		server.Start("localhost", "8080")
	}()

	// Give the server some time to start
	time.Sleep(2 * time.Second)
	return server
}

// The actual test case that starts the server and sends a payload
func TestServerCanReceivesPayload(t *testing.T) {
	// Start the server
	server := startTestServer(t)

	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", server.Host, server.Port))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	payload := &api.Payload{
		Headers: api.Headers{
			PathLength: 5,
			Path:       "/get",
		},
		Data: []api.KeyValue{
			{Key: api.Number([]byte{0x00, 0x00, 0x00, 0x0F}), Value: api.String("hello")},
			{Key: api.Number([]byte{0x00, 0x00, 0x00, 0x02}), Value: api.Binary([]byte{0x00, 0x00, 0x00, 0x01})}, // Changed from Number to Binary)},
		},
	}

	data, err := payload.Serialize()
	if err != nil {
		t.Fatalf("Failed to serialize payload: %v", err)
	}

	_, err = conn.Write(data)
	conn.Close()

	if err != nil {
		t.Fatalf("Failed to send payload: %v", err)
	}

	t.Log("Payload sent to server.")
}
