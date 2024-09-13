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

func (s *Server) Start(host string, port string, ctx context.Context, cancel context.CancelFunc, logger *zap.Logger) {
	// if this function fails then close the whole context
	defer cancel()
	logger.Info("Server is starting")
	Listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		defer cancel()
		logger.Fatal("Failed to start server", zap.Error(err))
	}
	defer func() {
		_ = Listener.Close()
	}()

	logger.Info("Server is Listening on",
		zap.String(host, host),
		zap.String(port, port),
	)
	// TODO: implement workers pool
	// TODO: connection heartbeats and timeout
	for {
		conn, err := Listener.Accept()
		if err != nil {
			logger.Fatal("Unable to accept connection", zap.Error(err))
		}

		// handle connection
		go func(c net.Conn) {
			defer func() {
				_ = c.Close()
			}()

		}(conn)
	}
}
