//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package testrail

import (
	"log"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/educlos/testrail"
	"github.com/spf13/viper"

	"github.com/insolar/testrail-cli/types"
)

var (
	autotestUserID = 10
	statusMap      = map[string]int{
		types.TestStatusPassed:       testrail.StatusPassed,
		types.TestStatusFailed:       testrail.StatusFailed,
		types.TestStatusSkipped:      6,
		types.TestStatusNotAvailable: 7,
	}
)


func TicketFromURL(url string) string {
	if strings.HasPrefix(url, "https") || strings.HasPrefix(url, "http") {
		s := strings.Split(url, "/")
		return s[len(s)-1]
	}
	return url
}

type Uploader struct {
	c   *testrail.Client
	run testrail.Run

	runID        int
	tests        map[int]testrail.SendableResult
	defaultTests types.TestCasesWithDescription
}

func NewUploader(url string, user string, password string) *Uploader {
	return &Uploader{
		c:     testrail.NewClient(url, user, password),
		tests: make(map[int]testrail.SendableResult),
	}
}

func (m Uploader) FormatURL(id int) string {
	return path.Join(viper.GetString("URL"), "/index.php?/cases/view/", strconv.Itoa(id))
}

func (m *Uploader) getCasesWithDescription(projectID int, suiteID int) types.TestCasesWithDescription {
	cases, err := m.c.GetCases(projectID, suiteID)
	if err != nil {
		log.Fatal(err)
	}

	var casesWithDescription types.TestCasesWithDescription
	for _, c := range cases {
		caseWithDescription := types.TestCaseWithDescription{
			ID:          c.ID,
			Description: c.Title,
		}
		casesWithDescription = append(casesWithDescription, caseWithDescription)
	}
	return casesWithDescription
}

func (m *Uploader) Init(runID int) {
	run, err := m.c.GetRun(runID)
	if err != nil {
		log.Fatal(err)
	}
	m.run = run

	m.runID = runID

	testCasesWithDescription := m.getCasesWithDescription(m.run.ProjectID, m.run.SuiteID)
	for _, testCase := range testCasesWithDescription {
		// update all cases with N/A status, we store all autotests in ONE run, so in case
		// someone delete particular case implementation status must be updated to N/A
		m.tests[testCase.ID] = testrail.SendableResult{
			AssignedToID: autotestUserID,
			StatusID:     statusMap[types.TestStatusNotAvailable],
			Comment:      "",
			Version:      "1",
			Elapsed:      *testrail.TimespanFromDuration(1 * time.Second),
			Defects:      "",
		}
	}
	m.defaultTests = testCasesWithDescription
}

func (m Uploader) GetCasesWithDescription() types.TestCasesWithDescription {
	return m.defaultTests
}

func (m *Uploader) AddTests(objects []*types.TestMatcher, ignoreNonExistent bool) {
	for _, object := range objects {
		if _, ok := m.tests[object.ID]; !ok && ignoreNonExistent {
			continue
		}
		m.tests[object.ID] = testrail.SendableResult{
			AssignedToID: autotestUserID,
			StatusID:     statusMap[object.Status],
			Comment:      "",
			Version:      "1",
			Elapsed:      *testrail.TimespanFromDuration(1 * time.Second),
			Defects:      TicketFromURL(object.IssueURL),
		}
	}
}

func (m *Uploader) Upload() {
	sendableResults := testrail.SendableResultsForCase{}

	for caseID, resultForCase := range m.tests {
		sendableResults.Results = append(sendableResults.Results, testrail.ResultsForCase{
			CaseID:         caseID,
			SendableResult: resultForCase,
		})
	}

	_, err := m.c.AddResultsForCases(m.runID, sendableResults)
	if err != nil {
		log.Fatal(err)
	}
}