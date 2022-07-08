package main 

const APP_NAME = "Pomadorik"

const APP_WIDTH = 250
const APP_HEIGHT = 250
const SOUND_FILE = "click1.mp3"

// pause name: seconds
var DEFAULT_TIMERS = map[string]int{ 
	"TOMATO": 1200, // 1200 sec = 20 min
	"SHORT": 300,
	"LONG": 600,
}
