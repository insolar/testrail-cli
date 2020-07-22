//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package parser

import (
	"io"
)

type Parser interface {
	Parse(io.Reader) []TestEvent
	GetParseIterator(io.Reader) EventReader
}

// TestEvent go test2json event object
type TestEvent struct {
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

func UniqueTestKeyFromEvent(e TestEvent) string {
	return UniqueTestKeyFromFields(e.Package, e.Test)
}

func UniqueTestKeyFromFields(pkgName, testName string) string {
	return pkgName + "|" + testName
}