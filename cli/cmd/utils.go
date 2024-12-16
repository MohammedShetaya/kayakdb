package cmd

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"go.uber.org/zap"
	"net"
)

func SendRequest(hostname string, port string, payload api.Payload, logger *zap.Logger) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", hostname, port))
	if err != nil {
		logger.Fatal("Failed to connect to server", zap.Error(err))
	}
	data, err := payload.Serialize()
	if err != nil {
		logger.Fatal("Failed to serialize payload", zap.Error(err))
	}
	_, err = conn.Write(data)
	conn.Close()
}
