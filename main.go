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
	"fyne.io/fyne/v2/dialog"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("RFT Tool | Real Flash Tool")

	// Set default theme to dark (black theme)
	myApp.Settings().SetTheme(theme.DarkTheme())

	// Log area
	logArea := widget.NewMultiLineEntry()
	logArea.Disable() // Make the log area read-only
	logArea.Wrapping = fyne.TextWrapWord // Enable text wrapping

	// Set text color for the log area

	

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
	// adb read info filter function to extract specific details
	filterAdbReadInfo := func(output string) string {
		// Define the keys and their corresponding labels
	keys := map[string]string{
		"ro.product.brand":                "Brand",
		"ro.product.model":                "Model",
		"ro.product.name":                 "Product",
		"ro.boot.hardware":				   "Cpu",
		"ro.boot.hwlevel": 				   "Hardware Level",
		"ro.hardware.info":                "Hardware info",
		"ro.secureboot.lockstate":		   "Bootloader",
		"ro.build.version.release":        "Version",
		"ro.build.version.incremental":	   "Android Version",
		"ro.board.platform":               "Android platform",
		"ro.build.version.security_patch": "Security patch",
		"ro.build.display.id":             "Software version",
		"persist.sys.timezone":            "Time Zone",
		"ro.secure":                       "Root Access",
		"ro.ril.miui.imei":                "Device IMEI",
		"ro.ril.miui.imei0":               "Device IMEI1",
		"ro.ril.miui.imei2":               "Device IMEI2",
		"ro.frp.pst":                      "FRP PST",
		"ro.boot.slot":            		   "Active Slot",
		"ro.crypto.state":                 "Crypto State",
	}

	// Split the output into lines
	lines := strings.Split(output, "\n")

	// Create a map to store the parsed values
	parsedValues := make(map[string]string)

	// Parse the output and store the values
	for _, line := range lines {
		for key, label := range keys {
			if strings.Contains(line, key) {
				value := strings.TrimSpace(strings.Split(line, ":")[1])
				parsedValues[label] = value
				break
			}
		}
	}

	// Add additional logic for Root Access
	if parsedValues["Root Access"] == "1" {
		parsedValues["Root Access"] = "Granted"
	} else {
		parsedValues["Root Access"] = "Denied"
	}

	// Format the adb read info output
	var formattedOutput strings.Builder
	formattedOutput.WriteString("Read Information\n")
	formattedOutput.WriteString(fmt.Sprintf("Brand: %s\n", parsedValues["Brand"]))
	formattedOutput.WriteString(fmt.Sprintf("Model: %s\n", parsedValues["Model"]))
	formattedOutput.WriteString(fmt.Sprintf("Product: %s\n", parsedValues["Product"]))
	formattedOutput.WriteString(fmt.Sprintf("Bootloader: %s\n", parsedValues["Bootloader"]))
	formattedOutput.WriteString(fmt.Sprintf("Cpu: %s\n", parsedValues["Cpu"]))
	formattedOutput.WriteString(fmt.Sprintf("Hardware Level: %s\n", parsedValues["Hardware Level"]))
	formattedOutput.WriteString(fmt.Sprintf("Hardware info: %s\n", parsedValues["Hardware Info"]))
	formattedOutput.WriteString(fmt.Sprintf("Version: %s\n", parsedValues["Version"]))
	formattedOutput.WriteString(fmt.Sprintf("Android Version: %s\n", parsedValues["Android Version"]))
	formattedOutput.WriteString(fmt.Sprintf("Android platform: %s\n", parsedValues["Android platform"]))
	formattedOutput.WriteString(fmt.Sprintf("Security patch: %s\n", parsedValues["Security patch"]))
	formattedOutput.WriteString(fmt.Sprintf("Software version: %s\n", parsedValues["Software version"]))
	formattedOutput.WriteString(fmt.Sprintf("Time Zone: %s\n", parsedValues["Time Zone"]))
	formattedOutput.WriteString(fmt.Sprintf("Root Access: %s\n", parsedValues["Root Access"]))
	formattedOutput.WriteString(fmt.Sprintf("Device IMEI: %s\n", parsedValues["Device IMEI"]))
	formattedOutput.WriteString(fmt.Sprintf("Device IMEI1: %s\n", parsedValues["Device IMEI1"]))
	formattedOutput.WriteString(fmt.Sprintf("Device IMEI2: %s\n", parsedValues["Device IMEI2"]))
	formattedOutput.WriteString(fmt.Sprintf("FRP PST: %s\n", parsedValues["FRP PST"]))
	formattedOutput.WriteString(fmt.Sprintf("Active Slot: %s\n", parsedValues["Active Slot"]))
	formattedOutput.WriteString(fmt.Sprintf("Crypto State: %s\n", parsedValues["Crypto State"]))
	formattedOutput.WriteString("Operation time: 00:03\n") // Hardcoded for now

	return formattedOutput.String()
}

// scan adb and fastboot devices button with pop up message
scanDevices := widget.NewButton("Scan Devices", func() {
	clearLog() // Clear the log before adding new content
	output, err := runCommand("adb", "devices")
	if err != nil {
		logArea.SetText(logArea.Text + "Error scanning devices: " + err.Error() + "\n")
		return
	}
	logArea.SetText(logArea.Text + "Scanning devices...\n" + output + "\n")

	// Check if the output contains "unauthorized"
	if strings.Contains(output, "unauthorized") {
		dialog.ShowInformation("Scan Devices", "Please allow USB debugging on your device.", myWindow)
	}
})


// ADB Tab

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
	adbReadInfo := widget.NewButton("Adb Read Info", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "shell", "getprop")
		if err != nil {
			logArea.SetText(logArea.Text + "Error reading adb info: " + err.Error() + "\n")
			return
		}
		filteredOutput := filterAdbReadInfo(output)
		logArea.SetText(logArea.Text + "Reading Adb Info...\n" + filteredOutput + "\n")
	})
	adbRebootRecovery := widget.NewButton("Adb reboot recovery", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "reboot", "recovery")
		if err != nil {
			logArea.SetText(logArea.Text + "Error rebooting to recovery: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Successfully rebooted to recovery...\n" + output + "\n")
	})
	diagWithoutRoot := widget.NewButton("Diag without root", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "shell", "am", "start", "-n", "com.android.settings/.DevelopmentSettings")
		if err != nil {
			logArea.SetText(logArea.Text + "Error opening Diag without root: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Successfully opened Diag without root...\n" + output + "\n")
	})
	diagWithRoot := widget.NewButton("Diag with root", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "shell", "am", "start", "-n", "com.android.settings/.DevelopmentSettings")
		if err != nil {
			logArea.SetText(logArea.Text + "Error opening Diag with root: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Successfully opened Diag with root...\n" + output + "\n")
	})
	// create wipe efs need backup then wipe and warrning message need yes to reset no to cancel

	
	WipeEfs := widget.NewButton("Wipe Efs", func() {
		// Show a confirmation dialog before proceeding
		clearLog() // Clear the log before adding new content
		dialog.ShowConfirm(
			"Wipe Efs", // Title
			"This will delete all data in the Efs partition. Make sure you have a backup! Do you want to continue?", // Message
			func(response bool) { // Callback function
				if response {
					// User clicked "Yes"
					logArea.SetText(logArea.Text + "Wiping Efs...\n")

					// Simulate running a command (replace with actual command execution)
					output, err := runCommand("adb", "shell", "rm", "-rf", "/efs")
					if err != nil {
						logArea.SetText(logArea.Text + "Error Wiping Efs: " + err.Error() + "\n")
						return
					}

					// Log success message
					logArea.SetText(logArea.Text + "Successfully Wiped Efs...\n" + output + "\n")
				} else {
					// User clicked "No" or closed the dialog
					logArea.SetText(logArea.Text + "Wipe Efs operation canceled.\n")
				}
			},
			myWindow, // Parent window for the dialog
		)
	})


	
	RebootEdl := widget.NewButton("Reboot Edl", func() {
		clearLog() // Clear the log before adding new content
		output, err := runCommand("adb", "reboot", "edl")
		if err != nil {
			logArea.SetText(logArea.Text + "Error rebooting to Edl: " + err.Error() + "\n")
			return
		}
		logArea.SetText(logArea.Text + "Successfully rebooted to Edl...\n" + output + "\n")
	})
	
	

	adbButtons := container.NewGridWithColumns(
		5, // Number of columns
		scanDevices,
		adbCheckButton,
		adbRebootButton,
		adbToBootloader,
		adbReadInfo,
		adbRebootRecovery,
		diagWithoutRoot,
		diagWithRoot,
		WipeEfs,
		RebootEdl,
	)
	
	adbTab := container.NewVBox(
		adbButtons, // Buttons arranged in a 5-column grid
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
	fastbootButton := container.NewGridWithColumns(
		5, // 5 columns
		fastbootCheckButton,
		fastbootReadInfoButton,
		fastbootRebootButton,
		fastbootFutureButton1,
		fastbootFutureButton2,
	)

	fastbootTab := container.NewVBox(
		fastbootButton, // Buttons side by side
	)

	// Xiaomi Special Tab
	xiaomiSpecialButton1 := widget.NewButton("Button 1", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 1 clicked...\n")
	})
	xiaomiSpecialButton2 := widget.NewButton("Button 2", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 2 clicked...\n")
	})
	xiaomiSpecialButton3 := widget.NewButton("Button 3", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 3 clicked...\n")
	})
	xiaomiSpecialButton4 := widget.NewButton("Button 4", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 4 clicked...\n")
	})
	xiaomiSpecialButton5 := widget.NewButton("Button 5", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 5 clicked...\n")
	})
	xiaomiSpecialButton6 := widget.NewButton("Button 6", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 6 clicked...\n")
	})
	xiaomiSpecialButton7 := widget.NewButton("Button 7", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 7 clicked...\n")
	})
	xiaomiSpecialButton8 := widget.NewButton("Button 8", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 8 clicked...\n")
	})
	xiaomiSpecialButton9 := widget.NewButton("Button 9", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 9 clicked...\n")
	})
	xiaomiSpecialButton10 := widget.NewButton("Button 10", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Special Button 10 clicked...\n")
	})

	// Set a smaller size for the buttons
	smallButtonSize := fyne.NewSize(80, 30) // Width: 80, Height: 30
	xiaomiSpecialButton1.Resize(smallButtonSize)
	xiaomiSpecialButton2.Resize(smallButtonSize)
	xiaomiSpecialButton3.Resize(smallButtonSize)
	xiaomiSpecialButton4.Resize(smallButtonSize)
	xiaomiSpecialButton5.Resize(smallButtonSize)
	xiaomiSpecialButton6.Resize(smallButtonSize)
	xiaomiSpecialButton7.Resize(smallButtonSize)
	xiaomiSpecialButton8.Resize(smallButtonSize)
	xiaomiSpecialButton9.Resize(smallButtonSize)
	xiaomiSpecialButton10.Resize(smallButtonSize)

	// Xiaomi Special Buttons arranged in a grid
	xiaomiSpecialTab := container.NewGridWithColumns(
		5, // 5 columns
		xiaomiSpecialButton1,
		xiaomiSpecialButton2,
		xiaomiSpecialButton3,
		xiaomiSpecialButton4,
		xiaomiSpecialButton5,
		xiaomiSpecialButton6,
		xiaomiSpecialButton7,
		xiaomiSpecialButton8,
		xiaomiSpecialButton9,
		xiaomiSpecialButton10,
	)

	// MediaTek Tab
	mediatekButton1 := widget.NewButton("Button 1", func() {
		clearLog()
		logArea.SetText(logArea.Text + "MediaTek Button 1 clicked...\n")
	})
	mediatekButton2 := widget.NewButton("Button 2", func() {
		clearLog()
		logArea.SetText(logArea.Text + "MediaTek Button 2 clicked...\n")
	})

	// Set smaller size for MediaTek buttons
	mediatekButton1.Resize(smallButtonSize)
	mediatekButton2.Resize(smallButtonSize)

	mediatekTab := container.NewGridWithColumns(
		2, // 2 columns
		mediatekButton1,
		mediatekButton2,
	)

	// Flash EDL Tab
	flashEDLButton1 := widget.NewButton("Button 1", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Flash EDL Button 1 clicked...\n")
	})
	flashEDLButton2 := widget.NewButton("Button 2", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Flash EDL Button 2 clicked...\n")
	})

	// Set smaller size for Flash EDL buttons
	flashEDLButton1.Resize(smallButtonSize)
	flashEDLButton2.Resize(smallButtonSize)

	flashEDLTab := container.NewGridWithColumns(
		2, // 2 columns
		flashEDLButton1,
		flashEDLButton2,
	)

	// Fastboot Tab (under Xiaomi)
	xiaomiFastbootButton1 := widget.NewButton("Button 1", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Fastboot Button 1 clicked...\n")
	})
	xiaomiFastbootButton2 := widget.NewButton("Button 2", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Xiaomi Fastboot Button 2 clicked...\n")
	})

	// Set smaller size for Xiaomi Fastboot buttons
	xiaomiFastbootButton1.Resize(smallButtonSize)
	xiaomiFastbootButton2.Resize(smallButtonSize)

	xiaomiFastbootTab := container.NewGridWithColumns(
		2, // 2 columns
		xiaomiFastbootButton1,
		xiaomiFastbootButton2,
	)

	// Xiaomi Main Tab (nested tabs)
	xiaomiTabs := container.NewAppTabs(
		container.NewTabItem("Xiaomi Special", xiaomiSpecialTab),
		container.NewTabItem("MediaTek", mediatekTab),
		container.NewTabItem("Flash EDL", flashEDLTab),
		container.NewTabItem("Fastboot", xiaomiFastbootTab),
	)

	// Samsung Tab
	samsungButton1 := widget.NewButton("Samsung Button 1", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Samsung Button 1 clicked...\n")
	})
	samsungButton2 := widget.NewButton("Samsung Button 2", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Samsung Button 2 clicked...\n")
	})
	samsungButton3 := widget.NewButton("Samsung Button 3", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Samsung Button 3 clicked...\n")
	})
	samsungButton4 := widget.NewButton("Samsung Button 4", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Samsung Button 4 clicked...\n")
	})
	samsungButton5 := widget.NewButton("Samsung Button 5", func() {
		clearLog()
		logArea.SetText(logArea.Text + "Samsung Button 5 clicked...\n")
	})

	// Samsung Buttons side by side
	samsungButtons := container.NewHBox(
		samsungButton1,
		samsungButton2,
		samsungButton3,
		samsungButton4,
		samsungButton5,
	)

	samsungTab := container.NewVBox(
		samsungButtons, // Buttons side by side
	)

	// Main Tabs (Navbar at the top)
	tabs := container.NewAppTabs(
		container.NewTabItem("ADB Device", adbTab),
		container.NewTabItem("Fastboot Device", fastbootTab),
		container.NewTabItem("Xiaomi", xiaomiTabs), // Nested tabs for Xiaomi
		container.NewTabItem("Samsung", samsungTab),
	)

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