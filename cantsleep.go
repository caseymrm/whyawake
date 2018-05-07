package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"

	"github.com/caseymrm/go-assertions"
	"github.com/caseymrm/go-statusbar/tray"
)

var sleepKeywords = map[string]bool{
	"PreventUserIdleDisplaySleep": true,
	//"PreventUserIdleSystemSleep":  true,
	"NoDisplaySleepAssertion": true,
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

func menuItems() []tray.MenuItem {
	processes := make([]tray.MenuItem, 0)
	pidAsserts := assertions.GetPIDAssertions()
	preventingSleep := false
	for key := range sleepKeywords {
		pids := pidAsserts[key]
		for _, pid := range pids {
			if pid.PID == cafPID {
				preventingSleep = true
				processes = append(processes, tray.MenuItem{
					Text:     "You have prevented sleep for " + preventionRemaining() + " longer",
					Callback: fmt.Sprintf("pid:%d", pid.PID),
				})
				continue
			}
			processes = append(processes, tray.MenuItem{
				Text:     pid.Name,
				Callback: fmt.Sprintf("pid:%d", pid.PID),
			})
		}
	}
	items := []tray.MenuItem{{Text: "Your laptop can sleep!"}}
	if len(processes) > 0 {
		text := "1 process is keeping your laptop awake"
		if len(processes) > 1 {
			text = fmt.Sprintf("%d processes are keeping your laptop awake", len(processes))
		}
		items = []tray.MenuItem{{
			Text:     text,
			Children: processes,
		}}
	}
	items = append(items, tray.MenuItem{
		Text:     "Prevent sleep",
		Children: sleepOptions,
		Callback: "prevent_sleep",
		State:    preventingSleep,
	})
	startupItem := tray.MenuItem{Text: "Run at start up", Callback: "startup"}
	if runningAtStartup() {
		startupItem.State = true
	}
	return append(items, startupItem)
}

func setMenuState() {
	title := cantSleepTitle
	if canSleep() {
		title = canSleepTitle
	}
	title += preventionRemaining()
	tray.App().SetMenuState(&tray.MenuState{
		Title: title,
		Items: menuItems(),
	})
}

func monitorAssertionChanges(channel chan assertions.AssertionChange) {
	for change := range channel {
		if sleepKeywords[change.Type] {
			setMenuState()
		}
	}
}

func handleClicks(callback chan string) {
	for clicked := range callback {
		go handleClick(clicked)
	}
}

func handleClick(clicked string) {
	switch clicked {
	case "startup":
		if runningAtStartup() {
			removeStartupItem()
		} else {
			addStartupItem()
		}
	case "prevent_sleep":
		cancelSleepPrevention()
	default:
		if strings.HasPrefix(clicked, "pid:") {
			pid, _ := strconv.Atoi(clicked[4:])
			switch tray.App().Alert("Kill process?", fmt.Sprintf("PID %d", pid), "Kill", "Force Kill", "Cancel") {
			case 0:
				fmt.Printf("Killing pid %d\n", pid)
				syscall.Kill(pid, syscall.SIGTERM)
			case 1:
				fmt.Printf("Killing -9 pid %d\n", pid)
				syscall.Kill(pid, syscall.SIGKILL)
			}
			return
		}
		if strings.HasPrefix(clicked, "prevent:") {
			minutes, _ := strconv.Atoi(clicked[8:])
			preventSleep(minutes)
			return
		}
		log.Printf("Other: %s", clicked)
	}
}

func main() {
	assertionsChannel := make(chan assertions.AssertionChange)
	trayChannel := make(chan string)
	assertions.SubscribeAssertionChanges(assertionsChannel)
	go monitorAssertionChanges(assertionsChannel)
	setMenuState()
	app := tray.App()
	app.Clicked = trayChannel
	app.MenuOpened = func() []tray.MenuItem {
		return menuItems()
	}
	go handleClicks(trayChannel)
	app.RunApplication()
}
