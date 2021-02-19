package dogtimer

import (
	"doghandler/config"
	"doghandler/notifiers"
	"doghandler/pkg/logger"
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
	Interval     time.Duration
	Maxcount     int
	Counter      int
	Lastreceived time.Time
	Timer        *time.Timer
	Receiver     string
	mu           sync.Mutex
}

// TouchDog refresh dog info
func (d *Dog) TouchDog() {
	d.mu.Lock()
	d.Counter = 0
	d.Lastreceived = time.Now()
	d.refreshTimer()
	d.mu.Unlock()
	logger.LogInfo("DogTimer", fmt.Sprintf("service %s received message, time: %s", d.ServiceID, d.Lastreceived))
}

func (d *Dog) refreshTimer() {
	if d.Timer != nil {
		d.Timer.Stop()
	}
	d.Timer = time.AfterFunc(d.Interval, d.CheckDog)
}

// CheckDog check dog status
func (d *Dog) CheckDog() {
	d.mu.Lock()
	d.Counter = d.Counter + 1
	logger.LogInfo("DogTimer", fmt.Sprintf("service %s receive message timeout, last received time: %s,  now counter is : %v, max count is: %v",
		d.ServiceID,
		d.Lastreceived,
		d.Counter,
		d.Maxcount))
	if d.Counter >= d.Maxcount {
		d.Alert()
		d.refreshTimer()
	} else {
		d.refreshTimer()
	}
	d.mu.Unlock()
}

// Alert Send the message that the dog is dead
func (d *Dog) Alert() {
	logger.LogInfo("DogTimer", fmt.Sprintf("service %s maximum number of not received, send alert", d.ServiceID))
	for _, n := range notifiers.NotifyMap[d.Receiver] {
		go n.Notify(&config.Message{
			Name:         d.Name,
			ServiceID:    d.ServiceID,
			Description:  d.Description,
			Lastreceived: d.Lastreceived,
		})
	}
	d.Counter = 0
	d.refreshTimer()
}

// InitDogs initiate Dogs slices
func InitDogs(global GlobalConf, services []ServiceConf) {
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
			Receiver:     receiver,
			Counter:      0,
			Lastreceived: time.Time{},
			Timer:        nil,
		}
	}
	logger.LogDebug("test", fmt.Sprintf("%+v", Dogs))
}

func ClearDogs() {
	logger.LogInfo("dogtimer", "stop all timer to clear dogs")
	for _, v := range Dogs {
		if v.Timer != nil {
			v.Timer.Stop()
		}
	}
	Dogs = make(map[string]*Dog)
}