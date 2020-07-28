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

	"github.com/insolar/testrail-cli/cmd/testrail-cli/internal"
	"github.com/insolar/testrail-cli/parser"
	"github.com/insolar/testrail-cli/parser/json"
	"github.com/insolar/testrail-cli/parser/text"
	"github.com/insolar/testrail-cli/testrail"
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
	flag.String("MATCHER", "default", "test output matcher")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	var (
		url      = viper.GetString("URL")
		user     = viper.GetString("USER")
		pass     = viper.GetString("PASSWORD")
		runID    = viper.GetInt("RUN_ID")
		file     = viper.GetString("FILE")
		skipDesc = viper.GetBool("SKIP-DESC")
	)

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

	var (
		parserName = viper.GetString("format")
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

	t := testrail.NewUploader(url, user, pass)
	t.Init(runID)

	filteredObjects := internal.FilterTestObjects(tObjects, t.GetCasesWithDescription(), skipDesc)
	filteredObjects.LogInvalidTests(t)

	t.AddTests(filteredObjects.Valid)
	t.Upload()
}
