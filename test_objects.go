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

// TestMatcher represents data differences between implementation and testrail case
type TestMatcher struct {
	Status     string
	CaseID     int
	Desc       string
	TRDesc     string
	GoTestName string
	IssueURL   string
}

type TestObjectSummary struct {
	Valid          []*TestMatcher
	NotFound       []*TestMatcher
	WrongDesc      []*TestMatcher
	SkippedNoIssue []*TestMatcher
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
func (m *TestRail) EventsToTestObjects(events map[string][]*TestEvent) []*TestMatcher {
	tests := make([]*TestMatcher, 0)
	for _, eventsBatch := range events {
		t := &TestMatcher{}
		for _, e := range eventsBatch {
			if e.Action == "output" {
				t.GoTestName = e.Test
				res := testCaseIdRe.FindAllStringSubmatch(e.Output, -1)
				if len(res) != 0 && len(res[0]) == 3 {
					d, err := strconv.Atoi(res[0][1])
					if err != nil {
						log.Fatal(err)
					}
					//TODO: Bad, should harden regex instead, debatable for now
					if t.CaseID != 0 {
						continue
					}
					t.CaseID = d
					t.Desc = res[0][2]
				}
				res = testStatusRe.FindAllStringSubmatch(e.Output, -1)
				if len(res) != 0 && len(res[0]) == 2 {
					t.Status = res[0][1]
				}
				res = testSkipIssueRe.FindAllStringSubmatch(e.Output, -1)
				if len(res) != 0 && len(res[0]) == 2 {
					t.IssueURL = res[0][1]
				}
			}
		}
		tests = append(tests, t)
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
			if UniqTestKey(e) == uniqTest {
				eventsByTest[uniqTest] = append(eventsByTest[uniqTest], e)
			}
		}
	}
	return eventsByTest
}

// JSONEventsToSendable convert events to test rail sendable format
func (m *TestRail) JSONEventsToSendable(events []*TestEvent) *testrail.SendableResultsForCase {
	testEventsBatch := m.GroupEventsByTest(events)
	tObjects := m.EventsToTestObjects(testEventsBatch)
	return m.TestObjectsToSendableResultsForCase(tObjects)
}

// TestObjectsToSendableResultsForCase converts TestObjects to sendable results
func (m *TestRail) TestObjectsToSendableResultsForCase(objs []*TestMatcher) *testrail.SendableResultsForCase {
	if len(objs) == 0 {
		log.Println("no valid tests found matching cases, skip sending")
		return nil
	}
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
	return &sendableResults
}

// NAResults generate payload to set all test rail run results to UNTESTED
func (m *TestRail) NAResults(cases []*CaseWithDesc) *testrail.SendableResultsForCase {
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
	sendableResults := &testrail.SendableResultsForCase{
		Results: results,
	}
	return sendableResults
}

func LogInvalidTests(objs *TestObjectSummary) {
	if len(objs.NotFound) > 0 {
		log.Println("Tests without testrail case ID:")
		for _, o := range objs.NotFound {
			log.Printf("  %s", o.GoTestName)
		}
	}
	if len(objs.WrongDesc) > 0 {
		log.Println("Test title discrepancy with testrail test-case title:")
		for _, o := range objs.WrongDesc {
			log.Printf("  %s", o.GoTestName)
			log.Printf("    Test: %s", o.Desc)
			log.Printf("    Testrail: %s", o.TRDesc)
			log.Printf("    Testcase: %s", TRTicket(o.CaseID))
		}
	}
	if len(objs.SkippedNoIssue) > 0 {
		log.Println("Skipped tests without issue:")
		for _, o := range objs.SkippedNoIssue {
			log.Printf("  %s", o.GoTestName)
		}
	}
}

// FilterTestObjects split test objects into groups: valid/not found/wrong description
func FilterTestObjects(objs []*TestMatcher, cases []*CaseWithDesc) *TestObjectSummary {
	wrongDescObjs := make([]*TestMatcher, 0)
	skipNoIssue := make([]*TestMatcher, 0)
	notFoundObjs := make([]*TestMatcher, 0)
	validObjs := make([]*TestMatcher, 0)
	for _, o := range objs {
		found := false
		for _, c := range cases {
			if o.CaseID == c.CaseID {
				found = true
				if o.Desc != c.Desc {
					o.TRDesc = c.Desc
					wrongDescObjs = append(wrongDescObjs, o)
					continue
				}
				validObjs = append(validObjs, o)
			}
		}
		if !found {
			notFoundObjs = append(notFoundObjs, o)
		}
		if o.Status == "SKIP" && o.IssueURL == "" {
			skipNoIssue = append(skipNoIssue, o)
		}
	}
	return &TestObjectSummary{
		Valid:          validObjs,
		NotFound:       notFoundObjs,
		WrongDesc:      wrongDescObjs,
		SkippedNoIssue: skipNoIssue,
	}
}
