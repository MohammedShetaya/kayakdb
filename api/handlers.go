package api

import (
	"context"
	"fmt"
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

// RequestHandler return a pointer to a Payload that represents the response and an error if any occurred while
// processing the request.
type RequestHandler func(raft *raft.Raft, logger *zap.Logger, payload *types.Payload) (*types.Payload, error)

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

func (c *HandlersController) HandleRequest(payload *types.Payload) (*types.Payload, error) {
	handler, exist := c.handlers[payload.Headers.Path.String()]
	if !exist {
		c.logger.Fatal("No Handler for the request path", zap.String("path", payload.Headers.Path.String()))
	}

	resp, err := handler(c.raft, c.logger, payload)
	if err != nil {
		c.logger.Error("Unable to handle Request", zap.Error(err))
		return nil, err
	}

	return resp, nil
}

func (c *HandlersController) RegisterHandlers() error {
	c.RegisterHandler("/get", GetHandler)
	c.RegisterHandler("/put", PutHandler)
	return nil
}

func GetHandler(r *raft.Raft, logger *zap.Logger, payload *types.Payload) (*types.Payload, error) {
	logger.Debug("Handling request", zap.String("path", payload.Headers.Path.String()))

	if len(payload.Data) == 0 {
		return nil, fmt.Errorf("get handler requires exactly one key in payload data")
	}

	key := payload.Data[0]
	value, err := r.Get(key)
	if err != nil {
		return nil, err
	}

	// At the moment the API only logs the value. A full implementation would
	// marshal a response back to the requester.
	if value == nil {
		return nil, fmt.Errorf("key not found. key: %v", key.String())
	} else {
		logger.Debug("Key retrieved", zap.String("key", key.String()), zap.String("value", value.String()))
	}

	// Build a response payload – even if value is nil we return an empty data slice
	var data []types.Type

	data = append(data, types.KeyValue{
		Key:   key,
		Value: value,
	})

	resp := &types.Payload{Data: data}

	return resp, nil
}

func PutHandler(r *raft.Raft, logger *zap.Logger, payload *types.Payload) (*types.Payload, error) {
	logger.Debug("Handling request", zap.String("path", payload.Headers.Path.String()))

	// Append the entries to the Raft log
	entries := r.Put(payload.Data)

	// Convert the committed log entries into a response payload that contains only the data field.
	var data []types.Type
	for _, entry := range entries {
		data = append(data, entry.Pair)
	}

	resp := &types.Payload{
		Data: data,
	}

	return resp, nil
}
