package main

import (
	"flag"
	"log"

	"github.com/paulgrammer/battery-reminder/internal/config"
	"github.com/paulgrammer/battery-reminder/internal/monitor"
	"github.com/paulgrammer/battery-reminder/internal/svc"
)

func main() {
	minPercentage := flag.Int("min", 20, "Minimum battery percentage")
	maxPercentage := flag.Int("max", 80, "Maximum battery percentage")
	soundPath := flag.String("sound", "", "Sound to play, for example: /System/Library/Sounds/Ping.aiff")
	serviceFlag := flag.String("service", "", "Control the system service (install, uninstall, start, stop, restart, run)")
	flag.Parse()

	cfg, err := config.NewConfig(*minPercentage, *maxPercentage, *soundPath)
	if err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	if *serviceFlag != "" {
		handleServiceCommand(cfg, *serviceFlag)
		return
	}

	runMonitor(cfg)
}

func handleServiceCommand(cfg *config.Config, command string) {
	sm, err := svc.NewServiceManager(cfg)
	if err != nil {
		log.Fatalf("Failed to create service manager: %v", err)
	}

	if err := sm.Control(command); err != nil {
		log.Fatalf("Service command '%s' failed: %v", command, err)
	}
}

func runMonitor(cfg *config.Config) {
	m := monitor.NewBatteryMonitor(cfg)
	m.Start()
}
