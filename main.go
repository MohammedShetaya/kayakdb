package main

import (
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

	server := api.NewServer(logger)
	server.Start(hostname, port)

}
