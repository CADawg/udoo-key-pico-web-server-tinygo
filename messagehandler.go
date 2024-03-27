package main

import (
	"machine"
)

func HandleMessage(message WireTransmission) {
	TurnOnLed() // so we can see transmissions

	switch message.Headers.Get("type") {
	// TODO: Implement message handling
	case "ping":
		wt := WireTransmission{
			Headers: Headers{
				{"type", "pong"},
			},
		}

		_ = SendMessage(wt)
	case "http":
		url := message.Headers.Get("url")

		_, _ = machine.Serial.Write([]byte("\n\rHTTP GET: " + url))

		file, err := HttpGetFile(url)

		_, _ = machine.Serial.Write([]byte("\n\rHTTP GET: " + url + " done"))

		if err != nil {
			_, _ = machine.Serial.Write([]byte("\n\rHTTP GET: " + url + " error: " + err.Error() + "\n\n"))

			wt := WireTransmission{
				Headers: Headers{
					{"type", "http"},
					{"url", url},
					{"status", "error"},
					{"responseTo", message.Headers.Get("id")},
				},
				Body: []byte(err.Error()),
			}

			_, _ = machine.Serial.Write([]byte("\n\rHTTP GET: " + url + " error response"))

			err = SendMessage(wt)

			if err != nil {
				_, _ = machine.Serial.Write([]byte("\n\rHTTP GET: " + url + " error response error: " + err.Error()))
			}
		} else {
			_, _ = machine.Serial.Write([]byte("\n\rHTTP GET: " + url + " success"))

			wt := WireTransmission{
				Headers: Headers{
					{"type", "http"},
					{"url", url},
					{"status", "ok"},
					{"responseTo", message.Headers.Get("id")},
				},
				Body: file,
			}

			_, _ = machine.Serial.Write([]byte("HTTP GET: " + url + " success response\n"))

			err = SendMessage(wt)

			if err != nil {
				_, _ = machine.Serial.Write([]byte("HTTP GET: " + url + " success response error: " + err.Error() + "\n"))
			}

			_, _ = machine.Serial.Write([]byte("HTTP GET: " + url + " success response sent\n"))
		}
	}
	return
}
