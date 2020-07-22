//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package parser

import (
	"io"
)

type EventReader interface {
	Next() (string, TestEvent, error)
}

type StreamingEventReader struct {
	events map[string][]TestEvent
	keys   []string

	currentKeyPos int
	currentValPos int
}

func NewStreamingEventReaderFromMap(events map[string][]TestEvent) *StreamingEventReader {
	keys := make([]string, len(events))

	keyPos := 0
	for key, _ := range events {
		keys[keyPos] = key
		keyPos++
	}

	return &StreamingEventReader{
		events: events,
		keys: keys,
		currentKeyPos: 0,
		currentValPos: 0,
	}
}

func (r *StreamingEventReader) Next() (string, TestEvent, error) {
	for {
		if r.currentKeyPos >= len(r.keys) {
			return "", TestEvent{}, io.EOF
		}

		currentKey := r.keys[r.currentKeyPos]
		if r.currentValPos >= len(r.events[currentKey]) {
			r.currentValPos = 0
			r.currentKeyPos++
			continue
		}

		currentEvent := r.events[currentKey][r.currentValPos]
		r.currentValPos++
		return currentKey, currentEvent, nil
	}
}

func IsEOF(err error) bool {
	return err == io.EOF
}