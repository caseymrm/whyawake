package main

import (
	"fmt"
	"syscall"

	"github.com/caseymrm/go-clamshell"
	"github.com/caseymrm/go-pmset"
	"github.com/caseymrm/menuet/v2"
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
		items = append(items, menuet.Regular{
			Text:     preventionRemaining(),
			FontSize: 12,
		})
		items = append(items,
			menuet.Regular{
				Text:    "Deactivate",
				Clicked: cancelSleepPrevention,
			},
			menuet.Separator{},
		)
	}

	processes := make([]menuet.MenuItem, 0)
	pidAsserts := pmset.GetPIDAssertions()
	for key := range sleepKeywords {
		pids := pidAsserts[key]
		for _, pid := range pids {
			if pid.PID == cafPID {
				continue
			}
			processes = append(processes, menuet.Regular{
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
		items = append(items, menuet.Regular{
			Text:     text,
			FontSize: 12,
		})
		items = append(items, processes...)
	} else if !preventingSleep() {
		items = append(items, menuet.Regular{
			Text: "Your Mac can sleep",
		})
	}

	items = append(items, menuet.Separator{})
	items = append(items, menuet.Regular{
		Text:     "Keep this Mac awake",
		FontSize: 12,
	})
	items = append(items, sleepOptions()...)

	items = append(items, menuet.Separator{})
	items = append(items, menuet.Regular{
		Text:     leftClickDefaultLabel(),
		FontSize: 12,
		Children: leftClickDefaultMenu,
	})

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

func monitorClamshell(channel chan bool) {
	for closed := range channel {
		if closed && preventingSleep() && cafMinutes == lidMode {
			cancelSleepPrevention()
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
	clamshellChannel := make(chan bool, 4)
	clamshell.SubscribeClamshellChanges(clamshellChannel)
	go monitorClamshell(clamshellChannel)
	setMenuState()
	app := menuet.App()
	app.Name = "Why Awake?"
	app.Label = "com.github.caseymrm.whyawake"
	app.Children = menuItems
	app.AutoUpdate.Version = "v0.9"
	app.AutoUpdate.Repo = "caseymrm/whyawake"
	refreshLeftClickHandler()
	app.RunApplication()
}
