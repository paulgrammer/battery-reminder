package monitor

import (
	"log"
	"time"

	"github.com/paulgrammer/battery-reminder/internal/battery"
	"github.com/paulgrammer/battery-reminder/internal/config"
	"github.com/paulgrammer/battery-reminder/internal/notifier"
)

// BatteryMonitor coordinates battery monitoring and notifications
type BatteryMonitor struct {
	config         *config.Config
	batteryMonitor *battery.Monitor
	notifier       *notifier.Notifier
}

// NewBatteryMonitor creates a new battery monitor
func NewBatteryMonitor(cfg *config.Config) *BatteryMonitor {
	return &BatteryMonitor{
		config:         cfg,
		batteryMonitor: battery.NewMonitor(),
		notifier:       notifier.NewNotifier(),
	}
}

// Start begins the battery monitoring loop
func (m *BatteryMonitor) Start() {
	log.Println(m.config.String())
	if err := m.notifier.SendNotification(m.config.String()); err != nil {
		log.Printf("Failed to send startup notification: %v", err)
	}

	for {
		m.checkAndAlert()
		sleepTime := m.calculateSleepTime()
		time.Sleep(sleepTime)
	}
}

// checkAndAlert checks battery status and sends alerts if needed
func (m *BatteryMonitor) checkAndAlert() {
	info, err := m.batteryMonitor.GetInfo()
	if err != nil {
		log.Printf("Failed to get battery info: %v", err)
		return
	}

	if m.shouldAlert(info) {
		m.sendAlert(info)
		m.waitAndRepeatAlerts(info)
	}
}

// shouldAlert determines if an alert should be sent based on battery status
func (m *BatteryMonitor) shouldAlert(info battery.Info) bool {
	isLowBattery := info.IsDischarging() && info.Percentage <= m.config.MinPercentage
	isHighBattery := info.IsCharging() && info.Percentage >= m.config.MaxPercentage
	return isLowBattery || isHighBattery
}

// sendAlert sends the appropriate alert based on battery status
func (m *BatteryMonitor) sendAlert(info battery.Info) {
	if info.IsCharging() {
		m.notifier.AlertHighBattery(info.Percentage, m.config.SoundPath)
	} else {
		m.notifier.AlertLowBattery(info.Percentage, m.config.SoundPath)
	}
}

// waitAndRepeatAlerts continues alerting while battery is outside range
func (m *BatteryMonitor) waitAndRepeatAlerts(info battery.Info) {
	for {
		time.Sleep(10 * time.Second)

		info, err := m.batteryMonitor.GetInfo()
		if err != nil {
			log.Printf("Failed to get battery info: %v", err)
			break
		}

		if !m.shouldAlert(info) {
			break
		}

		m.sendAlert(info)
	}
}

// calculateSleepTime determines how long to wait before the next check
func (m *BatteryMonitor) calculateSleepTime() time.Duration {
	info, err := m.batteryMonitor.GetInfo()
	if err != nil {
		log.Printf("Failed to get battery info: %v", err)
		return 60 * time.Second
	}

	var sleepTime time.Duration

	switch info.State {
	case battery.StateCharging:
		if info.Percentage >= m.config.MaxPercentage {
			sleepTime = 1 * time.Second
		} else {
			sleepTime = time.Duration(m.config.MaxPercentage-info.Percentage) * time.Second
		}
	case battery.StateDischarging:
		sleepTime = time.Duration(info.Percentage-m.config.MinPercentage) * time.Second
	default:
		sleepTime = 60 * time.Second
	}

	if sleepTime < 1*time.Second {
		sleepTime = 1 * time.Second
	}

	return sleepTime
}
