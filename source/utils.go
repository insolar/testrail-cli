//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package source

import (
	"bytes"
	"fmt"
	"strconv"
)

var (
	fourSpaces = []byte("    ")
	testPrefix = []byte("    Test")
)

func TryExtractTest(line []byte) (string, bool) {
	if bytes.HasPrefix(line, testPrefix) && bytes.Contains(line[len(testPrefix):], []byte(":")) {
		// Example:    TestExample: test output
		// we are interested in test name
		// it is not accurate and could capture something wrong
		// but it is better then just assume previous test
		part := line[len(fourSpaces):bytes.Index(line, []byte(":"))]
		return string(part), true
	}
	return "", false
}

func DumpEvents(events []TestEvent) {
	for _, e := range events {
		fmt.Println(`{Action:"` + e.Action + `", Package:"` + e.Package + `", Test:"` + e.Test + `", Elapsed:` + fmt.Sprintf("%f", e.Elapsed) + `, Output:` + strconv.Quote(e.Output) + `},`)
	}
}
