package main

import (
    "encoding/json"
    "fmt"
    "os"
    "sync"
    "time"
	"os/exec"
)

type ServiceConfig struct {
    Name          string `json:"name"`
    CheckURL      string `json:"check_url"`
    StartTimeout  int    `json:"start_timeout"`
    RetryAttempts int    `json:"retry_attempts"`
    UseSudo       bool   `json:"use_sudo"`
    StartCmd      string `json:"start_cmd"`
    StopCmd       string `json:"stop_cmd"`
    StopAfter     string `json:"stop_after"`
}

type Service struct {
    Config     ServiceConfig
    IsRunning  bool
    IsStarting bool
	IsDown     bool
    mu         sync.Mutex
}

var serviceMap = map[string]*Service{}
var globalMu sync.Mutex



func LoadServiceConfigs(path string) error {
    file, err := os.Open(path)
    if err != nil {
        return err
    }
    defer file.Close()

    var configs []ServiceConfig
    if err := json.NewDecoder(file).Decode(&configs); err != nil {
        return err
    }

    globalMu.Lock()
    defer globalMu.Unlock()
    for _, cfg := range configs {
        serviceMap[cfg.Name] = &Service{Config: cfg}
    }

    return nil
}


func GetService(name string) (*Service, bool) {
    globalMu.Lock()
    defer globalMu.Unlock()
    svc, ok := serviceMap[name]
    return svc, ok
}



func StartServiceIfNeeded(svc *Service) {
	//fmt.Println("heere22")
    svc.mu.Lock()
    if svc.IsStarting {
        svc.mu.Unlock()
        return
    }

    svc.IsStarting = true
    svc.mu.Unlock()

    go func() {
        //fmt.Printf("Starting service %s...\n", svc.Config.Name)
        //time.Sleep(10 * time.Second) 
		cmd := exec.Command("systemctl", "start","jellyfin.service")

		err := cmd.Run()
    	if err != nil {
        	fmt.Println("Error:", err)
        	return
    	}
        time.Sleep(10 * time.Second)
        svc.mu.Lock()
        svc.IsStarting = false
        svc.mu.Unlock()

        fmt.Printf("Service %s is now running.\n", svc.Config.Name)
    }()



}
