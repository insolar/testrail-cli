//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package convlog

import (
	"bufio"
	"io"
)

// LineScanner scans lines and keep track of line numbers
type LineScanner struct {
	*bufio.Reader
	lnum int    // Current line number.
	text []byte // Content of current line of text.
	err  error  // Error from latest operation.
}

// NewLineScanner creates a new line scanner from r
func NewLineScanner(r io.Reader) *LineScanner {
	br := bufio.NewReader(r)
	ls := &LineScanner{
		Reader: br,
	}
	return ls
}

// Scan advances to next line.
func (ls *LineScanner) Scan() bool {
	if ls.text, ls.err = ls.Reader.ReadBytes('\n'); ls.err != nil {
		if ls.err == io.EOF {
			ls.err = nil
		}
		return false
	}
	ls.lnum++
	return true
}

// Text returns the current line
func (ls *LineScanner) Text() []byte {
	return ls.text
}

// Err returns the current error (nil if no error)
func (ls *LineScanner) Err() error {
	return ls.err
}

// Line returns the current line number
func (ls *LineScanner) Line() int {
	return ls.lnum
}
