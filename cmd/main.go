package main

import (
	"fmt"
	"github.com/MohammedShetaya/kayakdb/api"
	"github.com/MohammedShetaya/kayakdb/config"
	"github.com/MohammedShetaya/kayakdb/utils"
	"os"
)

func main() {
	// load configurations
	c := &config.Configuration{}
	_, err := utils.LoadConfigurations(c)
	if err != nil {
		_ = fmt.Errorf("Failed to load configurations: %v\n", err)
		os.Exit(1)
	}

	logger := utils.InitLogger(c.LogLevel)
	defer func() {
		_ = logger.Sync()
	}()

	server := api.NewServer(c, logger)
	server.Start()
}
