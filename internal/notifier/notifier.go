package notifier

import (
	"fmt"
	"log"
	"os/exec"
	"time"
)

// Notifier handles all notification methods
type Notifier struct{}

// NewNotifier creates a new notifier instance
func NewNotifier() *Notifier {
	return &Notifier{}
}

// SendNotification displays a macOS notification
func (n *Notifier) SendNotification(message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" sound name "default"`, message, "Battery Reminder")
	cmd := exec.Command("osascript", "-e", script)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}
	return nil
}

// Speak uses text-to-speech to read the message
func (n *Notifier) Speak(message string) error {
	cmd := exec.Command("say", message)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to speak: %w", err)
	}
	return nil
}

// Unmute sets the system volume to maximum
func (n *Notifier) Unmute() error {
	cmd := exec.Command("osascript", "-e", "set volume output volume 100")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unmute: %w", err)
	}
	return nil
}

// PlaySound plays the specified sound file
func (n *Notifier) PlaySound(soundPath string) error {
	if soundPath == "" {
		return fmt.Errorf("no sound path provided")
	}

	cmd := exec.Command("afplay", soundPath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to play sound: %w", err)
	}
	return nil
}

// AlertLowBattery sends alerts for low battery
func (n *Notifier) AlertLowBattery(percentage int, soundPath string) {
	n.Unmute()

	message := fmt.Sprintf("Battery is at %d%%. Please plugin your charger.", percentage)

	if err := n.SendNotification(message); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}

	if soundPath != "" {
		if err := n.PlaySound(soundPath); err != nil {
			log.Printf("Sound playback failed, using speech: %v", err)
			n.speakWithDelay(message)
		}
	} else {
		n.speakWithDelay(message)
	}
}

// AlertHighBattery sends alerts for high battery (fully charged)
func (n *Notifier) AlertHighBattery(percentage int, soundPath string) {
	n.Unmute()

	message := fmt.Sprintf("Battery is %d%% charged. Please unplug your charger.", percentage)

	if err := n.SendNotification(message); err != nil {
		log.Printf("Failed to send notification: %v", err)
	}

	if soundPath != "" {
		if err := n.PlaySound(soundPath); err != nil {
			log.Printf("Sound playback failed, using speech: %v", err)
			n.speakWithDelay(message)
		}
	} else {
		n.speakWithDelay(message)
	}
}

// speakWithDelay waits briefly before speaking to avoid overlap with notification
func (n *Notifier) speakWithDelay(message string) {
	time.Sleep(2 * time.Second)
	if err := n.Speak(message); err != nil {
		log.Printf("Failed to speak: %v", err)
	}
}
