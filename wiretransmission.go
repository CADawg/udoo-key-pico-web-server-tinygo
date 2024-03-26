package main

import (
	"fmt"
	"hash/crc32"
	"strconv"
)

var ErrChecksumInvalid = fmt.Errorf("checksum invalid")
var ErrChecksumInvalidRequestedRebroadcast = fmt.Errorf("checksum invalid, requested rebroadcast")

type WireTransmission struct {
	Headers  Headers
	Body     string
	Checksum []byte
}

func (w *WireTransmission) Deserialize(headers string, body string, checksum []byte) error {
	// deserialize headers
	err := w.Headers.Deserialize(headers)

	// check checksum
	if bytesToCrc32(checksum) != crc32.ChecksumIEEE([]byte(headers+body)) {
		id := w.Headers.Get("id")

		if id != "" {
			idInt, err := strconv.Atoi(id)

			if err != nil {
				return err
			}

			err = RequestRebroadcast(idInt)

			if err != nil {
				return err
			}

			return ErrChecksumInvalidRequestedRebroadcast
		}

		return ErrChecksumInvalid
	}

	if err != nil {
		return err
	}

	w.Body = body

	return nil
}

func crc32ToBytes(data []byte) []byte {
	checksum := crc32.ChecksumIEEE(data)

	return []byte{byte(checksum >> 24), byte(checksum >> 16), byte(checksum >> 8), byte(checksum)}
}

func bytesToCrc32(data []byte) uint32 {
	if len(data) != 4 {
		return 0
	}

	return uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
}

func (w *WireTransmission) Serialize() ([]byte, error) {
	// serialize headers
	headersSerialized, err := w.Headers.Serialize()

	w.Checksum = crc32ToBytes([]byte(headersSerialized + w.Body))

	if err != nil {
		return nil, err
	}

	// convert text to byte arrays (can't send int32 over the wire directly)
	byteArrayHeader := []byte(headersSerialized)
	byteArrayBody := []byte(w.Body)

	// send the header
	encoded := []byte{StartHeader}
	for _, b := range byteArrayHeader {
		encoded = append(encoded, EncodeByteSafe(b)...)
	}
	encoded = append(encoded, StartText)
	for _, b := range byteArrayBody {
		encoded = append(encoded, EncodeByteSafe(b)...)
	}
	encoded = append(encoded, EndText)
	for _, b := range w.Checksum {
		encoded = append(encoded, EncodeByteSafe(b)...)
	}
	encoded = append(encoded, EndTransmission)

	return encoded, nil
}

func EncodeByteSafe(b byte) []byte {
	if b == StartHeader || b == StartText || b == EndText || b == EndTransmission || b == EscapeChar {
		return []byte{EscapeChar, b}
	}

	return []byte{b}
}
