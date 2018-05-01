package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/caseymrm/go-assertions"
	"github.com/caseymrm/go-statusbar/tray"
)

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

func getStartupPath() string {
	u, err := user.Current()
	if err != nil {
		log.Printf("user.Current: %v", err)
		return ""
	}
	return u.HomeDir + "/Library/LaunchAgents/com.github.caseymrm.CantSleep.plist"
}

func runningAtStartup() bool {
	_, err := os.Stat(getStartupPath())
	if err == nil {
		return true
	}
	return false
}

func removeStartupItem() {
	err := os.Remove(getStartupPath())
	if err != nil {
		log.Printf("os.Remove: %v", err)
	}
}

func addStartupItem() {
	path := getStartupPath()
	// Make sure ~/Library/LaunchAgents exists
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		log.Printf("os.MkdirAll: %v", err)
		return
	}
	executable, err := os.Executable()
	if err != nil {
		log.Printf("os.Executable: %v", err)
		return
	}
	fmt.Println(executable)
	f, err := os.Create(path)
	if err != nil {
		log.Printf("os.Create: %v", err)
		return
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf(launchdConfig, executable))
	if err != nil {
		log.Printf("f.WriteString: %v", err)
		return
	}
}

func menuItems() []tray.MenuItem {
	items := make([]tray.MenuItem, 0)
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
	startupItem := tray.MenuItem{Text: "Run at start up", Callback: "startup"}
	if runningAtStartup() {
		startupItem.State = true
	}
	items = append(items, startupItem)

	return append(preAmble, items...)
}

func menuState() *tray.MenuState {
	if canSleep() {
		return &tray.MenuState{
			Title: canSleepTitle,
		}
	}
	return &tray.MenuState{
		Title: cantSleepTitle,
	}
}

func monitorAssertionChanges(channel chan assertions.AssertionChange) {
	for change := range channel {
		if sleepKeywords[change.Type] {
			tray.App().SetMenuState(menuState())
		}
	}
}

func handleClicks(callback chan string) {
	for clicked := range callback {
		switch clicked {
		case "startup":
			if runningAtStartup() {
				removeStartupItem()
			} else {
				addStartupItem()
			}
		default:
			pid, _ := strconv.Atoi(clicked)
			go func() {
				switch tray.App().Alert("Kill process?", fmt.Sprintf("PID %d", pid), "Kill", "Kill -9", "Cancel") {
				case 0:
					fmt.Printf("Killing pid %d\n", pid)
					syscall.Kill(pid, syscall.SIGTERM)
				case 1:
					fmt.Printf("Killing -9 pid %d\n", pid)
					syscall.Kill(pid, syscall.SIGKILL)
				}
			}()
		}
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
		return menuItems()
	}
	go handleClicks(trayChannel)
	app.RunApplication()
}

var launchdConfig = `
<?xml version='1.0' encoding='UTF-8'?>
 <!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
 <plist version='1.0'>
   <dict>
     <key>Label</key><string>CantSleep</string>
     <key>Program</key><string>%s</string>
     <key>StandardOutPath</key><string>/tmp/CantSleep-out.log</string>
     <key>StandardErrorPath</key><string>/tmp/CantSleep-err.log</string>
     <key>KeepAlive</key><true/>
     <key>RunAtLoad</key><true/>
   </dict>
</plist>
`
