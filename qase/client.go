//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package qase

import (
	"context"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/antihax/optional"
	"github.com/spf13/viper"
	qClient "go.qase.io/client"

	"github.com/insolar/testrail-cli/types"
)

var (
	// autotestUserID = 10
	statusMap = map[string]string{
		types.TestStatusPassed:  "passed",
		types.TestStatusFailed:  "failed",
		types.TestStatusSkipped: "skipped",
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

	runId     int32
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
	cases, http, err := m.c.CasesApi.GetCases(context.Background(),
		m.projectID,
		&qClient.CasesApiGetCasesOpts{FiltersSuiteId: optional.NewInt32(m.suiteID)})
	checkResponse(http, err)

	if cases.Result == nil {
		log.Fatal("Cases response empty")
	}
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

func (m *Uploader) Init(projectID string, suiteID int32, title string) {
	ctx := context.Background()
	m.projectID = projectID
	m.suiteID = suiteID

	// prepare cases for run
	suiteCases, http, err := m.c.CasesApi.GetCases(ctx, m.projectID,
		&qClient.CasesApiGetCasesOpts{FiltersSuiteId: optional.NewInt32(m.suiteID)})
	checkResponse(http, err)

	casesForRun := make([]int64, 0)
	if suiteCases.Result != nil {
		for _, c := range suiteCases.Result.Entities {
			casesForRun = append(casesForRun, c.Id)
		}
	}
	if len(casesForRun) == 0 {
		log.Fatal("Cases is empty")
	}

	// create run
	run, _, err := m.c.RunsApi.CreateRun(ctx,
		qClient.RunCreate{
			Title:      title,
			Cases:      casesForRun,
			IsAutotest: true,
		}, m.projectID)
	checkResponse(http, err)

	if run.Result != nil {
		m.runId = int32(run.Result.Id)
		log.Printf("Created run id = %d", run.Result.Id)
	}
}

func (m *Uploader) AddTests(objects []*types.TestMatcher) {
	for _, object := range objects {
		// todo посмотреть что постить
		m.tests[object.ID] = qClient.ResultCreate{
			CaseId: int64(object.ID),
			Status: object.Status,
			// Time:        0,
			// TimeMs:      0,
			// Defect:      false,
			// Attachments: nil,
			// Stacktrace:  "",
			Comment: "Autotest " + object.GoTestName,
		}
	}
}

func (m *Uploader) Upload() {
	ctx := context.Background()
	results := make([]qClient.ResultCreate, 0)

	for _, resultForCase := range m.tests {
		results = append(results, resultForCase)
	}

	_, http, err := m.c.ResultsApi.CreateResultBulk(ctx, qClient.ResultCreateBulk{Results: results}, m.projectID, m.runId)
	checkResponse(http, err)

	_, _, err = m.c.RunsApi.CompleteRun(ctx, m.projectID, m.runId)
	checkResponse(http, err)
}

func checkResponse(http *http.Response, err error) {
	if http.StatusCode > 200 {
		log.Fatal("Unexpected status code")
		log.Printf("Actual status code: %d", http.StatusCode)
	}
	if err != nil {
		log.Fatal(err)
	}
}
