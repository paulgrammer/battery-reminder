# Battery Reminder

A lightweight Go application that monitors your MacBook's battery level and reminds you when it's time to plug in or unplug your charger. This helps optimize battery usage and extend the lifespan of your MacBook's battery.

## Purpose

This program helps MacBook users maintain their battery health by:
- Alerting when battery drops to a minimum threshold (default: 20%)
- Alerting when battery reaches a maximum threshold (default: 80%)
- Preventing overcharging and deep discharging, which can reduce battery lifespan

## Features

- **Customizable thresholds**: Set your own minimum and maximum battery percentages
- **Multiple alert methods**: Visual notifications, sound alerts, and text-to-speech
- **System service support**: Run as a background service (LaunchAgent on macOS)
- **Smart monitoring**: Adjusts check frequency based on battery level
- **Modular architecture**: Clean, maintainable codebase with separated concerns


## Installation

### Prerequisites

- Go 1.25.0 or later
- macOS (uses pmset, osascript, and afplay)

### Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/paulgrammer/battery-reminder.git
   cd battery-reminder
   ```

2. Build the application:
   ```bash
   go build -o bin/battery-reminder ./cmd/battery-reminder
   ```

## Usage

### Basic Usage

Run the program with default settings (20% minimum, 80% maximum):

```bash
./bin/battery-reminder
```

### Custom Thresholds

Specify custom battery percentage thresholds:

```bash
./bin/battery-reminder -min 15 -max 85
```

### Custom Sound Alert

Use a custom sound file for alerts:

```bash
./bin/battery-reminder -sound /System/Library/Sounds/Ping.aiff
```

### Command-line Options

```
-min int
    Minimum battery percentage (default: 20)
-max int
    Maximum battery percentage (default: 80)
-sound string
    Sound file to play for alerts (e.g., /System/Library/Sounds/Ping.aiff)
-service string
    Control the system service: install, uninstall, start, stop, restart, run
```

## Running as a System Service

### Install the Service

Install battery-reminder as a system service that starts automatically:

```bash
./bin/battery-reminder -service install -min 20 -max 80
```

### Start the Service

```bash
./bin/battery-reminder -service start
```

### Stop the Service

```bash
./bin/battery-reminder -service stop
```

### Restart the Service

```bash
./bin/battery-reminder -service restart
```

### Uninstall the Service

```bash
./bin/battery-reminder -service uninstall
```

### Service Management

Once installed, the service runs in the background and starts automatically on system boot. Logs can be viewed using:

```bash
# macOS
log show --predicate 'processImagePath contains "battery-reminder"' --info --last 1h
```

## How It Works

1. **Battery Monitoring**: Uses `pmset -g batt` to retrieve current battery status
2. **State Detection**: Determines if the battery is charging, discharging, or in an unknown state
3. **Smart Intervals**: Calculates sleep time based on distance from thresholds
4. **Alert System**:
   - Shows macOS notification
   - Unmutes system volume
   - Plays custom sound or uses text-to-speech
   - Repeats alerts every 10 seconds while outside threshold range

## Architecture

The application is built with a modular architecture:

- **battery**: Low-level battery information retrieval
- **config**: Configuration validation and management
- **notifier**: All notification mechanisms (visual, audio, speech)
- **monitor**: High-level monitoring logic and coordination
- **svc**: System service integration using kardianos/service

This design ensures:
- Easy testing of individual components
- Clear separation of concerns
- Maintainable and extensible codebase
- Proper error handling throughout

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o bin/battery-reminder ./cmd/battery-reminder

# Build with version info
go build -ldflags "-X main.version=1.0.0" -o bin/battery-reminder ./cmd/battery-reminder
```

## Compatibility

- **macOS**: Fully supported (10.9+)
- **Other OS**: Not supported (relies on macOS-specific commands)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Disclaimer

This program is a monitoring tool and not a replacement for built-in battery management features provided by Apple. Use at your own discretion.

## Acknowledgments

- [kardianos/service](https://github.com/kardianos/service) - Cross-platform service management library
