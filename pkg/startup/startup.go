package startup

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"text/template"
)

type Startup struct {
	AppLabel string
	AppName  string

	launchdOnce     sync.Once
	launchdTemplate *template.Template
}

const launchdString = `
<?xml version='1.0' encoding='UTF-8'?>
 <!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\" >
 <plist version='1.0'>
   <dict>
     <key>Label</key>
     <string>{{.Name}}</string>
     <key>Program</key>
     <string>{{.Executable}}</string>
     <key>StandardOutPath</key>
     <string>/tmp/{{.Label}}-out.log</string>
     <key>StandardErrorPath</key>
     <string>/tmp/{{.Label}}-err.log</string>
     <key>RunAtLoad</key>
     <true/>
   </dict>
</plist>
`

func (s *Startup) getStartupPath() string {
	if s.AppLabel == "" {
		log.Fatal("Need to set a Label for the app")
	}
	u, err := user.Current()
	if err != nil {
		log.Printf("user.Current: %v", err)
		return ""
	}
	return u.HomeDir + "/Library/LaunchAgents/" + s.AppLabel + ".plist"
}

func (s *Startup) RunningAtStartup() bool {
	if s.AppLabel == "" {
		log.Println("Warning: no application Label set")
		return false
	}
	_, err := os.Stat(s.getStartupPath())
	return err == nil
}

func (s *Startup) RemoveStartupItem() {
	err := os.Remove(s.getStartupPath())
	if err != nil {
		log.Printf("os.Remove: %v", err)
	}
}

func (s *Startup) AddStartupItem() {
	path := s.getStartupPath()
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
	f, err := os.Create(path)
	if err != nil {
		log.Printf("os.Create: %v", err)
		return
	}
	defer f.Close()
	s.launchdOnce.Do(func() {
		fmt.Println("called once")
		s.launchdTemplate = template.Must(template.New("launchdConfig").Parse(launchdString))
	})
	err = s.launchdTemplate.Execute(f,
		struct {
			Name       string
			Label      string
			Executable string
		}{
			s.AppName,
			s.AppLabel,
			executable,
		})
	if err != nil {
		log.Printf("template.Execute: %v", err)
		return
	}
}
