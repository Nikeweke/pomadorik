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

	"github.com/getlantern/systray"
	// "github.com/skratchdot/open-golang/open"
	// "github.com/getlantern/systray/example/icon"
	"pomadorik/icon"

	"io/ioutil"
	"path/filepath"
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

var TIMER = DEFAULT_TIMERS["TOMATO"] 
var TICKER *time.Ticker = nil
var TIMER_TXT *canvas.Text = nil 

type BtnHandlerFn func(string, *canvas.Text) func()
var mainWindow fyne.Window
var App fyne.App 

func main() {
	App = app.NewWithID(APP_NAME)

	// setuping window 
	mainWindow = App.NewWindow(APP_NAME)
	mainWindow.Resize(fyne.NewSize(APP_WIDTH, APP_HEIGHT))
	mainWindow.SetMaster()

	// set icon 
	r, _ := LoadResourceFromPath("./icon/app-icon.png")
	mainWindow.SetIcon(r)

	// Register systray's starting and exiting functions
	systray.Register(onReady, onExit)

	// setting intercept not to close app, but hide window,
	// and close only via tray 
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	content := buildContent(func (timerName string, timerTxt *canvas.Text) func() {
		TIMER_TXT = timerTxt

		// set on "space" start a tomato timer
		// https://developer.fyne.io/api/v1.4/keyname.html
		mainWindow.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
			switch k.Name {
			case fyne.KeySpace:
				startCountdown(DEFAULT_TIMERS["TOMATO"])
			}
		})

		return func() {
			startCountdown(DEFAULT_TIMERS[timerName])
		}
	})

	mainWindow.SetContent(content)
	fmt.Println("window init...")

	mainWindow.Show()
	App.Run()
}

// // will show desktop notification
// func showNotification(a fyne.App) {
// 	time.Sleep(time.Second * 2)
// 	a.SendNotification(fyne.NewNotification("Example Title", "Example Content"))
// }


// ==============================================> SYSTRAY
// https://pkg.go.dev/github.com/getlantern/systray
func onExit() {} 
func onReady() {
	// Set up menu
	systray.SetTemplateIcon(icon.Data, icon.Data)

	systray.SetTitle(APP_NAME)
	systray.SetTooltip(APP_NAME)

	mTimer := systray.AddMenuItem("Open", "Open") // returns *MenuItem and has title 
	systray.AddSeparator()

	mTomato := systray.AddMenuItem("Tomato", "Starts timer") // title, tooltip
	mShort := systray.AddMenuItem("Short break", "Starts timer of short break")
	mLong := systray.AddMenuItem("Long break", "Starts timer of long break")
	systray.AddSeparator()

	mQuit := systray.AddMenuItem("Quit", "Quit this")

	// Display Menu
	go func() {
		for {
			select {
			case <-mTimer.ClickedCh:
				mainWindow.Show()

			case <-mTomato.ClickedCh:
				startCountdown(DEFAULT_TIMERS["TOMATO"])
			
			case <-mShort.ClickedCh:
				startCountdown(DEFAULT_TIMERS["SHORT"])

			case <-mLong.ClickedCh:
				startCountdown(DEFAULT_TIMERS["LONG"])

			case <-mQuit.ClickedCh:
				// systray.Quit()
				App.Quit()
				return
			}
		}
	}()
}

// ==============================================> ICON
type Resource interface {
	Name() string
	Content() []byte
}
type StaticResource struct {
	StaticName    string
	StaticContent []byte
}
func (r *StaticResource) Name() string {
	return r.StaticName
}
func (r *StaticResource) Content() []byte {
	return r.StaticContent
}

func LoadResourceFromPath(path string) (Resource, error) {
	bytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
			return nil, err
	}
	name := filepath.Base(path)
	return NewStaticResource(name, bytes), nil
}

func NewStaticResource(name string, content []byte) *StaticResource {
	return &StaticResource{
			StaticName:    name,
			StaticContent: content,
	}
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
		buildSpace(),

		container.New(
			layout.NewCenterLayout(), 
			// container.New(layout.NewVBoxLayout(),
				buildTxtWithStyle(
					"Press \"Space\" to start Tomato",
					TextColors["grey"],
					10,
				),
			// ),

		),
			
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
			
	systray.SetTitle(fmt.Sprintf("%s (%s)", APP_NAME, timerTxt.Text))
	systray.SetTooltip(fmt.Sprintf("%s (%s)", APP_NAME, timerTxt.Text))
}

func startCountdown(defaultTime int) {
	if TICKER != nil {
		TICKER.Stop()
	}

	// if timer already started, at again start, just stop it 
	TIMER = defaultTime
	updateTimerTxt(TIMER, TIMER_TXT)

	TICKER = startTimer(func (ticker *time.Ticker) {
		updateTimerTxt(TIMER, TIMER_TXT)

		if TIMER == 0 {
			playSound()
			ticker.Stop()
			TICKER = nil
			mainWindow.Show()
			mainWindow.RequestFocus()
		}

		TIMER--
	})
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
		Volume: 1.6,
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