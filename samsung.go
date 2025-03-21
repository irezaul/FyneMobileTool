package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "os/exec" 
)

// Function to create the Samsung Tab
func createSamsungTab(logArea *widget.Entry, clearLog func(), myWindow fyne.Window) fyne.CanvasObject {
    // General Samsung Buttons
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

    // File selection buttons
    selectPITButton := widget.NewButton("Select PIT File", func() {
        dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                logArea.SetText(logArea.Text + "PIT File Selected: " + reader.URI().Path() + "\n")
            }
        }, myWindow)
    })

    selectBLButton := widget.NewButton("Select BL File", func() {
        dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                logArea.SetText(logArea.Text + "BL File Selected: " + reader.URI().Path() + "\n")
            }
        }, myWindow)
    })

    selectAPButton := widget.NewButton("Select AP File", func() {
        dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                logArea.SetText(logArea.Text + "AP File Selected: " + reader.URI().Path() + "\n")
            }
        }, myWindow)
    })

    selectCPButton := widget.NewButton("Select CP File", func() {
        dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                logArea.SetText(logArea.Text + "CP File Selected: " + reader.URI().Path() + "\n")
            }
        }, myWindow)
    })

    selectCSCButton := widget.NewButton("Select CSC File", func() {
        dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
            if err == nil && reader != nil {
                logArea.SetText(logArea.Text + "CSC File Selected: " + reader.URI().Path() + "\n")
            }
        }, myWindow)
    })

    // Flash button
    flashButton := widget.NewButton("Flash", func() {
        clearLog()
        logArea.SetText(logArea.Text + "Flashing started...\n")
    
        // Replace "your-flashing-tool" with the actual tool name or full path
        toolPath := "C:\\path\\to\\your-flashing-tool.exe" // Update this with the correct path
    
        // Replace "arguments" with the actual arguments for the tool
        cmd := exec.Command(toolPath, "--pit", "path/to/pit/file", "--bl", "path/to/bl/file", "--ap", "path/to/ap/file", "--cp", "path/to/cp/file", "--csc", "path/to/csc/file")
    
        // Run the command
        err := cmd.Run()
        if err != nil {
            logArea.SetText(logArea.Text + "Flashing failed: " + err.Error() + "\n")
        } else {
            logArea.SetText(logArea.Text + "Flashing completed successfully.\n")
        }
    })

    // File selection buttons layout
    fileSelectionButtons := container.NewVBox(
        selectPITButton,
        selectBLButton,
        selectAPButton,
        selectCPButton,
        selectCSCButton,
    )

    // Flash button layout
    flashButtonContainer := container.NewHBox(
        widget.NewLabel(""),
        flashButton,
        widget.NewLabel(""),
    )

    // General Samsung Buttons layout
    samsungButtons := container.NewHBox(
        samsungButton1,
        samsungButton2,
        samsungButton3,
        samsungButton4,
        samsungButton5,
    )

    // Samsung Tab layout
    return container.NewVBox(
        samsungButtons,         // General Samsung buttons
        fileSelectionButtons,   // File selection buttons
        flashButtonContainer,   // Flash button
    )
}