package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/caseymrm/go-caffeinate"
	"github.com/caseymrm/menuet"
)

func sleepOptions() []menuet.MenuItem {
	return []menuet.MenuItem{
		{Text: "Indefinitely", Clicked: func() { preventSleep(0) }, State: sleepOptionSelected(0)},
		//{Text: "1 min (testing)", Clicked: func() {preventSleep(1) }, State: sleepOptionSelected(1)},
		{Text: "10 minutes", Clicked: func() { preventSleep(10) }, State: sleepOptionSelected(10)},
		{Text: "30 minutes", Clicked: func() { preventSleep(30) }, State: sleepOptionSelected(30)},
		{Text: "1 hour", Clicked: func() { preventSleep(60) }, State: sleepOptionSelected(60)},
		{Text: "3 hours", Clicked: func() { preventSleep(180) }, State: sleepOptionSelected(180)},
	}
}

var caf caffeinate.Caffeinate
var cafOnce sync.Once

var cafPID int
var cafExpire time.Time
var cafMinutes int

func preventSleep(minutes int) {
	cafOnce.Do(func() {
		caf = caffeinate.Caffeinate{
			Display:    true,
			IdleSystem: true,
			PID:        os.Getpid(),
		}
	})
	caf.Timeout = 60 * minutes
	caf.Start()
	cafMinutes = minutes
	if minutes > 0 {
		cafExpire = time.Now().Add(time.Duration(minutes) * time.Minute)
	} else {
		cafExpire = time.Time{}
	}
	if cafPID == 0 {
		menuet.App().Notification(menuet.Notification{
			Title:    "Preventing sleep",
			Subtitle: fmt.Sprintf("Your computer will not sleep for %d minutes", minutes),
			Message:  "Deactivate in the Why Awake? menu",
		})
	}
	cafPID = caf.CaffeinatePID()
	caf.Wait()
	cafPID = 0
	cafExpire = time.Time{}
	setMenuState()
	menuet.App().Notification(menuet.Notification{
		Title:    "Your computer can sleep again",
		Subtitle: fmt.Sprintf("It was kept awake for %d minutes", minutes),
		Message:  "Keep preventing sleep in the Why Awake? menu",
	})
}

func cancelSleepPrevention() {
	caf.Stop()
}

func preventingSleep() bool {
	return cafPID != 0
}

func sleepOptionSelected(minutes int) bool {
	if cafPID == 0 {
		return false
	}
	return minutes == cafMinutes
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
