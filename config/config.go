package config

import (
	"doghandler/modules"
	"doghandler/notifiers"
	"doghandler/utils"
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// GlobalConf global as default config
type GlobalConf struct {
	Interval int    `mapstructure:"interval"`
	Maxcount int    `mapstructure:"maxcount"`
	Receiver string `mapstructure:"receiver"`
}

// ServiceConf service config
type ServiceConf struct {
	Name        string `mapstructure:"name"`
	ServiceID   string `mapstructure:"serviceid"`
	Description string `mapstructure:"description,omitempty"`
	Interval    int    `mapstructure:"interval,omitempty"`
	Maxcount    int    `mapstructure:"maxcount,omitempty"`
	Receiver    string `mapstructure:"receiver,omitempty"`
}

// ReceiverConf receiver config
type ReceiverConf struct {
	Name           string              `mapstructure:"name"`
	WebhookConfigs []notifiers.Webhook `mapstructure:"webhook_config,omitempty"`
	WechatConfigs  []notifiers.Wechat  `mapstructure:"wechat_config,omitempty"`
}

// InitDogs initiate Dogs slices
func InitDogs() {
	global, services, receiversmap := readConfig()
	for _, s := range services {
		interval := global.Interval
		maxcount := global.Maxcount
		receiver := global.Receiver
		if s.Interval != 0 {
			interval = s.Interval
		}
		if s.Maxcount != 0 {
			maxcount = s.Maxcount
		}
		if s.Receiver != "" {
			receiver = s.Receiver
		}

		modules.Dogs[s.ServiceID] = &modules.Dog{
			Name:         s.Name,
			Description:  s.Description,
			ServiceID:    s.ServiceID,
			Interval:     interval,
			Maxcount:     maxcount,
			Notifiers:    receiversmap[receiver].([]modules.Notifier),
			Counter:      0,
			Lastreceived: time.Time{},
			Timer:        nil,
		}
	}
}

func transReceiver(rcs []ReceiverConf) map[string]interface{} {
	var receivermap = make(map[string]interface{})
	for _, rc := range rcs {

		r := []modules.Notifier{}

		for _, n := range rc.WebhookConfigs {
			r = append(r, n)
		}
		for _, n := range rc.WechatConfigs {
			r = append(r, n)
		}
		receivermap[rc.Name] = r
	}
	return receivermap
}

func readConfig() (GlobalConf, []ServiceConf, map[string]interface{}) {
	var global GlobalConf
	var services []ServiceConf
	var receivers []ReceiverConf
	if err := viper.UnmarshalKey("global", &global); err != nil {
		utils.LogError("DogTimer", fmt.Sprintf("Error reading config file, %s", err))
	}

	if err := viper.UnmarshalKey("services", &services); err != nil {
		utils.LogError("DogTimer", fmt.Sprintf("Error reading config file, %s", err))
	}

	temp := make(map[string]bool)
	for _, v := range services {
		t := v.ServiceID
		if !temp[t] {
			temp[t] = true
		} else {
			utils.LogError("DogTimer", fmt.Sprintf("Find duplicate serviceid %s", t))
			os.Exit(1)
		}
	}

	if err := viper.UnmarshalKey("receivers", &receivers); err != nil {
		utils.LogError("DogTimer", fmt.Sprintf("Error reading config file, %s", err))
	}

	temp = make(map[string]bool)
	for _, v := range receivers {
		t := v.Name
		if !temp[t] {
			temp[t] = true
		} else {
			utils.LogError("DogTimer", fmt.Sprintf("Find duplicate receiver name %s", t))
			os.Exit(1)
		}
	}

	receiversmap := transReceiver(receivers)
	return global, services, receiversmap

}
