package main

import (
    "fmt"
    "strings"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
)

// ADB-related functions and UI components

// Function to filter specific details from `adb shell getprop` output
func filterAdbReadInfo(output string) string {
    keys := map[string]string{
        "ro.product.brand":                "Brand",
        "ro.product.model":                "Model",
        "ro.product.name":                 "Product",
        "ro.boot.hardware":                "Cpu",
        "ro.boot.hwlevel":                 "Hardware Level",
        "ro.hardware.info":                "Hardware info",
        "ro.secureboot.lockstate":         "Bootloader",
        "ro.build.version.release":        "Version",
        "ro.build.version.incremental":    "Android Version",
        "ro.board.platform":               "Android platform",
        "ro.build.version.security_patch": "Security patch",
        "ro.build.display.id":             "Software version",
        "persist.sys.timezone":            "Time Zone",
        "ro.secure":                       "Root Access",
        "ro.ril.miui.imei":                "Device IMEI",
        "ro.ril.miui.imei0":               "Device IMEI1",
        "ro.ril.miui.imei2":               "Device IMEI2",
        "ro.frp.pst":                      "FRP PST",
        "ro.boot.slot":                    "Active Slot",
        "ro.crypto.state":                 "Crypto State",
    }

    lines := strings.Split(output, "\n")
    parsedValues := make(map[string]string)

    for _, line := range lines {
        for key, label := range keys {
            if strings.Contains(line, key) {
                value := strings.TrimSpace(strings.Split(line, ":")[1])
                parsedValues[label] = value
                break
            }
        }
    }

    if parsedValues["Root Access"] == "1" {
        parsedValues["Root Access"] = "Granted"
    } else {
        parsedValues["Root Access"] = "Denied"
    }

    var formattedOutput strings.Builder
    formattedOutput.WriteString("Read Information\n")
    for key, value := range parsedValues {
        formattedOutput.WriteString(fmt.Sprintf("%s: %s\n", key, value))
    }
    formattedOutput.WriteString("Operation time: 00:03\n") // Hardcoded for now

    return formattedOutput.String()
}

// ADB Tab UI components
func createAdbTab(logArea *widget.Entry, clearLog func(), runCommand func(string, ...string) (string, error), myWindow fyne.Window) fyne.CanvasObject {
    // ADB Buttons
    scanDevices := widget.NewButton("Scan Devices", func() {
        clearLog()
        output, err := runCommand("adb", "devices")
        if err != nil {
            logArea.SetText(logArea.Text + "Error scanning devices: " + err.Error() + "\n")
            return
        }
        logArea.SetText(logArea.Text + "Scanning devices...\n" + output + "\n")

        if strings.Contains(output, "unauthorized") {
            dialog.ShowInformation("Scan Devices", "Please allow USB debugging on your device.", myWindow)
        }
    })

    checkADB := func() {
        clearLog()
        output, err := runCommand("adb", "devices")
        if err != nil {
            logArea.SetText(logArea.Text + "Error checking ADB: " + err.Error() + "\n")
            return
        }
        logArea.SetText(logArea.Text + "Checking ADB...\n" + output + "\n")

        if strings.Contains(output, "unauthorized") {
            logArea.SetText(logArea.Text + "Please allow USB debugging on your device.\n")
            logArea.SetText(logArea.Text + "1. Check your device screen.\n")
            logArea.SetText(logArea.Text + "2. Tap 'Allow' for USB debugging.\n")
        }
    }

    adbCheckButton := widget.NewButton("Check ADB", func() {
        checkADB()
    })
    adbRebootButton := widget.NewButton("Reboot Device", func() {
        clearLog()
        output, err := runCommand("adb", "reboot")
        if err != nil {
            logArea.SetText(logArea.Text + "Error rebooting device: " + err.Error() + "\n")
            return
        }
        logArea.SetText(logArea.Text + "Successfully rebooted device...\n" + output + "\n")
    })
    adbToBootloader := widget.NewButton("Adb to Bootloader", func() {
        clearLog()
        output, err := runCommand("adb", "reboot", "bootloader")
        if err != nil {
            logArea.SetText(logArea.Text + "Error rebooting to bootloader: " + err.Error() + "\n")
            return
        }
        logArea.SetText(logArea.Text + "Successfully rebooted to bootloader...\n" + output + "\n")
    })
    adbReadInfo := widget.NewButton("Adb Read Info", func() {
        clearLog()
        output, err := runCommand("adb", "shell", "getprop")
        if err != nil {
            logArea.SetText(logArea.Text + "Error reading adb info: " + err.Error() + "\n")
            return
        }
        filteredOutput := filterAdbReadInfo(output)
        logArea.SetText(logArea.Text + "Reading Adb Info...\n" + filteredOutput + "\n")
    })
    adbRebootRecovery := widget.NewButton("Adb reboot recovery", func() {
        clearLog()
        output, err := runCommand("adb", "reboot", "recovery")
        if err != nil {
            logArea.SetText(logArea.Text + "Error rebooting to recovery: " + err.Error() + "\n")
            return
        }
        logArea.SetText(logArea.Text + "Successfully rebooted to recovery...\n" + output + "\n")
    })
    
    
    

    adbButtons := container.NewGridWithColumns(
        5,
        scanDevices,
        adbCheckButton,
        adbRebootButton,
        adbToBootloader,
        adbReadInfo,
        adbRebootRecovery,
        
        
        
    )

    return container.NewVBox(adbButtons)
}