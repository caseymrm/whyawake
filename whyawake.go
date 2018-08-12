package main

import (
	"fmt"
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

func menuItems(item menuet.MenuItem) []menuet.MenuItem {
	items := make([]menuet.MenuItem, 0)

	if preventingSleep() {
		items = append(items, menuet.MenuItem{
			Text:     preventionRemaining(),
			FontSize: 12,
		})
		items = append(items, menuet.MenuItem{
			Text:    "Deactivate",
			Clicked: cancelSleepPrevention,
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
				Text: pid.Name,
				Clicked: func() {
					killProcess(pid.PID)
				},
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
	items = append(items, sleepOptions()...)

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
	})
	menuet.App().MenuChanged()
}

func monitorAssertionChanges(channel chan pmset.AssertionChange) {
	for change := range channel {
		if sleepKeywords[change.Type] || otherChangeKeywords[change.Type] {
			setMenuState()
		}
	}
}

func killProcess(pid int) {
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
}

func main() {
	assertionsChannel := make(chan pmset.AssertionChange)
	pmset.SubscribeAssertionChanges(assertionsChannel)
	go monitorAssertionChanges(assertionsChannel)
	setMenuState()
	app := menuet.App()
	app.Name = "Why Awake?"
	app.Label = "com.github.caseymrm.whyawake"
	app.Children = menuItems
	app.AutoUpdate.Version = "v0.5"
	app.AutoUpdate.Repo = "caseymrm/whyawake"
	app.RunApplication()
}
