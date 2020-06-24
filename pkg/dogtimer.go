package pkg

import (
	"doghandler/notify"
	"doghandler/utils"
	"fmt"
	"sync"
	"time"
)

// Dogs dog map
var Dogs = make(map[string]*Dog)

// Dog Watchdog config
type Dog struct {
	Name         string
	ServiceID    string
	Description  string
	Interval     int
	Maxcount     int
	Counter      int
	Lastreceived time.Time
	Timer        *time.Timer
	Notifiers    []notify.Notifier
	mu           sync.Mutex
}

// ServiceInfo service infomation
// type ServiceInfo struct {
// 	name        string
// 	description string
// }

// TouchDog refresh dog info
func (d *Dog) TouchDog() {
	d.mu.Lock()
	d.Counter = 0
	d.Lastreceived = time.Now()
	d.refreshtimer()
	d.mu.Unlock()
	utils.LogInfo("DogTimer", fmt.Sprintf("service %s last received time: %s", d.ServiceID, d.Lastreceived))
}

func (d *Dog) refreshtimer() {
	if d.Timer != nil {
		d.Timer.Stop()
	}
	d.Timer = time.AfterFunc(time.Duration(d.Interval)*time.Second, d.CheckDog)
}

// CheckDog check dog status
func (d *Dog) CheckDog() {
	d.mu.Lock()
	d.Counter = d.Counter + 1
	utils.LogInfo("DogTimer", fmt.Sprintf("service %s did not receiver, now counter is : %v, max count is: %v", d.ServiceID, d.Counter, d.Maxcount))
	if d.Counter >= d.Maxcount {
		d.Alert()
		d.refreshtimer()
	} else {
		d.refreshtimer()
	}
	d.mu.Unlock()
}

// Alert Send the message that the dog is dead
func (d *Dog) Alert() {
	utils.LogInfo("DogTimer", fmt.Sprintf("service %s maximum number of not received, send alert", d.ServiceID))
	for _, n := range d.Notifiers {
		n.SendMessage(d.Name, d.Description)
	}
	d.Counter = 0
	d.refreshtimer()
}

// InitDogs initiate Dogs slices
func InitDogs(global notify.GlobalConf, services []notify.ServiceConf, receiversmap map[string]interface{}) {
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

		Dogs[s.ServiceID] = &Dog{
			Name:         s.Name,
			Description:  s.Description,
			ServiceID:    s.ServiceID,
			Interval:     interval,
			Maxcount:     maxcount,
			Notifiers:    receiversmap[receiver].([]notify.Notifier),
			Counter:      0,
			Lastreceived: time.Time{},
			Timer:        nil,
		}
	}
}
