#!/bin/bash

# Set variables
BATTERY_REMINDER="battery-reminder"
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
    nohup "$PWD"/build/"$BATTERY_REMINDER" --min "$MIN" --max "$MAX" >"$BATTERY_REMINDER".log 2>&1 &
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

remove_from_zshrc() {
  # Remove existing lines between <Battery Reminder> and </Battery Reminder>
  sed -i '' '/# <Battery Reminder>/,/# <\/Battery Reminder>/d' ~/.zshrc
}

add_battery_reminder_to_zshrc() {
  local B_MIN="$1"
  local B_MAX="$2"

  local COMMAND=$(
    cat <<EOF

# <Battery Reminder>
BATTERY_REMINDER="$BATTERY_REMINDER"
B_APP_ROOT="$PWD"
B_MIN="$B_MIN"
B_MAX="$B_MAX"
if ! pgrep -f "$BATTERY_REMINDER" >/dev/null; then
  nohup "\$B_APP_ROOT/build/\$BATTERY_REMINDER" --min "\$B_MIN" --max "\$B_MAX" >"\$B_APP_ROOT/\$BATTERY_REMINDER.log" 2>&1 &
fi
# </Battery Reminder>
EOF
  )

  # Remove
  remove_from_zshrc
  # Append the updated command to .zshrc
  echo "$COMMAND" >>~/.zshrc
}

# Determine the command
case "$COMMAND" in
"stop")
  stop_battery_reminder
  remove_from_zshrc
  ;;
*)

  # Set defaults for min and max
  if [[ -z "$MIN" ]]; then
    MIN=$DEFAULT_MIN
  fi
  if [[ -z "$MAX" ]]; then
    MAX=$DEFAULT_MAX
  fi

  # Stop the program if it is already running
  if is_battery_reminder_running; then
    stop_battery_reminder
  fi

  # Start the program
  start_battery_reminder

  # add_battery_reminder_to_zshrc
  add_battery_reminder_to_zshrc $MIN $MAX

  echo "Installation complete. Battery Reminder will now run at startup with min=$MIN and max=$MAX."
  ;;
esac
