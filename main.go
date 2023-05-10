package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type batteryInfo struct {
	Percentage int
	State      string
}

// getBattery retrieves the battery information for the first battery found on the system.
func getBattery() batteryInfo {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		log.Println("Failed to get battery information:", err)
	}

	re := regexp.MustCompile(`[0-9]+%`)
	percentageStr := re.FindString(string(out))
	percentage, err := strconv.Atoi(percentageStr[:len(percentageStr)-1])
	if err != nil {
		log.Println("Failed to get battery percantage:", err)
	}

	var state string
	if regexp.MustCompile(`discharging`).MatchString(string(out)) {
		state = "Discharging"
	} else if regexp.MustCompile(`charging`).MatchString(string(out)) {
		state = "Charging"
	} else {
		state = "Unknown"
	}

	return batteryInfo{
		Percentage: percentage,
		State:      state,
	}
}

// playWhenPercentageOutsideRange plays a sound repeatedly when the battery percentage is
// below or equal to the minimum threshold and discharging, or above or equal to the
// maximum threshold and charging. The sound to play is specified by the 'sound' parameter.
func playWhenPercentageOutsideRange(minPercentage, maxPercentage int, sound string) {
	for {
		batteryInfo := getBattery()
		percentage := batteryInfo.Percentage
		isDischarging := batteryInfo.State == "Discharging"
		isCharging := batteryInfo.State == "Charging"

		if (isDischarging && percentage <= minPercentage) || (isCharging && percentage >= maxPercentage) {
			playSound(sound, percentage, isCharging)
			// Wait 10 seconds, and play sound again
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
}

func main() {
	log.Println("Battery Reminder Started.")

	// Define command-line flags for minimum and maximum battery percentage, and sound file path.
	minPercentage := flag.Int("min", 20, "Minimum battery percentage")
	maxPercentage := flag.Int("max", 80, "Maximum battery percentage")
	sound := flag.String("sound", "", "Sound to play, for example: /System/Library/Sounds/Ping.aiff")

	// Parse the command-line flags.
	flag.Parse()

	for {
		// Play sound when the battery percentage is outside the specified range.
		playWhenPercentageOutsideRange(*minPercentage, *maxPercentage, *sound)

		// Get current battery state
		batteryInfo := getBattery()

		// Calculate sleep time based on how close we are to the thresholds
		var sleepTime time.Duration
		if batteryInfo.State == "Charging" {
			// if charging, and already at maxPercentage. don't wait.
			if batteryInfo.Percentage >= *maxPercentage {
				sleepTime = 0
			}

			// If charging, the closer we are to maxPercentage, the shorter the sleepTime
			sleepTime = time.Duration(*maxPercentage-batteryInfo.Percentage) * time.Second
		} else if batteryInfo.State == "Discharging" {
			// If discharging, the closer we are to minPercentage, the shorter the sleepTime
			sleepTime = time.Duration(batteryInfo.Percentage-*minPercentage) * time.Second
		}

		// Ensure that sleepTime is never less than 1 second
		if sleepTime < 1*time.Second {
			sleepTime = 1 * time.Second
		}

		// Wait for sleepTime seconds before checking the battery percentage again.
		time.Sleep(sleepTime)
	}
}

// unmute sets the system volume to maximum.
func unmute() {
	cmd := exec.Command("osascript", "-e", "set volume output volume 100")
	cmd.Run()
}

// playSound plays the specified sound file.
func playSound(sound string, percentage int, isCharging bool) {
	unmute()

	if sound != "" {
		// Play the sound using the afplay command
		cmd := exec.Command("afplay", sound)
		err := cmd.Run()
		if err != nil {
			log.Println("Failed to play sound:", err)
			speak(percentage, isCharging)
		}
	} else {
		speak(percentage, isCharging)
	}
}

func speak(percentage int, isCharging bool) {
	message := fmt.Sprintf("Battery is at %d%%. Please plugin your charger.", percentage)

	if isCharging {
		message = fmt.Sprintf("Battery is %d%% charged. Please unplug your charger.", percentage)
	}

	exec.Command("say", message).Run()
}
