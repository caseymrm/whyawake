package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/caseymrm/go-caffeinate"
	"github.com/caseymrm/go-statusbar/tray"
)

var sleepOptions = []tray.MenuItem{
	{Text: "Indefinitely", Callback: "prevent:0"},
	{Text: "1 min testing", Callback: "prevent:1"},
	{Text: "5 minutes", Callback: "prevent:5"},
	{Text: "10 minutes", Callback: "prevent:10"},
	{Text: "15 minutes", Callback: "prevent:15"},
	{Text: "30 minutes", Callback: "prevent:30"},
	{Text: "1 hour", Callback: "prevent:60"},
	{Text: "2 hours", Callback: "prevent:120"},
	{Text: "5 hours", Callback: "prevent:300"},
}

var caf caffeinate.Caffeinate
var cafOnce sync.Once

var cafPID int
var cafExpire time.Time
var cafMinutes int

func preventSleep(minutes int) {
	log.Printf("PreventSleep %d", minutes)
	cafOnce.Do(func() {
		caf = caffeinate.Caffeinate{
			Display:    true,
			IdleSystem: true,
			PID:        os.Getpid(),
		}
	})
	caf.Timeout = 60 * minutes
	log.Printf("Starting")
	caf.Start()
	log.Printf("Done starting")
	cafMinutes = minutes
	if minutes > 0 {
		cafExpire = time.Now().Add(time.Duration(minutes) * time.Minute)
	} else {
		cafExpire = time.Time{}
	}
	cafPID = caf.CaffeinatePID()
	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
		for range ticker.C {
			setMenuState()
		}
	}()
	log.Printf("Waiting... %d", cafPID)
	caf.Wait()
	log.Printf("Done caffeinating... %d", cafPID)
	ticker.Stop()
	cafPID = 0
	cafExpire = time.Time{}
	setMenuState()
}

func cancelSleepPrevention() {
	caf.Stop()
}

func preventionRemaining() string {
	if cafPID == 0 {
		return ""
	}
	if cafMinutes == 0 {
		return "∞️"
	}
	remaining := int(time.Until(cafExpire).Seconds())
	if remaining > 60 {
		return fmt.Sprintf("%dm", remaining/60)
	}
	return fmt.Sprintf("%ds", remaining)
}
