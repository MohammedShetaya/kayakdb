package api

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
)

type Server struct {
	Host               string
	Port               string
	handlersController *HandlersController
	listener           *net.Listener
	logger             *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	server := new(Server)
	server.logger = logger
	return server
}

func (s *Server) Start(host string, port string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		s.logger.Error("Server is Down")
		// if this function fails then close the whole context
		cancel()
	}()
	// initialize the protocol types
	InitProtocol()
	s.logger.Info("Server is starting")
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		s.logger.Fatal("Failed to start server", zap.Error(err))
	}
	defer func() {
		_ = listener.Close()
	}()
	s.listener = &listener
	s.Host = host
	s.Port = port
	s.handlersController = NewHandlerController(&ctx, s.logger)

	s.logger.Info("Server is Listening on",
		zap.String("host", host),
		zap.String("port", port),
	)
	// TODO: implement workers pool
	for {

		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Unable to Accept connection", zap.Error(err))
		}

		// handle connection
		go s.handleConnection(&ctx, s.logger, conn)
	}
}

func (s *Server) handleConnection(ctx *context.Context, logger *zap.Logger, conn net.Conn) {
	// TODO: connection heartbeats and timeout
	defer func() {
		_ = conn.Close()
	}()

	logger.Info("Received request", zap.String("from", conn.RemoteAddr().String()))

	var data []byte
	buffer := make([]byte, 1024)

	for {
		// if the server context is canceled then exit
		select {
		case <-(*ctx).Done():
			return
		default:
		}

		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				logger.Fatal("Error occurred while reading payload data", zap.Error(err))
			}
			break
		}
		data = append(data, buffer[:n]...)

		if uint32(len(data)*8) > MaxPayloadSize {
			logger.Fatal(fmt.Sprintf("Exceeded maximum payload size of %v", MaxPayloadSize))
		}
	}

	logger.Info("data received", zap.String("data", string(data)))

	var payload Payload
	err := payload.Deserialize(data)

	if err != nil {
		logger.Fatal("Failed to deserialize payload", zap.Error(err))
	}

	logger.Info("Payload Successfully Deserialized", zap.String("payload", payload.String()))

	// TODO: handle request based on header path.
	s.handlersController.HandleRequest(&payload)
}
