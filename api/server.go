package api

import (
	"context"
	"go.uber.org/zap"
	"net"
)

type Server struct {
	listener *net.Listener
	logger   *zap.Logger
}

func NewServer(logger *zap.Logger) Server {
	var s Server = Server{}
	s.logger = logger
	return s
}

func (s *Server) Start(host string, port string) {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		s.logger.Error("Server is Down")
		// if this function fails then close the whole context
		cancel()
	}()

	s.logger.Info("Server is starting")
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		s.logger.Fatal("Failed to start server", zap.Error(err))
	}
	defer func() {
		_ = listener.Close()
	}()
	s.listener = &listener

	s.logger.Info("Server is Listening on",
		zap.String("host", host),
		zap.String("port", port),
	)
	// TODO: implement workers pool
	for {

		conn, err := listener.Accept()
		if err != nil {
			s.logger.Fatal("Unable to accept connection", zap.Error(err))
		}

		// handle connection
		go handleConnection(ctx, s.logger, conn)
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
	logger.Info("Received request", zap.String("from", conn.RemoteAddr().String()))
	payload, err := decode(logger, conn)
	if err != nil {
		logger.Fatal("Unable to decode the request payload", zap.Error(err))
	}

	logger.Info(payload.String())
	// TODO: handle request based on header path.
}
