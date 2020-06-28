package main

import (
	"doghandler/config"
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
	utils.InitLogger()

	config.InitDogs()

	router := routers.InitRouter()

	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	s.ListenAndServe()
}
