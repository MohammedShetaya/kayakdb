package tests

import (
	"encoding/binary"
	"net"
	"os"
	"testing"
	"time"

	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/utils/log"
)

func startTestServer(t *testing.T) {
	// Set environment variables for the server
	os.Setenv("KAYAK_HOSTNAME", "localhost")
	os.Setenv("KAYAK_PORT", "8080")

	// Initialize the logger
	logger := log.InitLogger()
	defer func() {
		_ = logger.Sync()
	}()

	// Start the server in a separate goroutine
	go func() {
		server := api.NewServer(logger)
		server.Start("localhost", "8080")
	}()

	// Give the server some time to start
	time.Sleep(2 * time.Second)
}

func TestClientIntegration(t *testing.T) {
	// Start the test server
	startTestServer(t)

	// Define the address to connect to
	address := "localhost:8080"

	// Connect to the TCP server
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	// Define the data to send
	path := "/get"
	var pathLength uint32 = uint32(len([]byte(path)))
	var keyType byte = 1
	var value uint32 = 20000
	var key uint32 = 8000

	// Create a byte slice to hold the data
	var data []byte

	// Convert pathLength to byte slice and append it
	pathLengthBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(pathLengthBytes, pathLength)
	data = append(data, pathLengthBytes...)

	// Convert path to byte slice and append it
	data = append(data, []byte(path)...)

	// Convert keyType to byte slice and append it
	data = append(data, keyType)

	// Convert key to byte slice and append it
	keyBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(keyBytes, key)
	data = append(data, keyBytes...)

	// Convert value to byte slice and append it
	valueBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(valueBytes, value)
	data = append(data, valueBytes...)

	// Send data to the server
	_, err = conn.Write(data)
	if err != nil {
		t.Fatalf("Failed to send data to server: %v", err)
	}
}
