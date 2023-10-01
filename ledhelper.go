package main

import (
	"machine"
	"time"
)

func TurnOnLed() {
	TurnOnLedForMs(100)
}

func TurnOnLedForMs(ms int) {
	// turn on led
	machine.LED.High()

	// set off time to 100 ms from now
	lightOffTime = time.Now().Add(time.Duration(ms) * time.Millisecond)

	// set light on to true
	isLightOn = true
}
