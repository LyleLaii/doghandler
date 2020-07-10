package webhook

import (
	"bytes"
	"doghandler/config"
	"doghandler/pkg/logger"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Webhook send alert by webhook
type Webhook struct {
	name   string
	config *config.WebhookConfig
	client *http.Client
}

// New return a http instance
// just return a simple client for test
func New(name string, c *config.WebhookConfig) (*Webhook, error) {
	return &Webhook{
		name:   name,
		config: c,
		client: &http.Client{},
	}, nil
}

// Notify send alert message
func (w Webhook) Notify(m *config.Message) {
	data := make(map[string]interface{})
	data["servicename"] = m.Name
	data["description"] = m.Description
	data["status"] = "Firing"
	data["time"] = time.Now()
	data["lastreceiverd"] = m.Lastreceived
	data["message"] = fmt.Sprintf("service %s:%s long time did not receive watchdog message,last received time: %v", m.ServiceID, m.Name, m.Lastreceived.Format("2006-01-02 15:04:05"))
	jsonData, _ := json.Marshal(data)

	requestPost, err := http.NewRequest("POST", w.config.URL, bytes.NewReader(jsonData))
	resp, err := w.client.Do(requestPost)
	if err != nil {
		logger.LogWarn("WebHook", fmt.Sprintf("server %s:%s post to %s error: %s!", m.ServiceID, m.Name, w.config.URL, err.Error()))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		logger.LogWarn("WebHook", fmt.Sprintf("server %s:%s post to %s status: %d , resdata: %s", m.ServiceID, m.Name, w.config.URL, resp.StatusCode, string(bodyContent)))
	} else {
		logger.LogInfo("WebHook", fmt.Sprintf("server %s:%s post to %s status: %d", m.ServiceID, m.Name, w.config.URL, resp.StatusCode))
	}
	return
}
