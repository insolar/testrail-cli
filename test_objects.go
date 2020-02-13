// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package testrail_cli

import (
	"github.com/educlos/testrail"
	"log"
	"strconv"
	"time"
)

// TestObject represents data required for testrail run result
type TestObject struct {
	Status     string
	CaseID     int
	Desc       string
	GoTestName string
	IssueURL   string
}

// TestEvent go test2json event object
type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

// EventsToTestObjects parses event batches to construct TestObjects, extracting caseID, Description, Status and IssueURL
func (m *TestRail) EventsToTestObjects(events map[string][]*TestEvent) []*TestObject {
	tests := make([]*TestObject, 0)
	for _, eventsBatch := range events {
		t := &TestObject{}
		for _, e := range eventsBatch {
			if e.Action == "output" {
				t.GoTestName = e.Test
				res := testCaseIdRe.FindAllStringSubmatch(e.Output, -1)
				if len(res) != 0 && len(res[0]) == 3 {
					d, err := strconv.Atoi(res[0][1])
					if err != nil {
						log.Fatal(err)
					}
					t.CaseID = d
					t.Desc = res[0][2]
				}
				res = testStatusRe.FindAllStringSubmatch(e.Output, -1)
				if len(res) != 0 && len(res[0]) == 2 {
					t.Status = res[0][1]
				}
				res = testIssueRe.FindAllStringSubmatch(e.Output, -1)
				if len(res) != 0 && len(res[0]) == 2 {
					t.IssueURL = res[0][1]
				}
			}
		}
		if t.CaseID != 0 && t.Status != "" {
			tests = append(tests, t)
		}
	}
	return tests
}

func UniqTestKey(e *TestEvent) string {
	return e.Test + "|" + e.Package
}

// GroupEventsByTest groups test2json events by test + package key
func (m *TestRail) GroupEventsByTest(events []*TestEvent) map[string][]*TestEvent {
	eventsByTest := make(map[string][]*TestEvent)
	testNames := make(map[string]int)
	for _, e := range events {
		if _, ok := testNames[UniqTestKey(e)]; !ok {
			testNames[UniqTestKey(e)] = 1
		}
	}
	for uniqTest := range testNames {
		for _, e := range events {
			if UniqTestKey(e) == uniqTest && e.Test != "" {
				eventsByTest[uniqTest] = append(eventsByTest[uniqTest], e)
			}
		}
	}
	return eventsByTest
}

// JSONEventsToSendable convert events to test rail sendable format
func (m *TestRail) JSONEventsToSendable(events []*TestEvent) testrail.SendableResultsForCase {
	testEventsBatch := m.GroupEventsByTest(events)
	tObjects := m.EventsToTestObjects(testEventsBatch)
	return m.TestObjectsToSendableResultsForCase(tObjects)
}

// TestObjectsToSendableResultsForCase converts TestObjects to sendable results
func (m *TestRail) TestObjectsToSendableResultsForCase(objs []*TestObject) testrail.SendableResultsForCase {
	results := make([]testrail.ResultsForCase, 0)
	for _, o := range objs {
		result := testrail.ResultsForCase{
			CaseID: o.CaseID,
			SendableResult: testrail.SendableResult{
				AssignedToID: autotestUserID,
				StatusID:     statusMap[o.Status],
				Comment:      "",
				Version:      "1",
				Elapsed:      *testrail.TimespanFromDuration(1 * time.Second),
				Defects:      TicketFromURL(o.IssueURL),
			},
		}
		results = append(results, result)
	}
	sendableResults := testrail.SendableResultsForCase{
		Results: results,
	}
	return sendableResults
}

// NAResults generate payload to set all test rail run results to UNTESTED
func (m *TestRail) NAResults(cases []*CaseWithDesc) testrail.SendableResultsForCase {
	results := make([]testrail.ResultsForCase, 0)
	for _, c := range cases {
		result := testrail.ResultsForCase{
			CaseID: c.CaseID,
			SendableResult: testrail.SendableResult{
				AssignedToID: autotestUserID,
				StatusID:     statusMap["N/A"],
				Comment:      "",
				Version:      "1",
				Elapsed:      *testrail.TimespanFromDuration(1 * time.Second),
				Defects:      "",
			},
		}
		results = append(results, result)
	}
	sendableResults := testrail.SendableResultsForCase{
		Results: results,
	}
	return sendableResults
}

// FilterValidTests add only tests with matching case id and description
func FilterValidTests(objs []*TestObject, cases []*CaseWithDesc) []*TestObject {
	filteredObjs := make([]*TestObject, 0)
	for _, o := range objs {
		found := false
		for _, c := range cases {
			if o.CaseID == c.CaseID {
				found = true
				if o.Desc != c.Desc {
					log.Printf(
						"case description doesn't match, N/A status will be sent:\ntest: %s\nhas: %s\nwant: %s\n",
						o.GoTestName,
						o.Desc,
						c.Desc,
					)
					continue
				}
				filteredObjs = append(filteredObjs, o)
			}
		}
		if !found {
			log.Printf("case not found: caseID: %d title: %s", o.CaseID, o.Desc)
		}
	}
	return filteredObjs
}
