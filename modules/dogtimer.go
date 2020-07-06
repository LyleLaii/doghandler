package modules

import (
	"doghandler/pkg/logger"
	"fmt"
	"sync"
	"time"
)

// Dogs dog map
var Dogs = make(map[string]*Dog)

// Notifier Notifier interface
type Notifier interface {
	SendMessage(*Dog) error
}

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
	Notifiers    []Notifier
	mu           sync.Mutex
}

// TouchDog refresh dog info
func (d *Dog) TouchDog() {
	d.mu.Lock()
	d.Counter = 0
	d.Lastreceived = time.Now()
	d.refreshtimer()
	d.mu.Unlock()
	logger.LogInfo("DogTimer", fmt.Sprintf("service %s last received time: %s", d.ServiceID, d.Lastreceived))
}

func (d *Dog) refreshtimer() {
	if d.Timer != nil {
		d.Timer.Stop()
	}
	d.Timer = time.AfterFunc(d.Interval, d.CheckDog)
}

// CheckDog check dog status
func (d *Dog) CheckDog() {
	d.mu.Lock()
	d.Counter = d.Counter + 1
	logger.LogInfo("DogTimer", fmt.Sprintf("service %s did not receiver, now counter is : %v, max count is: %v", d.ServiceID, d.Counter, d.Maxcount))
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
	logger.LogInfo("DogTimer", fmt.Sprintf("service %s maximum number of not received, send alert", d.ServiceID))
	for _, n := range d.Notifiers {
		go n.SendMessage(d)
	}
	d.Counter = 0
	d.refreshtimer()
}
