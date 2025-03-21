package main

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

// Function to create the Samsung Tab
func createSamsungTab(logArea *widget.Entry, clearLog func()) fyne.CanvasObject {
    // Samsung Buttons
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

    // Samsung Tab layout
    return container.NewVBox(
        samsungButtons, // Buttons side by side
    )
}