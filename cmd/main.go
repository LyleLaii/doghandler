package main

import (
	"doghandler/notify"
	"doghandler/pkg"
	"doghandler/routers"
	"doghandler/utils"
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg = pflag.StringP("config", "c", "", "doghandler config file path")
)

func main() {
	pflag.Parse()
	utils.InitConfig(*cfg)
	utils.InitLogger()
	var global notify.GlobalConf
	var services []notify.ServiceConf
	var receivers []notify.ReceiverConf

	if err := viper.UnmarshalKey("global", &global); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.UnmarshalKey("services", &services); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.UnmarshalKey("receivers", &receivers); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	receivermap := notify.TransReceiver(receivers)

	pkg.InitDogs(global, services, receivermap)
	// fmt.Printf("%+v", pkg.Dogs)

	router := routers.InitRouter()

	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	s.ListenAndServe()
}
