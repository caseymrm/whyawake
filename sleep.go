package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caseymrm/go-caffeinate"
	"github.com/caseymrm/go-statusbar/tray"
)

var sleepOptions = []tray.MenuItem{
	{Text: "Indefinitely", Callback: "prevent:0"},
	//{Text: "1 min test", Callback: "prevent:1"},
	{Text: "10 minutes", Callback: "prevent:10"},
	{Text: "30 minutes", Callback: "prevent:30"},
	{Text: "1 hour", Callback: "prevent:60"},
	{Text: "3 hours", Callback: "prevent:180"},
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
	if cafMinutes > 0 {
		go func() {
			for range ticker.C {
				setMenuState()
			}
		}()
	}
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

func preventingSleep() bool {
	return cafPID != 0
}

func sleepOptionSelected(item tray.MenuItem) bool {
	if cafPID == 0 {
		return false
	}
	return strings.HasSuffix(item.Callback, fmt.Sprintf(":%d", cafMinutes))
}

func preventionRemaining() string {
	if cafMinutes == 0 {
		return "Staying awake indefinitely"
	}
	remaining := int(time.Until(cafExpire).Seconds())
	if remaining > 60 {
		return fmt.Sprintf("%d minutes remaining", remaining/60)
	}
	return fmt.Sprintf("%d seconds remaining", remaining)
}
