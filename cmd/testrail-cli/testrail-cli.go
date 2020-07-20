// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package main

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	tr "github.com/insolar/testrail-cli"
	"github.com/insolar/testrail-cli/source"
	"github.com/insolar/testrail-cli/source/convlog"
	"github.com/insolar/testrail-cli/source/json"
	"github.com/insolar/testrail-cli/source/text"
)

func main() {

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TR")
	flag.String("URL", "", "testrail url")
	flag.String("USER", "", "testrail username")
	flag.String("PASSWORD", "", "testrail password/token")
	flag.String("FILE", "", "go test json file")
	flag.Int("RUN_ID", 0, "testrail run id")
	flag.Bool("SKIP-DESC", false, "skip description check")
	flag.String("FORMAT", "json", "test output format")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	url := viper.GetString("URL")
	user := viper.GetString("USER")
	pass := viper.GetString("PASSWORD")
	runID := viper.GetInt("RUN_ID")
	file := viper.GetString("FILE")
	skipDesc := viper.GetBool("SKIP-DESC")
	format := viper.GetString("format")

	if url == "" {
		log.Fatal("provide TestRail url")
	}
	if runID == 0 {
		log.Fatal("provide run id, ex.: --RUN_ID=54, or env TR_RUN_ID=54")
	}
	if user == "" {
		log.Fatal("provide user for TestRail authentication")
	}
	if pass == "" {
		log.Fatal("provide password/token for TestRail authentication")
	}
	if format == "" {
		log.Fatal("provide input format")
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

	var parser source.Parser
	switch format {
	case "json":
		parser = json.Parser{}
	case "text":
		parser = text.Parser{}
	case "convlog":
		parser = convlog.Parser{}
	default:
		log.Fatalf("Unsupported format %s", format)
	}

	events := parser.Parse(stream)

	t := tr.NewTestRail(url, user, pass)
	run := t.GetRun(runID)
	casesWithDescs := t.GetCasesWithDescs(run.ProjectID, run.SuiteID)
	// update all cases with N/A status, we store all autotests in ONE run, so in case
	// someone delete particular case implementation status must be updated to N/A
	untested := t.NAResults(casesWithDescs)
	t.UpdateRunForCases(runID, untested)

	testEventsBatch := t.GroupEventsByTest(events)
	tObjects := t.EventsToTestObjects(testEventsBatch)
	filteredObjects := tr.FilterTestObjects(tObjects, casesWithDescs, skipDesc)
	tr.LogInvalidTests(filteredObjects)

	sendableResults := t.TestObjectsToSendableResultsForCase(filteredObjects.Valid)
	if sendableResults != nil {
		t.UpdateRunForCases(runID, sendableResults)
	}
}
