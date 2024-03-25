package main

import "machine"

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

		machine.Serial.Write([]byte("HTTP GET: " + url + "\n"))

		file, err := HttpGetFile(url)

		machine.Serial.Write([]byte("HTTP GET: " + url + " done\n"))

		if err != nil {
			machine.Serial.Write([]byte("HTTP GET: " + url + " error: " + err.Error() + "\n"))

			wt := WireTransmission{
				Headers: Headers{
					{"type", "http"},
					{"url", url},
					{"status", "error"},
					{"responseTo", message.Headers.Get("id")},
				},
				Body: err.Error(),
			}

			machine.Serial.Write([]byte("HTTP GET: " + url + " error response\n"))

			err = SendMessage(wt)

			if err != nil {
				machine.Serial.Write([]byte("HTTP GET: " + url + " error response error: " + err.Error() + "\n"))
			}
		} else {
			machine.Serial.Write([]byte("HTTP GET: " + url + " success\n"))

			wt := WireTransmission{
				Headers: Headers{
					{"type", "http"},
					{"url", url},
					{"status", "ok"},
					{"responseTo", message.Headers.Get("id")},
				},
				Body: string(file),
			}

			machine.Serial.Write([]byte("HTTP GET: " + url + " success response\n"))

			err = SendMessage(wt)

			if err != nil {
				machine.Serial.Write([]byte("HTTP GET: " + url + " success response error: " + err.Error() + "\n"))
			}

			machine.Serial.Write([]byte("HTTP GET: " + url + " success response sent\n"))
		}
	}
	return
}
