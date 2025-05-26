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

	config := &config.Configuration{}
	_, err := utils.LoadConfigurations(config)
	if err != nil {
		_ = fmt.Errorf("Failed to load configurations: %v\n", err)
		os.Exit(1)
	}

	logger := log.InitLogger(config.LogLevel)
	defer func() {
		_ = logger.Sync()
	}()

	server := api.NewServer(config, logger)
	server.Start()
}

//
//func main() {
//
//	f1 := func(a int) error {
//		if a == 5 {
//			return fmt.Errorf("error happened")
//		}
//		return nil
//	}
//
//	f2 := func(e interface{}) {
//		c := e.(error)
//		fmt.Println("executing")
//		fmt.Println(c)
//	}
//
//	i1 := []reflect.Value{reflect.ValueOf(2)}
//
//	out1 := reflect.ValueOf(f1).Call(i1)
//
//	out2 := reflect.ValueOf(f2).Call(out1)
//
//	fmt.Println(out2)
//}
