package main

import (
	"fmt"
	"strings"
)

const MagicHeaderStart = "HEADERS\n"

type Header struct {
	Key   string
	Value string
}

func SerializeHeader(header Header) (string, error) {
	if header.Key == "" {
		return "", fmt.Errorf("header key is empty")
	}

	if strings.Contains(header.Key, ":") {
		return "", fmt.Errorf("header key contains colon")
	}

	return header.Key + ":" + header.Value, nil
}

type Headers []Header

func (h *Headers) Get(key string) string {
	for _, header := range *h {
		if header.Key == key {
			return header.Value
		}
	}

	return ""
}

// Serialize - join headers together, split by new line \n
func (h *Headers) Serialize() (string, error) {
	var serialisedHeaders string

	// add magic start chars
	serialisedHeaders += MagicHeaderStart

	for _, header := range *h {
		serialisedHeader, err := SerializeHeader(header)
		if err != nil {
			return "", err
		}

		serialisedHeaders += serialisedHeader + "\n"
	}

	return serialisedHeaders, nil
}

func (h *Headers) Deserialize(data string) error {
	// check if start is correct
	if string(data[:len(MagicHeaderStart)]) != MagicHeaderStart {
		return fmt.Errorf("invalid headers start")
	}

	// remove magic start chars
	data = data[len(MagicHeaderStart):]

	// split by new line
	headers := strings.Split(string(data), "\n")

	// remove last empty header
	headers = headers[:len(headers)-1]

	// iterate over headers
	for _, header := range headers {
		// split by colon
		headerParts := strings.Split(header, ":")

		// if header has more than 1 :, join all the remaining parts (support ipv6 and other : containing headers)
		if len(headerParts) > 2 {
			headerParts = []string{headerParts[0], strings.Join(headerParts[1:], ":")}
		}

		// add header to headers
		*h = append(*h, Header{headerParts[0], headerParts[1]})
	}

	return nil
}
