#!/bin/bash

# Set variables
BATTERY_REMINDER="battery-reminder"
LAUNCH_AGENT_FILE="com.user.$BATTERY_REMINDER.plist"
LAUNCH_AGENT_PATH="$HOME/Library/LaunchAgents/$LAUNCH_AGENT_FILE"
COMMAND="$1"
MIN="$2"
MAX="$3"
DEFAULT_MIN=20
DEFAULT_MAX=80

# Check if the program is already running
is_battery_reminder_running() {
  pgrep -x "$BATTERY_REMINDER" >/dev/null
}

# Start the program
start_battery_reminder() {
  if is_battery_reminder_running; then
    echo "Battery Reminder is already running."
  else
    nohup ./build/"$BATTERY_REMINDER" --min "$MIN" --max "$MAX" > "$BATTERY_REMINDER".log 2>&1 &
    echo "Battery Reminder started."
  fi
}

# Stop the program
stop_battery_reminder() {
  if is_battery_reminder_running; then
    pkill -x "$BATTERY_REMINDER"
    echo "Battery Reminder stopped."
  else
    echo "Battery Reminder is not running."
  fi
}

# Create the LaunchAgent plist file
create_launch_agent_plist() {
  cat <<EOT >"$LAUNCH_AGENT_PATH"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>$LAUNCH_AGENT_FILE</string>
  <key>ProgramArguments</key>
  <array>
    <string>$(pwd)/build/$BATTERY_REMINDER</string>
    <string>--min</string>
    <string>"$MIN"</string>
    <string>--max</string>
    <string>"$MAX"</string>
  </array>
  <key>RunAtLoad</key>
  <true/>
</dict>
</plist>
EOT
}

# Load the LaunchAgent plist file
load_launch_agent_plist() {
  launchctl load "$LAUNCH_AGENT_PATH"
}

# Determine the command
case "$COMMAND" in
"stop")
  stop_battery_reminder
  ;;
*)
  # Delete existing files if they exist
  if [[ -f "$LAUNCH_AGENT_PATH" ]]; then
    launchctl unload "$LAUNCH_AGENT_PATH"
    rm "$LAUNCH_AGENT_PATH"
  fi

  # Set defaults for min and max
  if [[ -z "$MIN" ]]; then
    MIN=$DEFAULT_MIN
  fi
  if [[ -z "$MAX" ]]; then
    MAX=$DEFAULT_MAX
  fi

  # Create the LaunchAgent plist file
  create_launch_agent_plist

  # Load the LaunchAgent plist file
  load_launch_agent_plist

  # Stop the program if it is already running
  if is_battery_reminder_running; then
    stop_battery_reminder
  fi

  # Start the program
  start_battery_reminder

  echo "Installation complete. Battery Reminder will now run at startup with min=$MIN and max=$MAX."
  ;;
esac

# ./setup.sh start 20 80
