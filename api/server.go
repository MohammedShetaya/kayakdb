package api

import (
	"context"
	"fmt"
	"github.com/MohammedShetaya/kayakdb/config"
	"github.com/MohammedShetaya/kayakdb/raft"
	"github.com/MohammedShetaya/kayakdb/types"
	"go.uber.org/zap"
	"io"
	"net"
)

type Server struct {
	handlersController *HandlersController
	logger             *zap.Logger
	config             *config.Configuration
}

func NewServer(config *config.Configuration, logger *zap.Logger) *Server {
	server := new(Server)
	server.logger = logger
	server.config = config
	return server
}

func (s *Server) Start() {

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		s.logger.Error("Server is Down")
		// if this function fails then close the whole context
		cancel()
	}()
	// initialize the protocol types
	types.RegisterDataTypes()

	s.logger.Info("Server is starting")
	listener, err := net.Listen("tcp", ":"+s.config.KayakPort)
	if err != nil {
		s.logger.Fatal("Failed to start server", zap.Error(err))
	}

	defer func() {
		_ = listener.Close()
	}()

	raftLib := raft.NewRaft(s.config, s.logger)
	go raftLib.Start()

	handler, e := NewHandlerController(&ctx, raftLib, s.logger)

	if e != nil {
		s.logger.Fatal("Failed to init request controller", zap.Error(err))
	}

	s.handlersController = handler

	s.logger.Info("Server is Listening on",
		zap.String("port", s.config.KayakPort),
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
				logger.Fatal("Error occurred while reading payload types", zap.Error(err))
			}
			break
		}
		data = append(data, buffer[:n]...)

		if uint32(len(data)*8) > types.MaxPayloadSize {
			logger.Fatal(fmt.Sprintf("Exceeded maximum payload size of %v", types.MaxPayloadSize))
		}
	}

	var payload types.Payload
	err := payload.Deserialize(data)

	if err != nil {
		logger.Fatal("Failed to deserialize payload", zap.Error(err))
	}

	logger.Info("Received Request", zap.String("from", conn.RemoteAddr().String()), zap.String("payload", payload.String()))

	// Handle request > build a response > send it back
	if res, err := s.handlersController.HandleRequest(&payload); err != nil {
		logger.Error("Failed to handle client request", zap.Error(err))
		failureRes := types.Payload{
			Headers: types.Headers{
				Status:  1,
				Message: err.Error(),
			},
		}
		serialized, _ := failureRes.Serialize()
		_, err = conn.Write(serialized)

	} else {
		res.Headers.Status = 0
		serialized, e := res.Serialize()

		if e != nil {
			logger.Error("Failed to serialize response", zap.Error(e))
		} else {
			if _, e := conn.Write(serialized); e != nil {
				logger.Error("Failed to send response to client", zap.Error(e))
			}
		}
	}
}
