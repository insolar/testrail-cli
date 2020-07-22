//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package convlog

import (
	"io"
	"log"

	"github.com/insolar/testrail-cli/parser"
)

var _ parser.Parser = (*Parser)(nil)

type Parser struct{}

func (p Parser) Parse(input io.Reader) []parser.TestEvent {
	var testEvents []parser.TestEvent

	iter := p.GetParseIterator(input)
	for {
		_, te, err := iter.Next()
		if parser.IsEOF(err) {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		testEvents = append(testEvents, te)
	}

	return testEvents
}

func (Parser) GetParseIterator(input io.Reader) parser.EventReader {
	return &iterativeReader{scanner: NewLineScanner(input)}
}