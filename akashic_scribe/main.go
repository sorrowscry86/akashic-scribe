package main

import (
	"akashic_scribe/core"
	"akashic_scribe/gui"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

func main() {
	// Use a unique app ID so Preferences API works
	myApp := app.NewWithID("com.voidcat.akashicscribe")
	// Set a temporary icon; replace with a custom one later
	myApp.SetIcon(theme.FyneLogo())
	myWindow := myApp.NewWindow("Akashic Scribe")
	// myApp.Settings().SetTheme(gui.NewAkashicTheme())

	// Initialize the core engine
	scribeEngine := core.NewRealScribeEngine()

	mainLayout := gui.CreateMainLayout(myWindow, scribeEngine)
	myWindow.SetContent(mainLayout)
	myWindow.Resize(fyne.NewSize(1200, 800))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
