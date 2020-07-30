// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package main

import (
	"io"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/insolar/testrail-cli/cmd/testrail-assuredledger-cli/internal"
	"github.com/insolar/testrail-cli/parser/convlog"
	"github.com/insolar/testrail-cli/testrail"
	"github.com/insolar/testrail-cli/types"
)

func main() {
	pflag.String("url", "", "testrail url")
	pflag.String("user", "", "testrail username")
	pflag.String("password", "", "testrail password/token")
	pflag.String("file", "", "go test json file")
	pflag.Int("run-id", 0, "testrail run id")
	pflag.Parse()

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TR")
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	var (
		url      = viper.GetString("url")
		user     = viper.GetString("user")
		pass     = viper.GetString("password")
		runID    = viper.GetInt("run-id")
		file     = viper.GetString("file")
	)

	if url == "" {
		log.Fatal("provide TestRail url")
	}
	if runID == 0 {
		log.Fatal("provide run id, ex.: --run-id 54, or env TR_RUN_ID=54")
	}
	if user == "" {
		log.Fatal("provide user for TestRail authentication")
	}
	if pass == "" {
		log.Fatal("provide password/token for TestRail authentication")
	}

	var stream io.Reader = os.Stdin
	if file != "" {
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		stream = f
	}
	parserInstance := convlog.Parser{}
	eventReader := parserInstance.GetParseIterator(stream)

	matcherInstance := internal.Converter{}
	tObjects := matcherInstance.ConvertEventsToMatcherObjects(eventReader)

	t := testrail.NewUploader(url, user, pass)
	t.Init(runID)

	for _, tObject := range tObjects {
		if !types.StatusKnown(tObject.Status) {
			tObject.Status = types.TestStatusFailed
		}
	}

	t.AddTests(tObjects, true)
	t.Upload()
}
