//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package json

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/insolar/testrail-cli/parser"
)

type iterativeReader struct {
	scanner *bufio.Scanner
}

func (i *iterativeReader) Next() (string, parser.TestEvent, error) {
	if !i.scanner.Scan() {
		if err := i.scanner.Err(); err != nil {
			log.Fatal(err)
		}
		return "", parser.TestEvent{}, io.EOF
	}

	var te parser.TestEvent

	if err := json.Unmarshal(i.scanner.Bytes(), &te); err != nil {
		return "", parser.TestEvent{}, fmt.Errorf("failed to unmarshal test json event: %w", err)
	}

	if te.Action == "output" {
		if testName, ok := parser.TryExtractTest([]byte(te.Output)); ok {
			te.Test = testName
		}
	}

	return parser.UniqueTestKeyFromEvent(te), te, nil
}

