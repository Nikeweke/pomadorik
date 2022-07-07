package main

import (
	// "fmt"
	// "time"

	"fyne.io/fyne/v2"
	app "fyne.io/fyne/v2/app"
	// container "fyne.io/fyne/v2/container"
	// layout "fyne.io/fyne/v2/layout"
	// widget "fyne.io/fyne/v2/widget"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"github.com/getlantern/systray/example/icon"
)

var thisApp fyne.App
var mainWindow fyne.Window
var sabWindow fyne.Window

const APP_NAME = "Pomadorik"

func main() {
	// Load Config
	// Load Auth tokens
    // Register systray's starting and exiting functions
	systray.Register(onReady, onExit)
	// Set up base GUI
	thisApp = app.NewWithID("Test Application")
	// thisApp.SetIcon(resourceIconPng)
	mainWindow = thisApp.NewWindow("Hidden Main Window")
	mainWindow.Resize(fyne.NewSize(800, 800))
	mainWindow.SetMaster()
	mainWindow.Hide()
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})
	sabWindow = thisApp.NewWindow("SAB Window")
	sabWindow.Resize(fyne.NewSize(640, 480))
	sabWindow.Hide()
	sabWindow.SetCloseIntercept(func() {
		sabWindow.Hide()
	})
	thisApp.Run()
}

func onExit() { 

}

func onReady() {
	// Refresh any expired tokens
	// Set up menu
	systray.SetTemplateIcon(icon.Data, icon.Data)
	
	systray.SetTitle(APP_NAME)
	systray.SetTooltip(APP_NAME)


	mGSM := systray.AddMenuItem("20:00", "Timer") // returns *MenuItem and has title 
	mGSM.Disable()
	systray.AddSeparator()

	mAbout := systray.AddMenuItem("Tomato", "Starts timer") // title, tooltip
	mPrefs := systray.AddMenuItem("Short break", "Starts timer of short break")
	mQuit2 := systray.AddMenuItem("Long break", "Starts timer of long break")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit this")

	// Display Menu
	go func() {
		for {
			select {
			case <-mGSM.ClickedCh:
				sabWindow.Show()
			case <-mPrefs.ClickedCh:
				open.Run("https://vonexplaino.com")
			case <-mAbout.ClickedCh:
				open.Run("https://vonexplaino.com/")
			case <-mQuit2.ClickedCh:
				open.Run("https://vonexplaino.com/")
			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}