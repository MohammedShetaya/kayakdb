package api

import (
	"context"
	"go.uber.org/zap"
	"net"
)

type Server struct {
	Listener net.Listener
	logger   *zap.Logger
}

func (s *Server) Start(host string, port string, logger *zap.Logger) {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Error("Server is Down")
		// if this function fails then close the whole context
		cancel()
	}()

	logger.Info("Server is starting")
	Listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	defer func() {
		_ = Listener.Close()
	}()

	logger.Info("Server is Listening on",
		zap.String("host", host),
		zap.String("port", port),
	)
	// TODO: implement workers pool
	for {

		conn, err := Listener.Accept()
		if err != nil {
			logger.Fatal("Unable to accept connection", zap.Error(err))
		}

		// handle connection
		go handleConnection(ctx, logger, conn)
	}
}

func handleConnection(ctx context.Context, logger *zap.Logger, conn net.Conn) {
	// TODO: connection heartbeats and timeout
	defer func() {
		_ = conn.Close()
	}()
	// if the server context is canceled then exit
	select {
	case <-ctx.Done():
		return
	default:
	}

	logger.Info("Handling request")
}
