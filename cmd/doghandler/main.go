package main

import (
	"doghandler/config"
	"doghandler/dogtimer"
	"doghandler/notifiers"
	"doghandler/pkg/logger"
	"doghandler/server/api"
	"doghandler/server/api/v1"
	"doghandler/server/middlewares"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/oklog/run"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const SERVERNAME = "doghandler"

var (
	cfg = pflag.StringP("config", "c", "", "doghandler config file path")
)

func main() {
	pflag.Parse()
	//utils.InitConfig(*cfg)
	logger.InitLogger(SERVERNAME)

	if *cfg != "" {
		viper.SetConfigFile(*cfg)

	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("settings")
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		logger.LogPanic(SERVERNAME, fmt.Sprintf("read config failed: %s", err))
	}


	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.LogInfo(SERVERNAME, "config file changed")
		//TODO: shouldn't reload config because config may be wrong, fix it.
		reloadDogs()
	})

	generateDogs()

	gin.SetMode(viper.GetString("runmode"))

	handler := gin.Default()
	handler.Use(middlewares.LoggerToFile())

	api.RegisterApi(handler)
	v1.RegisterApi(handler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	var g run.Group

	g.Add(
		func() error {
			return s.ListenAndServe()
		},
		func(err error) {
			logger.LogInfo(SERVERNAME, "server close")
			s.Close()
		})

	{
		cancel := make(chan struct{})
		g.Add(func() error {
			return interrupt(cancel)
		}, func(error) {
			close(cancel)
		})
	}
	if err := g.Run(); err != nil {
		logger.LogInfo(SERVERNAME, fmt.Sprintf("%+v", errors.Wrapf(err, "command run failed")))
		os.Exit(1)
	}
	logger.LogInfo(SERVERNAME, "exiting")
}

func generateDogs() {
	receivers := config.InitReceivers()
	notifiers.InitNotifyIntergrationMap(receivers)
	dogtimer.InitDogs(dogtimer.ReadConfig())
}

func reloadDogs() {
	dogtimer.ClearDogs()
	generateDogs()
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case s := <-c:
		logger.LogInfo(SERVERNAME, fmt.Sprintf("caught signal %s . Beginng to exit.", s))
		return nil
	case <-cancel:
		return errors.New("canceled")
	}
}