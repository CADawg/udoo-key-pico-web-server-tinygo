package main

func HandleMessage(message WireTransmission) {
	TurnOnLed() // so we can see transmissions

	switch message.Headers.Get("type") {
	// TODO: Implement message handling
	}
	return
}
