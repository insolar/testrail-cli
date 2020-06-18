//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package source_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/testrail-cli/source"
	"github.com/insolar/testrail-cli/source/json"
	"github.com/insolar/testrail-cli/source/text"
)

func TestParsers(t *testing.T) {
	jsonFile, err := os.Open("json/example_test.log")
	require.NoError(t, err)

	jsonRes := json.Parser{}.Parse(jsonFile)

	textFile, err := os.Open("text/example_test.log")
	require.NoError(t, err)

	textRes := text.Parser{}.Parse(textFile)

	var JSONDiff []source.TestEvent
	for _, e := range jsonRes {
		found := false
		for _, te := range textRes {
			e.Elapsed = 0
			te.Elapsed = 0
			found = te == e
			if found {
				break
			}
		}

		// json captures one extra line
		if e.Test == "" && e.Output == "FAIL\n" && e.Package == "github.com/insolar/testrail-cli/package1" {
			continue
		}

		if !found {
			JSONDiff = append(JSONDiff, e)
		}
	}

	source.DumpEvents(JSONDiff)

	assert.Nil(t, JSONDiff)
}
