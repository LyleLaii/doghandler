package notifiers

import (
	"bytes"
	"doghandler/modules"
	"doghandler/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Webhook send alert by webhook
type Webhook struct {
	URL string `mapstructure:"url"`
}

// SendMessage send alert message
func (w Webhook) SendMessage(d *modules.Dog) (err error) {
	client := &http.Client{}
	data := make(map[string]interface{})
	data["servicename"] = d.Name
	data["description"] = d.Description
	data["status"] = "Firing"
	data["time"] = time.Now()
	data["lastreceiverd"] = d.Lastreceived
	data["message"] = fmt.Sprintf("service %s:%s long time did not receive watchdog message,last received time: %v", d.ServiceID, d.Name, d.Lastreceived.Format("2006-01-02 15:04:05"))
	jsonData, _ := json.Marshal(data)

	requestPost, err := http.NewRequest("POST", w.URL, bytes.NewReader(jsonData))
	resp, err := client.Do(requestPost)
	if err != nil {
		utils.LogWarn("WebHook", fmt.Sprintf("server %s:%s post to %s error: %s!", d.ServiceID, d.Name, w.URL, err.Error()))

		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		bodyContent, _ := ioutil.ReadAll(resp.Body)
		utils.LogWarn("WebHook", fmt.Sprintf("server %s:%s post to %s status: %d , resdata: %s", d.ServiceID, d.Name, w.URL, resp.StatusCode, string(bodyContent)))
	} else {
		utils.LogInfo("WebHook", fmt.Sprintf("server %s:%s post to %s status: %d", d.ServiceID, d.Name, w.URL, resp.StatusCode))
	}
	return
}
