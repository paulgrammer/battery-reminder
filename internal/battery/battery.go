package battery

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

// State represents the battery charging state
type State string

const (
	StateCharging    State = "Charging"
	StateDischarging State = "Discharging"
	StateUnknown     State = "Unknown"
)

// Info contains battery information
type Info struct {
	Percentage int
	State      State
}

// Monitor provides battery information retrieval
type Monitor struct{}

// NewMonitor creates a new battery monitor
func NewMonitor() *Monitor {
	return &Monitor{}
}

// GetInfo retrieves the current battery information
func (m *Monitor) GetInfo() (Info, error) {
	out, err := exec.Command("pmset", "-g", "batt").Output()
	if err != nil {
		return Info{}, fmt.Errorf("failed to get battery information: %w", err)
	}

	percentage, err := parsePercentage(string(out))
	if err != nil {
		return Info{}, fmt.Errorf("failed to parse battery percentage: %w", err)
	}

	state := parseState(string(out))

	return Info{
		Percentage: percentage,
		State:      state,
	}, nil
}

// parsePercentage extracts the battery percentage from pmset output
func parsePercentage(output string) (int, error) {
	re := regexp.MustCompile(`[0-9]+%`)
	percentageStr := re.FindString(output)
	if percentageStr == "" {
		return 0, fmt.Errorf("no percentage found in output")
	}

	percentage, err := strconv.Atoi(percentageStr[:len(percentageStr)-1])
	if err != nil {
		return 0, fmt.Errorf("failed to convert percentage: %w", err)
	}

	return percentage, nil
}

// parseState extracts the battery state from pmset output
func parseState(output string) State {
	if regexp.MustCompile(`discharging`).MatchString(output) {
		return StateDischarging
	} else if regexp.MustCompile(`charging`).MatchString(output) {
		return StateCharging
	}
	return StateUnknown
}

// IsCharging returns true if the battery is charging
func (i Info) IsCharging() bool {
	return i.State == StateCharging
}

// IsDischarging returns true if the battery is discharging
func (i Info) IsDischarging() bool {
	return i.State == StateDischarging
}

// LogInfo logs the battery information
func (i Info) LogInfo() {
	log.Printf("Battery: %d%% - %s", i.Percentage, i.State)
}
