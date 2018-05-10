package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"

	"github.com/caseymrm/go-assertions"
	"github.com/caseymrm/menuet/tray"
)

var sleepKeywords = map[string]bool{
	"PreventUserIdleDisplaySleep": true,
	//"PreventUserIdleSystemSleep":  true,
	"NoDisplaySleepAssertion": true,
}
var canSleepTitle = "💤"
var cantSleepTitle = "😳"
var caffeinatedTitle = "🤪"

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
	items := make([]tray.MenuItem, 0)

	if preventingSleep() {
		items = append(items, tray.MenuItem{
			Text:     preventionRemaining(),
			FontSize: 12,
		})
		items = append(items, tray.MenuItem{
			Text:     "Deactivate",
			Callback: "deactivate",
		}, tray.MenuItem{
			Text: "---",
		})
	}

	processes := make([]tray.MenuItem, 0)
	pidAsserts := assertions.GetPIDAssertions()
	for key := range sleepKeywords {
		pids := pidAsserts[key]
		for _, pid := range pids {
			if pid.PID == cafPID {
				continue
			}
			processes = append(processes, tray.MenuItem{
				Text:     pid.Name,
				Callback: fmt.Sprintf("pid:%d", pid.PID),
			})
		}
	}
	if len(processes) > 0 {
		text := "1 process keeping your Mac awake"
		if len(processes) > 1 {
			text = fmt.Sprintf("%d processes keeping your Mac awake", len(processes))
		}
		items = append(items, tray.MenuItem{
			Text:     text,
			FontSize: 12,
		})
		items = append(items, processes...)
	} else if !preventingSleep() {
		items = append(items, tray.MenuItem{
			Text: "Your Mac can sleep",
		})
	}

	items = append(items, tray.MenuItem{
		Text: "---",
	})
	items = append(items, tray.MenuItem{
		Text:     "Keep this Mac awake",
		FontSize: 12,
	})
	for _, option := range sleepOptions {
		option.State = sleepOptionSelected(option)
		items = append(items, option)
	}

	return items
}

func setMenuState() {
	image := "Red Eye.pdf"
	if canSleep() {
		image = "Eye.pdf"
	}
	if preventingSleep() {
		image = "Awake Eye.pdf"
	}
	tray.App().SetMenuState(&tray.MenuState{
		Image: image,
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
	case "deactivate":
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
	app.Name = "Why Awake?"
	app.Label = "com.github.caseymrm.whyawake"
	app.Clicked = trayChannel
	app.MenuOpened = func() []tray.MenuItem {
		return menuItems()
	}
	go handleClicks(trayChannel)
	app.RunApplication()
}
