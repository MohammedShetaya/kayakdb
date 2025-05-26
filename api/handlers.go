package api

import (
	"context"
	"github.com/MohammedShetaya/kayakdb/raft"
	"github.com/MohammedShetaya/kayakdb/types"
	"go.uber.org/zap"
)

type HandlersController struct {
	handlers map[string]RequestHandler
	raft     *raft.Raft
	ctx      *context.Context
	logger   *zap.Logger
}

// TODO: change the return type to be of type response
type RequestHandler func(ctx *context.Context, logger *zap.Logger, payload *types.Payload) error

func NewHandlerController(ctx *context.Context, raft *raft.Raft, logger *zap.Logger) *HandlersController {
	controller := &HandlersController{
		handlers: make(map[string]RequestHandler),
		raft:     raft,
		ctx:      ctx,
		logger:   logger,
	}

	err := controller.RegisterHandlers()
	if err != nil {
		logger.Fatal("Error initializing handlers", zap.Error(err))
	}
	return controller
}

func (c *HandlersController) RegisterHandler(path string, handler RequestHandler) {
	c.handlers[path] = handler
}

func (c *HandlersController) HandleRequest(payload *types.Payload) {
	handler, exist := c.handlers[payload.Headers.Path]
	if !exist {
		c.logger.Fatal("No Handler for the request path", zap.String("path", payload.Headers.Path))
	}

	err := handler(c.ctx, c.logger, payload)
	if err != nil {
		c.logger.Fatal("Unable to handle Request", zap.Error(err))
	}
}

func (c *HandlersController) RegisterHandlers() error {
	c.RegisterHandler("/get", GetHandler)
	c.RegisterHandler("/put", PutHandler)
	return nil
}

func GetHandler(ctx *context.Context, logger *zap.Logger, payload *types.Payload) error {
	// TODO: implement the get logic
	logger.Info("Handling request", zap.String("path", payload.Headers.Path))
	return nil
}

func PutHandler(ctx *context.Context, logger *zap.Logger, payload *types.Payload) error {
	// TODO: implement the put logic
	logger.Info("Handling request", zap.String("path", payload.Headers.Path))
	return nil
}

func (c *HandlersController) PutRPC(ctx *context.Context, logger *zap.Logger, payload *types.Payload) error {
	logger.Info("Handling request", zap.String("path", payload.Headers.Path))
	return nil
}
