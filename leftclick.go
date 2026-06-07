package main

import "github.com/caseymrm/menuet/v2"

// Persisted as NSUserDefaults integer "leftClickDuration". 0 (the default
// for a missing key) means "left-click opens the menu, like before". The
// other sentinel (-2 = indefinite) shifts preventSleep's 0=indefinite out
// of the way so we can reserve 0 for "off" at the storage layer.
const (
	leftClickDefaultsKey = "leftClickDuration"

	storedOff        = 0
	storedLid        = -1
	storedIndefinite = -2
)

type leftClickChoice struct {
	label  string
	stored int
}

var leftClickChoices = []leftClickChoice{
	{"Off", storedOff},
	{"Until I close the lid", storedLid},
	{"Indefinitely", storedIndefinite},
	{"10 minutes", 10},
	{"30 minutes", 30},
	{"1 hour", 60},
	{"3 hours", 180},
}

// readLeftClickStored returns the raw persisted value.
func readLeftClickStored() int {
	return menuet.Defaults().Integer(leftClickDefaultsKey)
}

// leftClickMinutes converts the persisted sentinel into the integer
// preventSleep expects: lidMode, 0 (indefinite), or N minutes.
// The bool reports whether left-click is wired up at all.
func leftClickMinutes() (minutes int, enabled bool) {
	switch v := readLeftClickStored(); v {
	case storedOff:
		return 0, false
	case storedLid:
		return lidMode, true
	case storedIndefinite:
		return 0, true
	default:
		return v, true
	}
}

func setLeftClickStored(stored int) {
	menuet.Defaults().SetInteger(leftClickDefaultsKey, stored)
	refreshLeftClickHandler()
}

// refreshLeftClickHandler sets or clears menuet's top-level click handler
// to match the current setting. Safe to call from any goroutine — menuet
// reads App().Clicked on each click.
func refreshLeftClickHandler() {
	_, enabled := leftClickMinutes()
	if enabled {
		menuet.App().Clicked = toggleLeftClick
	} else {
		menuet.App().Clicked = nil
	}
}

// toggleLeftClick is what menuet calls on a top-level left click. If sleep
// prevention is already running, this cancels it; otherwise it starts a new
// session at the configured default duration.
func toggleLeftClick() {
	if preventingSleep() {
		cancelSleepPrevention()
		return
	}
	minutes, enabled := leftClickMinutes()
	if !enabled {
		return
	}
	go preventSleep(minutes)
}

// leftClickDefaultLabel produces the parent menu item's text so the current
// choice is visible without opening the submenu.
func leftClickDefaultLabel() string {
	current := readLeftClickStored()
	for _, c := range leftClickChoices {
		if c.stored == current {
			return "Left-click: " + c.label
		}
	}
	return "Left-click: Off"
}

func leftClickDefaultMenu() []menuet.MenuItem {
	current := readLeftClickStored()
	items := make([]menuet.MenuItem, len(leftClickChoices))
	for i, c := range leftClickChoices {
		stored := c.stored
		items[i] = menuet.MenuItem{
			Text:    c.label,
			State:   current == stored,
			Clicked: func() { setLeftClickStored(stored) },
		}
	}
	return items
}
