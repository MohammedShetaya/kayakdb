package main

import "github.com/MohammedShetaya/kayakdb/utils/log"

func main() {
	logger := log.InitLogger()
	logger.Info("Hello World!")
}
