package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Wechat send alert by wechat
type Wechat struct {
	Address string `mapstructure:"address"`
}

// SendMessage send alert message
func (w Wechat) SendMessage(name string, description string) (err error) {
	client := &http.Client{}

	data := make(map[string]interface{})
	data["servicename"] = name
	data["description"] = description
	data["message"] = "service long time did not receive watchdog message"
	jsonData, _ := json.Marshal(data)

	requestPost, err := http.NewRequest("POST", w.Address, bytes.NewReader(jsonData))
	resp, err := client.Do(requestPost)
	if err != nil {
		fmt.Printf("get request failed, err:[%s]", err.Error())
		return
	}
	defer resp.Body.Close()

	bodyContent, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("resp status code:[%d]\n", resp.StatusCode)
	fmt.Printf("resp body data:[%s]\n", string(bodyContent))
	return
}
