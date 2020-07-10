package notifiers

import (
	"doghandler/config"
	"doghandler/notifiers/webhook"
	"doghandler/pkg/logger"
	"fmt"
)

// NotifyIntergrations notify integration slice
var NotifyIntergrations []NotifyIntegration

// Notifier notifier instance
type Notifier interface {
	Notify(m *config.Message)
}

// NotifyIntegration notify integration
type NotifyIntegration struct {
	name         string
	receitertype string
	notifier     Notifier
}

type errorInfo struct {
	name string
	err  string
}

// Name return notify name
func (n NotifyIntegration) Name() string {
	return n.name
}

// ReceiverType return notift receiver type
func (n NotifyIntegration) ReceiverType() string {
	return n.receitertype
}

// Notify notify instance send message
func (n NotifyIntegration) Notify(m *config.Message) {
	n.notifier.Notify(m)
}

// NotifyMap notify integration map
var NotifyMap map[string][]NotifyIntegration

// InitNotifyIntergrationMap generate notify integration by map format
func InitNotifyIntergrationMap(rs []config.Receiver) {
	var (
		errs        []errorInfo
		notifiermap = make(map[string][]NotifyIntegration)
		add         = func(name string, receivertype string, f func() (Notifier, error)) {
			n, e := f()
			if e != nil {
				errs = append(errs, errorInfo{
					name: name,
					err:  fmt.Sprintf("%v", e),
				})
				return
			}
			notifiermap[name] = append(notifiermap[name], NotifyIntegration{
				name:         name,
				receitertype: receivertype,
				notifier:     n,
			})
		}
	)

	for _, cs := range rs {
		for _, c := range cs.WebhookConfigs {
			add(cs.Name, cs.ReceiverType, func() (Notifier, error) { return webhook.New(cs.Name, &c) })
		}
	}

	if len(errs) > 0 {
		logger.LogPanic("Notify", fmt.Sprintf("Init notify intergration get errors: %+v", errs))
	}
	NotifyMap = notifiermap
}
