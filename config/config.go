package config

import (
	"doghandler/pkg/logger"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// GlobalConf global as default config
type GlobalConf struct {
	Interval time.Duration `mapstructure:"interval"`
	Maxcount int           `mapstructure:"maxcount"`
	Receiver string        `mapstructure:"receiver"`
}

// ServiceConf service config
type ServiceConf struct {
	Name        string        `mapstructure:"name"`
	ServiceID   string        `mapstructure:"serviceid"`
	Description string        `mapstructure:"description,omitempty"`
	Interval    time.Duration `mapstructure:"interval,omitempty"`
	Maxcount    int           `mapstructure:"maxcount,omitempty"`
	Receiver    string        `mapstructure:"receiver,omitempty"`
}

// ReadConfig analyze config file
func ReadConfig() (GlobalConf, []ServiceConf) {
	var global GlobalConf
	var services []ServiceConf
	// var receivers []Receiver
	if err := viper.UnmarshalKey("global", &global); err != nil {
		logger.LogPanic("DogTimer", fmt.Sprintf("Error reading config file, %s", err))
	}

	if err := viper.UnmarshalKey("services", &services); err != nil {
		logger.LogPanic("DogTimer", fmt.Sprintf("Error reading config file, %s", err))
	}

	temp := make(map[string]bool)
	for _, v := range services {
		t := v.ServiceID
		if !temp[t] {
			temp[t] = true
		} else {
			logger.LogPanic("DogTimer", fmt.Sprintf("Find duplicate serviceid %s", t))
		}
	}

	return global, services
}
