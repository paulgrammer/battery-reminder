package config

import "fmt"

// Config holds the application configuration
type Config struct {
	MinPercentage int
	MaxPercentage int
	SoundPath     string
}

// NewConfig creates a new configuration with validation
func NewConfig(minPercentage, maxPercentage int, soundPath string) (*Config, error) {
	if minPercentage < 0 || minPercentage > 100 {
		return nil, fmt.Errorf("minimum percentage must be between 0 and 100, got %d", minPercentage)
	}
	if maxPercentage < 0 || maxPercentage > 100 {
		return nil, fmt.Errorf("maximum percentage must be between 0 and 100, got %d", maxPercentage)
	}
	if minPercentage >= maxPercentage {
		return nil, fmt.Errorf("minimum percentage (%d) must be less than maximum percentage (%d)", minPercentage, maxPercentage)
	}

	return &Config{
		MinPercentage: minPercentage,
		MaxPercentage: maxPercentage,
		SoundPath:     soundPath,
	}, nil
}

// String returns a string representation of the config
func (c *Config) String() string {
	return fmt.Sprintf("Battery Reminder Started. \nMinimum: %d%% \nMaximum: %d%%", c.MinPercentage, c.MaxPercentage)
}
