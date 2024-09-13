package main

import (
	"context"
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/utils/log"
	"os"
)

func main() {
	logger := log.InitLogger()

	defer func() {
		_ = logger.Sync()
	}()

	hostname := os.Getenv("KAYAK_HOSTNAME")
	port := os.Getenv("KAYAK_PORT")

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		logger.Error("Server is shut down")
		cancel()
	}()

	server := api.Server{}

	go server.Start(hostname, port, ctx, cancel, logger)

	<-ctx.Done()

}
