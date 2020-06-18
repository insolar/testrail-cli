// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package json

import (
	"bufio"
	"encoding/json"
	"io"
	"log"

	"github.com/insolar/testrail-cli/source"
)

var _ source.Parser = (*Parser)(nil)

type Parser struct{}

func (Parser) Parse(input io.Reader) []source.TestEvent {
	var testEvents []source.TestEvent
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		var te source.TestEvent
		if err := json.Unmarshal(scanner.Bytes(), &te); err != nil {
			log.Fatalf("failed to unmarshal test event json: %s\n", err)
		}
		if te.Action == "output" {
			if testName, ok := source.TryExtractTest([]byte(te.Output)); ok {
				te.Test = testName
			}
		}
		testEvents = append(testEvents, te)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return testEvents
}
