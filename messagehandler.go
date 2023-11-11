package main

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

		file, err := HttpGetFile(url)

		if err != nil {
			wt := WireTransmission{
				Headers: Headers{
					{"type", "http"},
					{"url", url},
					{"status", "error"},
					{"responseTo", message.Headers.Get("id")},
				},
				Body: err.Error(),
			}

			_ = SendMessage(wt)
		} else {
			wt := WireTransmission{
				Headers: Headers{
					{"type", "http"},
					{"url", url},
					{"status", "ok"},
					{"responseTo", message.Headers.Get("id")},
				},
				Body: string(file),
			}

			_ = SendMessage(wt)
		}
	}
	return
}
