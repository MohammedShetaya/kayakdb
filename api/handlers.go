package api

import (
	"context"
	"go.uber.org/zap"
)

type HandlersController struct {
	handlers map[string]RequestHandler
	ctx      *context.Context
	logger   *zap.Logger
}

// TODO: change the return type to be of type response
type RequestHandler func(ctx *context.Context, logger *zap.Logger, payload *Payload) error

func NewHandlerController(ctx *context.Context, logger *zap.Logger) *HandlersController {
	controller := &HandlersController{
		handlers: make(map[string]RequestHandler),
		ctx:      ctx,
		logger:   logger,
	}

	err := controller.InitHandlers()
	if err != nil {
		logger.Fatal("Error initializing handlers", zap.Error(err))
	}
	return controller
}

func (c *HandlersController) RegisterHandler(path string, handler RequestHandler) {
	c.handlers[path] = handler
}

func (c *HandlersController) HandleRequest(payload *Payload) {
	handler, exist := c.handlers[payload.Headers.Path]
	if !exist {
		c.logger.Fatal("No Handler for the request path", zap.String("path", payload.Headers.Path))
	}

	err := handler(c.ctx, c.logger, payload)
	if err != nil {
		c.logger.Fatal("Unable to handle Request", zap.Error(err))
	}
}

func (c *HandlersController) InitHandlers() error {
	c.RegisterHandler("/get", GetHandler)
	c.RegisterHandler("/put", PutHandler)
	return nil
}

func GetHandler(ctx *context.Context, logger *zap.Logger, payload *Payload) error {
	// TODO: implement the get logic
	logger.Info("Handling request", zap.String("path", payload.Headers.Path))
	return nil
}

func PutHandler(ctx *context.Context, logger *zap.Logger, payload *Payload) error {
	// TODO: implement the put logic
	logger.Info("Handling request", zap.String("path", payload.Headers.Path))
	return nil
}
