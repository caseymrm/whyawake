package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

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
