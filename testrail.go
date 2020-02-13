// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package testrail_cli

import (
	"github.com/educlos/testrail"
	"log"
	"regexp"
)

var (
	autotestUserID = 10
	statusMap      = map[string]int{
		"PASS": 1,
		"FAIL": 5,
		"N/A":  7,
		"SKIP": 6,
	}
	testStatusRe = regexp.MustCompile(`--- (.*):`)
	testIssueRe  = regexp.MustCompile(`.*issue:\s(.*?)\s`)
	testCaseIdRe = regexp.MustCompile(`C(\d{1,8})\s(.*)`)
)

type TestRail struct {
	c *testrail.Client
}

func NewTestRail(url string, user string, password string) *TestRail {
	return &TestRail{
		c: testrail.NewClient(url, user, password),
	}
}

type CaseWithDesc struct {
	CaseID int
	Desc   string
}

func (m *TestRail) GetCasesWithDescs(projectID int, suiteID int) []*CaseWithDesc {
	casesWithDescs := make([]*CaseWithDesc, 0)
	cases, err := m.c.GetCases(projectID, suiteID)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range cases {
		casesWithDescs = append(casesWithDescs, &CaseWithDesc{
			CaseID: c.ID,
			Desc:   c.Title,
		})
	}
	return casesWithDescs
}

func (m *TestRail) GetRun(id int) testrail.Run {
	run, err := m.c.GetRun(id)
	if err != nil {
		log.Fatal(err)
	}
	return run
}

func (m *TestRail) UpdateRunForCases(runId int, results testrail.SendableResultsForCase) []testrail.Result {
	res, err := m.c.AddResultsForCases(runId, results)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
