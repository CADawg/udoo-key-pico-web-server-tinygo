package main

import (
	"fmt"
	"strconv"
)

func IsStreaming() bool {
	return isStreamingHeader || isStreamingBody || isStreamingChecksum
}

func GetAvailableID() int {
	sendCacheIndex += 1
	if sendCacheIndex >= len(sendCache) {
		sendCacheIndex = 0
	}

	// wipe old data
	sendCache[sendCacheIndex] = nil

	// we always return the next number because we'll never know if the transmission failed (otherwise we'd use the entire bandwidth for acknowledgements of acknowledgements)
	return sendCacheIndex
}

func Rebroadcast(id int) error {
	if id < 0 || id >= len(sendCache) {
		return fmt.Errorf("invalid id")
	}

	if sendCache[id] == nil {
		return fmt.Errorf("no transmission with that id")
	}

	sendChannel <- *sendCache[id]

	return nil
}

func RequestRebroadcast(id int) error {
	wt := WireTransmission{
		Headers: Headers{
			{"type", "requestRebroadcast"},
		},
		Body: []byte(strconv.Itoa(id)),
	}

	return SendMessage(wt)
}

func SendMessage(wt WireTransmission) error {
	// get ID
	id := GetAvailableID()

	// if we couldn't get an ID, just ignore it (if we lose this packet it's just never going to be answered)
	if id != -1 {
		wt.Headers = append(wt.Headers, Header{"id", strconv.Itoa(id)})
	}

	// check the headers are valid
	_, err := wt.Headers.Serialize()

	if err != nil {
		return err
	}

	sendChannel <- wt

	// add to cache (evicting oldest)
	sendCache[id] = &wt

	return nil
}

func HandleReceiveByte(data byte) {
	if data == StartHeader && !isEscaping {
		// set header to streaming and body to not streaming
		isStreamingHeader = true
		isStreamingBody = false
		isStreamingChecksum = false
		// reset the current header
		currentHeader = []byte{}
	} else if data == StartText && !isEscaping {
		// set the body to streaming
		isStreamingBody = true
		isStreamingHeader = false
		isStreamingChecksum = false
		// reset the current body
		currentBody = []byte{}
	} else if data == EndText && isStreamingBody && !isEscaping {
		isStreamingBody = false
		isStreamingHeader = false
		isStreamingChecksum = true

		// reset the current checksum
		currentChecksum = []byte{}
	} else if data == EndTransmission && !isEscaping {
		// set the body to not streaming
		isStreamingBody = false
		isStreamingHeader = false
		isStreamingChecksum = false

		// casting to string should be enough to utf-8 decode the bytes
		wt := WireTransmission{}
		err := wt.Deserialize(string(currentHeader), currentBody, currentChecksum)

		if err != nil {
			ErrSerialLog(err)
		} else {
			receiveChannel <- wt
		}
	} else if data == EscapeChar && !isEscaping {
		isEscaping = true
	} else if isStreamingHeader {
		// Add the next byte to the current header
		currentHeader = append(currentHeader, data)
		// reset the escaping flag
		isEscaping = false
	} else if isStreamingBody {
		// Add the next byte to the current body
		currentBody = append(currentBody, data)
		// reset the escaping flag
		isEscaping = false
	} else if isStreamingChecksum {
		currentChecksum = append(currentChecksum, data)
		isEscaping = false
	}
}
