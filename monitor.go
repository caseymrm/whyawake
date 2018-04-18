package main

import (
	"fmt"
	"regexp"
	"time"

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

func pmset() *tray.MenuState {
	asserts := assertions.GetAssertions()
	pidAsserts := assertions.GetPIDAssertions()
	canSleep := true
	for key, val := range asserts {
		if val == 1 && sleepKeywords[key] {
			canSleep = false
		}
	}
	ms := tray.MenuState{
		Items: make([]tray.MenuItem, 0, 1),
	}
	for key := range sleepKeywords {
		pids := pidAsserts[key]
		for _, pid := range pids {
			ms.Items = append(ms.Items, tray.MenuItem{
				Text:     fmt.Sprintf("%s (pid %d)", pid.Name, pid.PID),
				Callback: fmt.Sprintf("%d", pid.PID),
			})
		}
	}
	if canSleep {
		ms.Title = canSleepTitle
	} else {
		ms.Title = cantSleepTitle
	}
	preAmble := []tray.MenuItem{{Text: "Your laptop can sleep!"}}
	if !canSleep {
		if len(ms.Items) == 1 {
			preAmble = []tray.MenuItem{{Text: "1 process is keeping your laptop awake:"}}
		} else {
			preAmble = []tray.MenuItem{{Text: fmt.Sprintf("%d processes are keeping your laptop awake:", len(ms.Items))}}
		}
	}
	if len(ms.Items) > 0 {
		preAmble = append(preAmble, tray.MenuItem{Text: "---"})
	}
	ms.Items = append(preAmble, ms.Items...)
	return &ms
}

func monitorPmSet() {
	for {
		tray.App().SetMenuState(pmset())
		time.Sleep(10 * time.Second)
	}
}

func handleClicks(callback chan string) {
	for pid := range callback {
		fmt.Printf("PID Clicked %s\n", pid)
	}
}

func main() {
	go monitorPmSet()
	callback := make(chan string)
	app := tray.App()
	app.Clicked = callback
	app.MenuOpened = func() []tray.MenuItem {
		ms := pmset()
		return ms.Items
	}
	go handleClicks(callback)
	app.RunApplication()
}
