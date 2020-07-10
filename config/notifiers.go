package config

import (
	"doghandler/pkg/logger"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Message Message format
type Message struct {
	Name         string
	ServiceID    string
	Description  string
	Lastreceived time.Time
}

// Receiver receiver config
type Receiver struct {
	Name           string          `mapstructure:"name"`
	ReceiverType   string          `mapstructure:"receivertype"`
	WebhookConfigs []WebhookConfig `mapstructure:"webhook_config,omitempty"`
	// WechatConfigs  []notifiers.Wechat  `mapstructure:"wechat_config,omitempty"`
}

// WebhookConfig send alert by webhook
type WebhookConfig struct {
	URL string `mapstructure:"url"`
}

// InitReceivers init receivers map
func InitReceivers() []Receiver {
	var receivers []Receiver

	if err := viper.UnmarshalKey("receivers", &receivers); err != nil {
		logger.LogError("GetReceivers", fmt.Sprintf("Error load receivers config: %s", err))
	}

	temp := make(map[string]bool)
	for _, v := range receivers {
		t := v.Name
		if !temp[t] {
			temp[t] = true
		} else {
			logger.LogPanic("DogTimer", fmt.Sprintf("Find duplicate receiver name %s", t))
		}
	}

	return receivers
}
