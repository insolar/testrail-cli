package main

import (
	"flag"
	tr "github.com/insolar/testrail-cli"
	"io"
	"log"
	"os"
)

func main() {
	url := flag.String("url", "https://insolar.testrail.io/", "testrail url")
	user := flag.String("u", "autotest@insolar.io", "testrail username")
	pass := flag.String("p", "", "testrail password/token")
	file := flag.String("f", "", "go test json file")
	runId := flag.Int("r", 0, "testrail run id")
	flag.Parse()

	if *runId == 0 {
		log.Fatal("provide run id, ex.: -r 54")
	}

	var stream io.Reader
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		stream = f
	} else {
		stream = os.Stdin
	}
	events := tr.Read(stream)

	t := tr.NewTestRail(*url, *user, *pass)
	run := t.GetRun(*runId)
	casesWithDescs := t.GetCasesWithDescs(run.ProjectID, run.SuiteID)
	// update all cases with N/A status, we store all autotests in ONE run, so in case
	// someone delete particular case implementation status must be updated to N/A
	untested := t.NAResults(casesWithDescs)
	t.UpdateRunForCases(*runId, untested)

	testEventsBatch := t.GroupEventsByTest(events)
	tObjects := t.EventsToTestObjects(testEventsBatch)
	filteredObjects := tr.FilterValidTests(tObjects, casesWithDescs)

	sendableResults := t.TestObjectsToSendableResultsForCase(filteredObjects)
	t.UpdateRunForCases(*runId, sendableResults)
}
