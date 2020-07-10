package main

import (
	"doghandler/config"
	"doghandler/modules"
	"doghandler/notifiers"
	"doghandler/pkg/logger"
	"doghandler/routers"
	"doghandler/utils"
	"net/http"

	"github.com/spf13/pflag"
)

var (
	cfg = pflag.StringP("config", "c", "", "doghandler config file path")
)

func main() {
	pflag.Parse()
	utils.InitConfig(*cfg)
	logger.InitLogger()

	receivers := config.InitReceivers()
	notifiers.InitNotifyIntergrationMap(receivers)
	modules.InitDogs(config.ReadConfig())

	router := routers.InitRouter()

	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	s.ListenAndServe()
}
