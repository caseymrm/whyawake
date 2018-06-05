package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"syscall"

	"github.com/caseymrm/go-pmset"
	"github.com/caseymrm/menuet"
)

var sleepKeywords = map[string]bool{
	"PreventUserIdleDisplaySleep": true,
	"NoDisplaySleepAssertion":     true,
}
var otherChangeKeywords = map[string]bool{
	"PreventUserIdleSystemSleep": true,
}
var canSleepTitle = "💤"
var cantSleepTitle = "😳"
var caffeinatedTitle = "🤪"

func canSleep() bool {
	asserts := pmset.GetAssertions()
	for key, val := range asserts {
		if val == 1 && sleepKeywords[key] {
			return false
		}
	}
	return true
}

func menuItems() []menuet.MenuItem {
	items := make([]menuet.MenuItem, 0)

	if preventingSleep() {
		items = append(items, menuet.MenuItem{
			Text:     preventionRemaining(),
			FontSize: 12,
		})
		items = append(items, menuet.MenuItem{
			Text:     "Deactivate",
			Callback: "deactivate",
		}, menuet.MenuItem{
			Type: menuet.Separator,
		})
	}

	processes := make([]menuet.MenuItem, 0)
	pidAsserts := pmset.GetPIDAssertions()
	for key := range sleepKeywords {
		pids := pidAsserts[key]
		for _, pid := range pids {
			if pid.PID == cafPID {
				continue
			}
			processes = append(processes, menuet.MenuItem{
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
		items = append(items, menuet.MenuItem{
			Text:     text,
			FontSize: 12,
		})
		items = append(items, processes...)
	} else if !preventingSleep() {
		items = append(items, menuet.MenuItem{
			Text: "Your Mac can sleep",
		})
	}

	items = append(items, menuet.MenuItem{
		Type: menuet.Separator,
	})
	items = append(items, menuet.MenuItem{
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
	menuet.App().SetMenuState(&menuet.MenuState{
		Image: image,
		Items: menuItems(),
	})
}

func monitorAssertionChanges(channel chan pmset.AssertionChange) {
	for change := range channel {
		if sleepKeywords[change.Type] || otherChangeKeywords[change.Type] {
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
			response := menuet.App().Alert(menuet.Alert{
				MessageText:     "Kill process?",
				InformativeText: fmt.Sprintf("PID %d", pid),
				Buttons:         []string{"Kill", "Force Kill", "Cancel"},
			})
			switch response.Button {
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
	assertionsChannel := make(chan pmset.AssertionChange)
	clickChannel := make(chan string)
	pmset.SubscribeAssertionChanges(assertionsChannel)
	go monitorAssertionChanges(assertionsChannel)
	setMenuState()
	app := menuet.App()
	app.Name = "Why Awake?"
	app.Label = "com.github.caseymrm.whyawake"
	app.Clicked = clickChannel
	app.MenuOpened = func() []menuet.MenuItem {
		return menuItems()
	}
	app.AutoUpdate.Version = "v0.5"
	app.AutoUpdate.Repo = "caseymrm/whyawake"
	go handleClicks(clickChannel)
	app.RunApplication()
}
