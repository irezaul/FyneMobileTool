package main

import (
    "strings"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// Function to filter specific details from `fastboot getvar all` output
func filterFastbootGetvar(output string) string {
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

// Fastboot Tab UI components
func createFastbootTab(logArea *widget.Entry, clearLog func(), runCommand func(string, ...string) (string, error)) fyne.CanvasObject {
    // Fastboot Buttons
    fastbootCheckButton := widget.NewButton("Check Fastboot", func() {
        clearLog() // Clear the log before adding new content
        output, err := runCommand("fastboot", "devices")
        if err != nil {
            logArea.SetText(logArea.Text + "Error checking Fastboot devices: " + err.Error() + "\n")
            return
        }
        logArea.SetText(logArea.Text + "Checking Fastboot Devices...\n" + output + "\n")
    })

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

    // Fastboot Tab layout
    return container.NewVBox(
        fastbootButton, // Buttons side by side
    )
}