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
    myWindow := myApp.NewWindow("RFT Tool | Real Flash Tool")

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

    // Create input boxes for ADB and Fastboot ports
   // Create a single input box for ADB and Fastboot port information
portInfoEntry := widget.NewEntry()
portInfoEntry.SetPlaceHolder("Port Information")

// Button to scan devices
scanDevicesButton := widget.NewButton("Scan Devices", func() {
    clearLog() // Clear the log area

    // Initialize port information message
    portInfoMessage := ""

    // Scan for ADB devices
    adbOutput, adbErr := runCommand("adb", "devices")
    if adbErr != nil {
        portInfoMessage += "ADB: Error scanning devices\n"
    } else {
        if strings.TrimSpace(adbOutput )!= "" {
            portInfoMessage += "ADB: Port 5037\n"
        } else {
            portInfoMessage += "ADB: No devices found\n"
        }
    }

    // Scan for Fastboot devices
    fastbootOutput, fastbootErr := runCommand("fastboot", "devices")
    if fastbootErr != nil {
        portInfoMessage += "Fastboot: Error scanning devices\n"
    } else {
        if strings.TrimSpace(fastbootOutput) != "" {
            portInfoMessage += "Fastboot: Port 5554\n"
        } else {
            portInfoMessage += "Fastboot: No devices found\n"
        }
    }

    // Display the port information in the input box
    portInfoEntry.SetText(portInfoMessage)
})

// Layout for the scan devices section
scanDevicesSection := container.NewVBox(
    container.NewHBox(
        widget.NewLabel("Port Info:"),
        portInfoEntry,
    ),
    scanDevicesButton,
)

    // Tabs
    adbTab := createAdbTab(logArea, clearLog, runCommand, myWindow)
    fastbootTab := createFastbootTab(logArea, clearLog, runCommand)
    samsungTab := createSamsungTab(logArea, clearLog)

    // Main Tabs (Navbar at the top)
    tabs := container.NewAppTabs(
        container.NewTabItem("ADB Device", adbTab),
        container.NewTabItem("Fastboot Device", fastbootTab),
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
        container.NewScroll(logArea),                 // Scrollable log area
    )

    // Layout
    content := container.NewBorder(
        container.NewVBox(
            scanDevicesSection, // Add the scan devices section here
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