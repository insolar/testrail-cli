//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package types

import (
	"github.com/insolar/testrail-cli/parser"
)

type Converter interface {
	// ConvertEventsToMatcherObjectsPreload parses event batches to construct TestObjects, extracting caseID, Description, Status and IssueURL
	ConvertEventsToMatcherObjectsPreload(events map[string][]parser.TestEvent) []*TestMatcher
	// ConvertEventsToMatcherObjects parses event stream to construct TestObject, extracting caseID, Description, Status and IssueURL
	ConvertEventsToMatcherObjects(reader parser.EventReader)[]*TestMatcher
}

// TestMatcher represents data differences between implementation and testrail case
type TestMatcher struct {
	ID                  int
	Status              string
	Description         string
	OriginalDescription string
	GoTestName          string
	IssueURL            string
}
