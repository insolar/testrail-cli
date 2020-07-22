//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package convlog

import (
	"fmt"
	"io"
	"regexp"

	"github.com/insolar/testrail-cli/parser"
)

var (
	convLogPrefixCutter = regexp.MustCompile(`(?m)^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+([+\-]\d{2}:\d{2}|Z)\s.{3}\s)`)
)

type iterativeReader struct {
	scanner   *LineScanner
}

func (i *iterativeReader) Next() (string, parser.TestEvent, error) {
	for i.scanner.Scan() {
		bytes := i.scanner.Text()
		if !convLogPrefixCutter.Match(bytes) {
			continue
		}
		text := string(convLogPrefixCutter.ReplaceAll(bytes, []byte(nil)))
		return "", parser.TestEvent{
			Action: "output",
			Output: text,
		}, nil
	}

	if err := i.scanner.Err(); err != nil {
		return "", parser.TestEvent{}, fmt.Errorf("failed to read text test event: %w", err)
	}

	return "", parser.TestEvent{}, io.EOF
}