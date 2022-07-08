package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/effects"
	
	"os"
	"log"
	"fmt"
	"image/color"
	"time"
)

var TextColors = map[string]color.RGBA{
  "green": color.RGBA{85, 165, 34, 1},
  // "grey": color.RGBA{82, 82, 82, 1},
  // "white": color.RGBA{255, 255, 255, 1},
  // "lightgrey": color.RGBA{57, 57, 57, 255},
  // "lightgrey2": color.RGBA{142, 142, 142, 255},
}

const APP_NAME = "Pomadorik"
const APP_WIDTH = 250
const APP_HEIGHT = 250
const SOUND_FILE = "click1.mp3"

// pause name: seconds
var DEFAULT_TIMERS = map[string]int{ 
	"TOMATO": 3, // 1200 sec = 20 min
	"SHORT": 300,
	"LONG": 600,
}

var TIMER = DEFAULT_TIMERS["TOMATO"] 
var TICKER *time.Ticker = nil

type BtnHandlerFn func(string, *canvas.Text) func()
var mainWindow fyne.Window

func main() {
	app := app.NewWithID(APP_NAME)
	mainWindow = app.NewWindow(APP_NAME)
	mainWindow.Resize(fyne.NewSize(APP_WIDTH, APP_HEIGHT))
	fmt.Println("window init...")

	content := buildContent(func (timerName string, timerTxt *canvas.Text) func() {
		// set on "space" start a tomato timer
		// https://developer.fyne.io/api/v1.4/keyname.html
		mainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			switch k.Name {
			case fyne.KeySpace:
				startCountdown(DEFAULT_TIMERS["TOMATO"], timerTxt)
			}
		})

		return func() {
			startCountdown(DEFAULT_TIMERS[timerName], timerTxt)
		}
	})

	mainWindow.SetContent(content)
	mainWindow.Show()

	// mainWindow.ShowAndRun()

	app.Run()
}

func showNotification(a fyne.App) {
	time.Sleep(time.Second * 2)
	a.SendNotification(fyne.NewNotification("Example Title", "Example Content"))
}

func buildContent(onBtnHandler BtnHandlerFn) fyne.CanvasObject {
	greenColor := TextColors["green"]

	timer := buildTxtWithStyle(formatTimer(TIMER), greenColor, 40)
	tomatoBtn := widget.NewButton("Tomato", onBtnHandler("TOMATO", timer))
	shortBrakeBtn := widget.NewButton("Short brake", onBtnHandler("SHORT", timer))
	longBrakeBtn := widget.NewButton("Long brake", onBtnHandler("LONG", timer))

	content := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),

		// header (timer)
		container.New(layout.NewCenterLayout(), timer),

		// btns 
		tomatoBtn,
		buildSpace(),

		shortBrakeBtn,
		longBrakeBtn,
	)
	return content
}



func buildTxtWithStyle(title string, textColor color.RGBA, textSize float32) *canvas.Text {
	txt := canvas.NewText(title, textColor)
	txt.TextSize = textSize
	// txt.Alignment = fyne.TextAlignTrailing 
	return txt
}

func buildLabelTxt(title string) *canvas.Text {
	txt := canvas.NewText(title, TextColors["grey"])
	txt.TextSize = 12
	return txt
}

func buildSpace() *canvas.Text {
	return buildLabelTxt("")
}

func updateTimerTxt(timer int, timerTxt *canvas.Text) {
	timerTxt.Text = formatTimer(timer) 
	timerTxt.Refresh() 
}

func startCountdown(defaultTime int, timerTxt *canvas.Text) {
		// if timer already started, at again start, just stop it 
		TIMER = defaultTime
		updateTimerTxt(TIMER, timerTxt)

		if TICKER != nil {
			TICKER.Stop()
		}

		ticker := startTimer(func (ticker *time.Ticker) {
			TIMER--
			if TIMER == 0 {
				playSound()
				ticker.Stop()
				TICKER = nil
				mainWindow.RequestFocus()
			}

			updateTimerTxt(TIMER, timerTxt)
		})
		TICKER = ticker
}

// https://gobyexample.com/tickers
func startTimer(onTickFn func(*time.Ticker)) *time.Ticker {
	ticker := time.NewTicker(1 * time.Second)
	done := make(chan bool)
	go func() {
		for {
			select {
				case <-done: return
				case <-ticker.C: 
					onTickFn(ticker)
			}
		}
	}()
	return ticker
}

func playSound() {
	f, err := os.Open("./sounds/" + SOUND_FILE)
	if err != nil {
		log.Fatal("Unable to open sound " + SOUND_FILE)
	}

	stream, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal("Unable to stream sound " + SOUND_FILE)
	}

	volume := effects.Volume{ 
		Streamer: stream,
		Base: 2,
		Volume: 1.8,
		Silent: false,
	}

	// activate speakers 
	speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)

	// play
	speaker.Play(&volume) 
}

func formatTimer(timer int) string {
	minutes := TIMER / 60
	seconds := TIMER % 60
	minZero := ""
	secZero := ""

	if minutes < 10 {
		minZero = "0"
	}
	if seconds < 10 {
		secZero = "0"
	}
	return fmt.Sprintf("%s%d:%s%d", minZero, minutes, secZero, seconds)
}