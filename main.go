package main

import (
	"log"
	"math"
	"os/exec"
	"time"

	"github.com/distatus/battery"
)

func main() {
	for {
		batteryInfo, err := battery.Get(0)
		if err != nil {
			log.Fatal("Failed to get battery information:", err)
		}

		parcentage := int(math.Round(batteryInfo.Current / batteryInfo.Full * 100))

		// Check if the battery percentage is 21% and not charging
		if batteryInfo.State.String() == "Discharging" && parcentage <= 21 {
			playSound()
		}

		// Check if the battery percentage is 79% and still charging
		if batteryInfo.State.String() == "Charging" && parcentage >= 79 {
			playSound()
		}

		// Wait for 5 minutes before checking the battery percentage again
		time.Sleep(5 * time.Minute)
	}
}

func unmute() {
	cmd := exec.Command("osascript", "-e", "set volume output volume 100")
	cmd.Run()
}

func playSound() {
	unmute()

	// Play the sound using the afplay command
	cmd := exec.Command("afplay", "/System/Library/Sounds/Ping.aiff")
	err := cmd.Run()
	if err != nil {
		log.Fatal("Failed to play sound:", err)
	}
}
