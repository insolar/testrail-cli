// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package main

import (
	"flag"
	tr "github.com/insolar/testrail-cli"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

func main() {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TR")
	flag.String("URL", "https://insolar.testrail.io/", "testrail url")
	flag.String("USER", "autotest@insolar.io", "testrail username")
	flag.String("PASSWORD", "", "testrail password/token")
	flag.String("FILE", "", "go test json file")
	flag.Int("RUN_ID", 0, "testrail run id")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)
	viper.SetDefault("USER", "autotest@insolar.io")
	viper.SetDefault("URL", "https://insolar.testrail.io/")

	url := viper.GetString("URL")
	user := viper.GetString("USER")
	pass := viper.GetString("PASSWORD")
	runID := viper.GetInt("RUN_ID")
	file := viper.GetString("FILE")

	if runID == 0 {
		log.Fatal("provide run id, ex.: --RUN_ID=54, or env TR_RUN_ID=54")
	}
	if pass == "" {
		log.Fatal("provide password/token for TestRail authentication")
	}

	var stream io.Reader
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		//defer f.Close()
		stream = f
	} else {
		stream = os.Stdin
	}
	events := tr.ReadFile(stream)

	t := tr.NewTestRail(url, user, pass)
	run := t.GetRun(runID)
	casesWithDescs := t.GetCasesWithDescs(run.ProjectID, run.SuiteID)
	// update all cases with N/A status, we store all autotests in ONE run, so in case
	// someone delete particular case implementation status must be updated to N/A
	untested := t.NAResults(casesWithDescs)
	t.UpdateRunForCases(runID, untested)

	testEventsBatch := t.GroupEventsByTest(events)
	tObjects := t.EventsToTestObjects(testEventsBatch)
	filteredObjects := tr.FilterTestObjects(tObjects, casesWithDescs)
	tr.LogInvalidTests(filteredObjects)

	sendableResults := t.TestObjectsToSendableResultsForCase(filteredObjects.Valid)
	if sendableResults != nil {
		t.UpdateRunForCases(runID, sendableResults)
	}
}
