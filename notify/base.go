package notify

// Notifier Notifier interface
type Notifier interface {
	SendMessage(name string, description string) error
}

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
	Name           string    `mapstructure:"name"`
	WebhookConfigs []Webhook `mapstructure:"webhook_config,omitempty"`
	WechatConfigs  []Wechat  `mapstructure:"wechat_config,omitempty"`
}

// TransReceiver Trans receiver config to map
func TransReceiver(rcs []ReceiverConf) map[string]interface{} {
	var receivermap = make(map[string]interface{})
	for _, rc := range rcs {

		r := []Notifier{}

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
