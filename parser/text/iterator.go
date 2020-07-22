//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package text

import (
	"fmt"
	"io"

	"github.com/insolar/testrail-cli/parser"
)

type iterativeReader struct {
	scanner   *LineScanner
	buffer    []parser.TestEvent
}

func (i *iterativeReader) popBuffer() (string, parser.TestEvent, error) {
	te := i.buffer[0]
	i.buffer = i.buffer[1:]
	return parser.UniqueTestKeyFromEvent(te), te, nil
}

func (i *iterativeReader) popBufferIfNotEmpty() (string, parser.TestEvent, error) {
	if len(i.buffer) > 0 {
		return i.popBuffer()
	}
	return "", parser.TestEvent{}, io.EOF
}

func (i *iterativeReader) Next() (string, parser.TestEvent, error) {
	if len(i.buffer) > 0 {
		return i.popBuffer()
	}

	converter := pkgConverter{}

	for i.scanner.Scan() {
		text := i.scanner.Text()

		i.buffer = append(i.buffer, converter.handleInputLine(text)...)

		if converter.finished {
			for pos := range i.buffer {
				if i.buffer[pos].Package == "" {
					i.buffer[pos].Package = converter.pkg
				}
			}

			i.buffer = append(i.buffer, converter.flushReport(0)...)
			if len(i.buffer) == 0 {
				continue
			}

			i.buffer = append(i.buffer, parser.TestEvent{
				Action:  converter.result,
				Package: converter.pkg,
				Elapsed: converter.elapsed,
			})

			return i.popBuffer()
		}
	}

	if err := i.scanner.Err(); err != nil {
		return "", parser.TestEvent{}, fmt.Errorf("failed to read text test event: %w", err)
	}

	if converter.result != "" {
		i.buffer = append(i.buffer, converter.flushReport(0)...)
		i.buffer = append(i.buffer, parser.TestEvent{
			Action:  converter.result,
			Package: converter.pkg,
			Elapsed: converter.elapsed,
		})
	}

	return i.popBufferIfNotEmpty()
}