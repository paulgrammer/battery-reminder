# Battery Reminder

This is a simple Go program that reminds MacBook Pro users when the battery level is at 20% or 80% charge, to help prolong the battery life of the MacBook.

## Purpose

The purpose of this program is to help MacBook Pro users optimize their battery usage and extend the lifespan of their MacBook's battery. By reminding the user when the battery level is at 20% or 80% charge, the user can avoid overcharging or completely discharging the battery, which can damage the battery and reduce its lifespan.

## How to use

1. Clone this repository to your local machine.
2. Make sure you have Go installed on your machine.
3. Open a terminal window and navigate to the directory where the program is located.
4. Run the following command to build the program:

   ```
   make build
   ```

5. Run the program using the following command:

   ```
   ./build/battery-reminder
   ```

   The program will run in the background and monitor the battery level of your MacBook Pro. It will remind you with a sound when the battery level is at 20% or 80% charge.

   You can customize the sound files by modifying the paths to the sound files in the `main()` function of the code.

## Running as a background process

To run the `battery-reminder` program as a background process using `nohup`, you can follow these steps:

1. Navigate to the directory where your `battery-reminder` program is located in a terminal window.

2. Run the following command to start your program as a background process:

   ```
   nohup ./build/battery-reminder > battery-reminder.log &
   ```

   This command starts your program in the background and redirects its output to a log file named `battery-reminder.log`. The `&` at the end of the command sends the process to the background.

3. To check the status of your program, you can run the following command:

   ```
   ps aux | grep battery-reminder
   ```

   This will show a list of all processes running on your system that contain the string "battery-reminder".

4. To stop your program, you can find its process ID (PID) using the `ps` command, and then use the `kill` command to stop it:

   ```
   ps aux | grep battery-reminder
   kill PID
   ```

   Replace `PID` with the process ID of your `battery-reminder` program.

That's it! Your `battery-reminder` program should now run as a background process using `nohup`. Note that this method does not automatically start your program when your Mac starts up, and you may need to manually start it each time you want to use it.

## Compatibility

This program has been tested on macOS, but it may work on other operating systems that support the `battery` package in Go.

## Disclaimer

This program is not intended to replace the built-in battery monitoring tools provided by Apple or to provide medical advice. Always consult with a qualified professional before making any changes to your usage patterns or lifestyle.

## License

This program is licensed under the MIT License. See the [LICENSE](https://opensource.org/licenses/MIT) file for details.