//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package qase

import (
	"context"
	"log"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/antihax/optional"
	"github.com/spf13/viper"
	qClient "go.qase.io/client"

	"github.com/insolar/testrail-cli/types"
)

var (
	autotestUserID = 10
	statusMap      = map[string]string{
		types.TestStatusPassed: "1",
		types.TestStatusFailed: "2",
		// types.TestStatusSkipped:      6,
		// types.TestStatusNotAvailable: 7,
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
	c *qClient.APIClient

	projectID string
	suiteID   int32
	tests     map[int]qClient.ResultCreate
	// defaultTests types.TestCasesWithDescription // todo?
}

func NewUploader(token string) *Uploader {
	configuration := qClient.NewConfiguration()
	configuration.AddDefaultHeader("Accept", "application/json")
	configuration.AddDefaultHeader("Content-Type", "application/json")
	configuration.AddDefaultHeader("Token", token)
	return &Uploader{
		c:     qClient.NewAPIClient(configuration),
		tests: make(map[int]qClient.ResultCreate),
	}
}

func (m Uploader) FormatURL(id int) string {
	return path.Join(viper.GetString("URL"), "/index.php?/cases/view/", strconv.Itoa(id)) // todo не так важно
}

func (m *Uploader) GetCasesWithDescription() types.TestCasesWithDescription {
	cases, _, err := m.c.CasesApi.GetCases(context.Background(),
		m.projectID,
		&qClient.CasesApiGetCasesOpts{FiltersSuiteId: optional.NewInt32(m.suiteID)})
	if err != nil {
		log.Fatal(err)
	}

	// todo дописать проверок
	var casesWithDescription types.TestCasesWithDescription
	for _, c := range cases.Result.Entities {
		caseWithDescription := types.TestCaseWithDescription{
			ID:          int(c.Id),
			Description: c.Title,
		}
		casesWithDescription = append(casesWithDescription, caseWithDescription)
	}
	return casesWithDescription
}

func (m *Uploader) Init(projectID string, suiteID int32) {
	ctx := context.Background()
	m.projectID = projectID
	m.suiteID = suiteID
	casesForRun := make([]int64, 0)
	suiteCases, _, err := m.c.CasesApi.GetCases(ctx, m.projectID,
		&qClient.CasesApiGetCasesOpts{FiltersSuiteId: optional.NewInt32(m.suiteID)})
	if err != nil && suiteCases.Result == nil {
		log.Fatal(err)
	}
	for _, c := range suiteCases.Result.Entities {
		casesForRun = append(casesForRun, c.Id)
	}

	run, _, err := m.c.RunsApi.CreateRun(ctx,
		qClient.RunCreate{
			Title:      "Automated test run " + time.Now().String(), // todo check
			Cases:      casesForRun,
			IsAutotest: true,
		}, m.projectID)
	if err != nil && run.Result == nil {
		log.Fatal(err)
	}
	log.Printf("Created run id = %d", run.Result.Id)
}

func (m *Uploader) AddTests(objects []*types.TestMatcher, ignoreNonExistent bool) {
	for _, object := range objects {
		if _, ok := m.tests[object.ID]; !ok && ignoreNonExistent {
			continue
		}
		// todo посмотреть что постить
		m.tests[object.ID] = qClient.ResultCreate{
			CaseId:      int64(object.ID),
			Case_:       nil,
			Status:      "",
			Time:        0,
			TimeMs:      0,
			Defect:      false,
			Attachments: nil,
			Stacktrace:  "",
			Comment:     "Autotest " + object.GoTestName,
		}
	}
}

func (m *Uploader) Upload() {
	ctx := context.Background()
	bulk := qClient.ResultCreateBulk{}

	for _, resultForCase := range m.tests {
		bulk.Results = append(bulk.Results, resultForCase)
	}

	_, _, err := m.c.ResultsApi.CreateResultBulk(ctx, bulk, m.projectID, m.suiteID)
	if err != nil {
		log.Fatal(err)
	}
}
