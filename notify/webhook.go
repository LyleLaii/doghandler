package notify

import (
	"bytes"
	"doghandler/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Webhook send alert by webhook
type Webhook struct {
	URL string `mapstructure:"url"`
}

// SendMessage send alert message
func (w Webhook) SendMessage(name string, description string) (err error) {
	client := &http.Client{}

	data := make(map[string]interface{})
	data["servicename"] = name
	data["description"] = description
	data["message"] = "service long time did not receive watchdog message"
	jsonData, _ := json.Marshal(data)

	requestPost, err := http.NewRequest("POST", w.URL, bytes.NewReader(jsonData))
	resp, err := client.Do(requestPost)
	if err != nil {
		utils.LogWarn("WebHook", fmt.Sprintf("server %s post to %s error: %s!", name, w.URL, err.Error()))

		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		utils.LogWarn("WebHook", fmt.Sprintf("server %s post to %s status: %d , resdata: %s", name, w.URL, resp.StatusCode, string(bodyContent)))
	} else {
		utils.LogInfo("WebHook", fmt.Sprintf("server %s post to %s status: %d , data: %s", name, w.URL, resp.StatusCode))
	}
	return
}
