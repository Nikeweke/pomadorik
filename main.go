package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	// "fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"

	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	
	"os"
	"log"
	"fmt"
	"image/color"
	"time"
)

var TextColors = map[string]color.RGBA{
  "green": color.RGBA{85, 165, 34, 1},
  "grey": color.RGBA{82, 82, 82, 1},
  "white": color.RGBA{255, 255, 255, 1},
  "lightgrey": color.RGBA{57, 57, 57, 255},
  "lightgrey2": color.RGBA{142, 142, 142, 255},
}

// 1. 3 buttons + timer - done 
// 2. timer - done 
// 2. tray hide
// 3. sound play - done

const DEFAULT_TIMER = 5
// const DEFAULT_TIMER = 1200 // 1200 sec = 20 min
const DEFAULT_LONG_BREAK = 600
const DEFAULT_SHORT_BREAK = 300 

var TIMER = DEFAULT_TIMER 
var TICKER *time.Ticker = nil
// var TIMER_TXT


func main() {
	app := app.New()
	window := app.NewWindow("Pomadoro")
	window.Resize(fyne.NewSize(300, 300))
	fmt.Println("window init...")

	content := buildContent()

	window.SetContent(content)
	window.ShowAndRun()
}

func buildContent() fyne.CanvasObject {
	greenColor := TextColors["green"]

	timer := buildTxtWithStyle(formatTimer(TIMER), greenColor, 40)
	tomatoBtn := widget.NewButton("Tomato", func() {
		startCountdown(DEFAULT_TIMER, timer)
	})
	shortBrakeBtn := widget.NewButton("Short brake", func() {
		startCountdown(DEFAULT_SHORT_BREAK, timer)
	})
	longBrakeBtn := widget.NewButton("Long brake", func() {
		startCountdown(DEFAULT_LONG_BREAK, timer)
	})

	content := fyne.NewContainerWithLayout(
		// layout.NewGridWrapLayout(fyne.NewSize(200, 200)),
		layout.NewVBoxLayout(),

		container.New(layout.NewCenterLayout(), timer),

		tomatoBtn,
		
		buildLabelTxt(""),

		longBrakeBtn,
		shortBrakeBtn,
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

func startCountdown(defaultTime int, timer *canvas.Text) {
		// fmt.Println("tapped")
		if TICKER != nil {
			TICKER.Stop()
			TIMER = defaultTime
		}

		ticker := startTimer(func (ticker *time.Ticker) {
			TIMER--
			if TIMER == 0 {
				playSound()
				ticker.Stop()
			}
			timer.Text = formatTimer(TIMER) 
			timer.Refresh() 
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
	soundFile := "click1.mp3"
	f, err := os.Open("./sounds/" + soundFile)
	if err != nil {
		log.Fatal("Unable to open sound " + soundFile)
	}

	stream, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal("Unable to stream sound " + soundFile)
	}

	// activate speakers 
	speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)

	// play
	speaker.Play(stream)
}

func formatTimer(timer int) string {
	minutes := TIMER / 60
	seconds := TIMER % 60
	return fmt.Sprintf("%d:%d", minutes, seconds)
}