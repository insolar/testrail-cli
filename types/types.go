//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package types

type TestCaseWithDescription struct {
	ID          int
	Description string
}

type TestCasesWithDescription []TestCaseWithDescription

const (
	TestStatusPassed       = "PASS"
	TestStatusFailed       = "FAIL"
	TestStatusNotAvailable = "N/A"
	TestStatusSkipped      = "SKIP"
)

type TestServer interface {
	FormatURL(id int) string

	Init(runID int)
	GetCasesWithDescription() TestCasesWithDescription
	AddTests(objects []*TestMatcher)
	Upload()
}

func StatusKnown(status string) bool {
	return status == TestStatusPassed ||
		status == TestStatusFailed ||
		status == TestStatusNotAvailable ||
		status == TestStatusSkipped
}
