package main

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/config"
	"github.com/MohammedShetaya/kayakdb/utils"
	"github.com/MohammedShetaya/kayakdb/utils/log"
	"os"
)

func main() {

	c := &config.Configuration{}
	_, err := utils.LoadConfigurations(c)
	if err != nil {
		_ = fmt.Errorf("Failed to load configurations: %v\n", err)
		os.Exit(1)
	}

	logger := log.InitLogger(c.LogLevel)
	defer func() {
		_ = logger.Sync()
	}()

	server := api.NewServer(c, logger)
	server.Start()
}
