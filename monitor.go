package main

import (
	"fmt"
	"regexp"

	"github.com/caseymrm/go-assertions"
	"github.com/caseymrm/go-statusbar/tray"
)

var processRe = regexp.MustCompile(`pid (\d+)\(([^)]+)\): \[0x[0-9a-f]+\] (\d\d:\d\d:\d\d) (\w+) named: "(.+)"`)
var sleepKeywords = map[string]bool{
	"PreventUserIdleDisplaySleep": true,
	//"PreventUserIdleSystemSleep":  true,
}
var canSleepTitle = "😴"
var cantSleepTitle = "😫"

func canSleep() bool {
	asserts := assertions.GetAssertions()
	for key, val := range asserts {
		if val == 1 && sleepKeywords[key] {
			return false
		}
	}
	return true
}

func wakingItems() []tray.MenuItem {
	items := make([]tray.MenuItem, 0, 1)
	pidAsserts := assertions.GetPIDAssertions()
	for key := range sleepKeywords {
		pids := pidAsserts[key]
		for _, pid := range pids {
			items = append(items, tray.MenuItem{
				Text:     fmt.Sprintf("%s (pid %d)", pid.Name, pid.PID),
				Callback: fmt.Sprintf("%d", pid.PID),
			})
		}
	}
	preAmble := []tray.MenuItem{{Text: "Your laptop can sleep!"}}
	if len(items) == 1 {
		preAmble = []tray.MenuItem{{Text: "1 process is keeping your laptop awake:"}}
	} else if len(items) > 1 {
		preAmble = []tray.MenuItem{{Text: fmt.Sprintf("%d processes are keeping your laptop awake:", len(items))}}
	}
	if len(items) > 0 {
		preAmble = append(preAmble, tray.MenuItem{Text: "---"})
	}
	return append(preAmble, items...)
}

func menuState() *tray.MenuState {
	canSleep := canSleep()
	ms := tray.MenuState{
		Items: wakingItems(),
	}
	if canSleep {
		ms.Title = canSleepTitle
	} else {
		ms.Title = cantSleepTitle
	}
	return &ms
}

func monitorAssertionChanges(channel chan assertions.AssertionChange) {
	for range channel {
		tray.App().SetMenuState(menuState())
	}
}

func handleClicks(callback chan string) {
	for pid := range callback {
		fmt.Printf("PID Clicked %s\n", pid)
	}
}

func main() {
	assertionsChannel := make(chan assertions.AssertionChange)
	trayChannel := make(chan string)
	assertions.SubscribeAssertionChanges(assertionsChannel)
	go monitorAssertionChanges(assertionsChannel)
	app := tray.App()
	app.SetMenuState(menuState())
	app.Clicked = trayChannel
	app.MenuOpened = func() []tray.MenuItem {
		return wakingItems()
	}
	go handleClicks(trayChannel)
	app.RunApplication()
}
