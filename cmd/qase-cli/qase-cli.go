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

	"github.com/insolar/testrail-cli/cmd/qase-cli/internal"
	"github.com/insolar/testrail-cli/parser"
	"github.com/insolar/testrail-cli/parser/json"
	"github.com/insolar/testrail-cli/parser/text"
	"github.com/insolar/testrail-cli/qase"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("QS")
	flag.String("API_TOKEN", "", "Qase apiToken")
	flag.String("FILE", "", "go test json file")
	flag.String("PROJECT_ID", "", "qase project id")
	flag.Int("SUITE_ID", 0, "qase suite id")
	flag.Bool("SKIP-DESC", false, "skip description check")
	flag.String("FORMAT", "json", "test output format")
	flag.String("MATCHER", "default", "test output matcher")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	var (
		apiToken  = viper.GetString("API_TOKEN")
		projectID = viper.GetString("PROJECT_ID")
		suiteID   = viper.GetInt32("SUITE_ID")
		file      = viper.GetString("FILE")
		skipDesc  = viper.GetBool("SKIP-DESC")
	)

	if apiToken == "" {
		log.Fatal("provide qase apiToken")
	}
	if projectID == "" {
		log.Fatal("provide qase project id")
	}
	if suiteID == 0 {
		log.Fatal("provide suite id, ex.: --SUITE_ID=54, or env QS_SUITE_ID=54")
	}

	var (
		parserName     = viper.GetString("format")
		parserInstance parser.Parser
	)
	switch parserName {
	case "json":
		parserInstance = json.Parser{}
	case "text":
		parserInstance = text.Parser{}
	default:
		log.Fatalf("Unsupported format %s", parserName)
	}

	matcherInstance := internal.Converter{}

	var stream io.Reader = os.Stdin
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		stream = f
	}
	eventReader := parserInstance.GetParseIterator(stream)
	tObjects := matcherInstance.ConvertEventsToMatcherObjects(eventReader)

	t := qase.NewUploader(apiToken)
	t.Init(projectID, suiteID)

	filteredObjects := internal.FilterTestObjects(tObjects, t.GetCasesWithDescription(), skipDesc)
	filteredObjects.LogInvalidTests(t)

	t.AddTests(filteredObjects.Valid, true)
	t.Upload()
}
