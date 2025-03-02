package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/canvas"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Mobile Flash Tool")

	// Set default theme to dark (black theme)
	myApp.Settings().SetTheme(theme.DarkTheme())

	// Log area
	logArea := widget.NewMultiLineEntry()
	logArea.Disable() // Make the log area read-only
	logArea.Wrapping = fyne.TextWrapWord // Enable text wrapping

	// Function to clear the log
	clearLog := func() {
		logArea.SetText("")
	}

	// Function to execute shell commands and log output with error handling
	runCommand := func(command string, args ...string) (string, error) {
		cmd := exec.Command(command, args...)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("command failed: %s\n%s", err, string(output))
		}
		return string(output), nil
	}

	// Function to filter specific details from `fastboot getvar all` output
	filterFastbootGetvar := func(output string) string {
		// Define the keys you want to extract
		keys := []string{"anti", "serialno", "product", "version", "secure", "unlocked"}

		// Split the output into lines
		lines := strings.Split(output, "\n")

		// Filter lines that contain the keys
		var filteredOutput strings.Builder
		for _, line := range lines {
			for _, key := range keys {
				if strings.Contains(line, key+":") {
					filteredOutput.WriteString(line + "\n")
					break
				}
			}
		}

		return filteredOutput.String()
	}

	// Function to check ADB devices and handle unauthorized state
	checkADB := func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "devices")
		if err != nil {
			logArea.SetText(logArea.Text + "Error checking ADB: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Checking ADB...\n" + output + "\n")

		// Check if the output contains "unauthorized"
		if strings.Contains(output, "unauthorized") {
			logArea.SetText(logArea.Text + "Please allow USB debugging on your device.\n")
			logArea.SetText(logArea.Text + "1. Check your device screen.\n")
			logArea.SetText(logArea.Text + "2. Tap 'Allow' for USB debugging.\n")
		}
	}

	// ADB Tab
	adbCheckButton := widget.NewButton("Check ADB", func() {
		checkADB() // Call the checkADB function
	})
	adbRebootButton := widget.NewButton("Reboot Device", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "reboot")
		if err != nil {
			logArea.SetText(logArea.Text + "Error rebooting device: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Successfully rebooted device...\n" + output + "\n")
	})
	adbToBootloader := widget.NewButton("Adb to Bootloader", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "reboot", "bootloader")
		if err != nil {
			logArea.SetText(logArea.Text + "Error rebooting to bootloader: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Successfully rebooted to bootloader...\n" + output + "\n")
	})
	adbFutureButton2 := widget.NewButton("Future Button 2", func() {
		clearLog() // Clear the log before adding new content
		logArea.SetText(logArea.Text + "Future Button 2 clicked...\n")
	})
	adbFutureButton3 := widget.NewButton("Future Button 3", func() {
		clearLog() // Clear the log before adding new content
		logArea.SetText(logArea.Text + "Future Button 3 clicked...\n")
	})

	// ADB Buttons side by side
	adbButtons := container.NewHBox(
		adbCheckButton,
		adbRebootButton,
		adbToBootloader,
		adbFutureButton2,
		adbFutureButton3,
	)

	adbTab := container.NewVBox(
		adbButtons, // Buttons side by side
	)

	// Fastboot Tab
	fastbootCheckButton := widget.NewButton("Check Fastboot", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("fastboot", "devices")
		if err != nil {
			logArea.SetText(logArea.Text + "Error checking Fastboot devices: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Checking Fastboot Devices...\n" + output + "\n")
	})

	// Fastboot Info Button
	fastbootReadInfoButton := widget.NewButton("Fastboot Read Info", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("fastboot", "getvar", "all")
		if err != nil {
			logArea.SetText(logArea.Text + "Error reading Fastboot info: " + err.Error() + "\n")
			return
		}
		filteredOutput := filterFastbootGetvar(output)
		logArea.SetText(logArea.Text + "Reading Fastboot Info...\n" + filteredOutput + "\n")
	})

	fastbootRebootButton := widget.NewButton("Reboot Fastboot", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("fastboot", "reboot")
		if err != nil {
			logArea.SetText(logArea.Text + "Error rebooting Fastboot devices: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Rebooting Fastboot Devices...\n" + output + "\n")
	})
	fastbootFutureButton1 := widget.NewButton("Future Button 1", func() {
		clearLog() // Clear the log before adding new content
		logArea.SetText(logArea.Text + "Future Button 1 clicked...\n")
	})
	fastbootFutureButton2 := widget.NewButton("Future Button 2", func() {
		clearLog() // Clear the log before adding new content
		logArea.SetText(logArea.Text + "Future Button 2 clicked...\n")
	})

	// Fastboot Buttons side by side
	fastbootButton := container.NewHBox(
		fastbootCheckButton,
		fastbootReadInfoButton,
		fastbootRebootButton,
		fastbootFutureButton1,
		fastbootFutureButton2,
	)

	fastbootTab := container.NewVBox(
		fastbootButton, // Buttons side by side
	)

	// Tabs (Navbar at the top)
	tabs := container.NewAppTabs(
		container.NewTabItem("ADB Device", adbTab),
		container.NewTabItem("Fastboot Device", fastbootTab),
	)

	// Xiaomi Tab
	xiaomiTab := container.NewVBox(
		widget.NewLabel("Xiaomi"),
	)
	// Declare Xiaomi tab
	Xiaomi := container.NewTabItem("Xiaomi", xiaomiTab)
	// Add Xiaomi tab to tabs
	tabs.Append(Xiaomi)

	// Samsung Tab
	samsungTab := container.NewVBox(
		widget.NewLabel("Samsung"),
	)
	// Declare Samsung tab
	Samsung := container.NewTabItem("Samsung", samsungTab)
	// Add Samsung tab to tabs
	tabs.Append(Samsung)

	// Time and Date
	timeLabel := widget.NewLabel("")
	go func() {
		for range time.Tick(time.Second) {
			timeLabel.SetText(time.Now().Format("2006-01-02 15:04:05"))
		}
	}()

	// Custom log area styling
	logContainer := container.NewStack(
		canvas.NewRectangle(theme.BackgroundColor()), // Background color
		container.NewScroll(logArea), // Scrollable log area
	)

	// Layout
	content := container.NewBorder(
		container.NewVBox(
			container.NewHBox(
				tabs, // Tabs on the left
			),
		),
		timeLabel, // Time and date at the bottom
		nil,
		nil,
		logContainer, // Log area in the center
	)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}