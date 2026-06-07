package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/caseymrm/go-caffeinate"
	"github.com/caseymrm/menuet/v2"
)

// lidMode is a sentinel value for cafMinutes meaning "keep awake until the
// lid is closed" — no time-based timeout, the clamshell watcher cancels.
const lidMode = -1

func sleepOptions() []menuet.MenuItem {
	return []menuet.MenuItem{
		{Text: "Until I close the lid", Clicked: func() { preventSleep(lidMode) }, State: sleepOptionSelected(lidMode)},
		{Text: "Indefinitely", Clicked: func() { preventSleep(0) }, State: sleepOptionSelected(0)},
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
	if minutes > 0 {
		caf.Timeout = 60 * minutes
	} else {
		caf.Timeout = 0
	}
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
			Subtitle: startSubtitle(minutes),
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
		Subtitle: stopSubtitle(minutes),
		Message:  "Keep preventing sleep in the Why Awake? menu",
	})
}

func startSubtitle(minutes int) string {
	switch {
	case minutes == lidMode:
		return "Your computer will stay awake until you close the lid"
	case minutes == 0:
		return "Your computer will not sleep until you deactivate"
	default:
		return fmt.Sprintf("Your computer will not sleep for %d minutes", minutes)
	}
}

func stopSubtitle(minutes int) string {
	if minutes == lidMode {
		return "The lid was closed"
	}
	return fmt.Sprintf("It was kept awake for %d minutes", minutes)
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
	switch {
	case cafMinutes == lidMode:
		return "Staying awake until lid closes"
	case cafMinutes == 0:
		return "Staying awake indefinitely"
	}
	remaining := int(time.Until(cafExpire).Seconds())
	if remaining > 60 {
		return fmt.Sprintf("%d minutes remaining", remaining/60)
	}
	return fmt.Sprintf("%d seconds remaining", remaining)
}
