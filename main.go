package main

import (
	"machine"
	"strconv"
	"time"
)

const StartHeader = 0x01
const StartText = 0x02
const EndText = 0x03
const EndTransmission = 0x04
const EscapeChar = 0x1B

var receiveChannel = make(chan WireTransmission, 64)
var sendChannel = make(chan WireTransmission, 64)

// sendCache remembers what we sent in case the transmission fails
var sendCache = make([]*WireTransmission, 128)
var sendCacheIndex = 0

var isStreamingBody = false
var isStreamingHeader = false
var isStreamingChecksum = false
var isEscaping = false
var currentHeader []byte
var currentBody []byte
var currentChecksum []byte

var isLightOn = false
var lightOffTime = time.Now()

var LastPacketReceived = time.Now()

func ErrSerialLog(err error) {
	if err != nil {
		_, _ = machine.Serial.Write([]byte(err.Error()))
	}
}

// need to check why it hangs and probably kill hung packets after a timeout

func main() {
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ErrSerialLog(machine.Serial.Configure(machine.UARTConfig{BaudRate: 115200}))
	ErrSerialLog(machine.UART0.Configure(machine.UARTConfig{BaudRate: 115200, RX: machine.UART0_RX_PIN, TX: machine.UART0_TX_PIN}))

	// Run loop forever
	for {
		// todo: add print debugging to http (/dev/ttyACM0 is the serial port for the pi when the button is held down on the UDOO key)

		if isLightOn {
			if time.Since(lightOffTime) > 0 {
				machine.LED.Low()
				isLightOn = false
			}
		}

		// RECEIVING CODE
		for machine.UART0.Buffered() > 0 {
			data, err := machine.UART0.ReadByte()

			LastPacketReceived = time.Now()

			if err == nil {
				HandleReceiveByte(data)
			}
		}

		if time.Since(LastPacketReceived) > 50*time.Millisecond && IsStreaming() {
			// might as well see if we can get any data out of it, cuz the checksum will tell us if it's corrupted
			HandleReceiveByte(EndTransmission)
		}

		// SENDING CODE
		for len(sendChannel) > 0 {
			wt := <-sendChannel
			toSend, err := wt.Serialize()

			if err != nil {
				ErrSerialLog(err)
			} else {
				for _, b := range toSend {
					ErrSerialLog(machine.UART0.WriteByte(b))
				}
			}
		}

		// PROCESSING CODE
		for len(receiveChannel) > 0 {
			msg := <-receiveChannel

			// process rebroadcast requests
			if msg.Headers.Get("type") == "requestRebroadcast" {
				intID, err := strconv.Atoi(string(msg.Body))

				if err != nil {
					ErrSerialLog(Rebroadcast(intID))
				}
			}

			// process other messages
			HandleMessage(msg)
		}
	}
}
