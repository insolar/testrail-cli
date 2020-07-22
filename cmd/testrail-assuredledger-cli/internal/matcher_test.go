//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_logLineParse(t *testing.T) {
	line := `testrail caller=testutils/investigation/testrail.go:126 ID=C5005 TestName=TestConstructor_SamePulse_AfterExecution TestPackage=github.com/insolar/assured-ledger/ledger-core/virtual/integration/deduplication testname=TestConstructor_SamePulse_WhileExecution`

	expectedFields := map[string]string{
		"caller":"testutils/investigation/testrail.go:126",
		"ID":"C5005",
		"TestName": "TestConstructor_SamePulse_AfterExecution",
		"TestPackage": "github.com/insolar/assured-ledger/ledger-core/virtual/integration/deduplication",
		"testname": "TestConstructor_SamePulse_WhileExecution",
	}

	t.Run("test basic", func(t *testing.T) {
		msg, fields, err := logLineParse(line)

		assert.Equal(t, "testrail", msg)
		assert.Equal(t, expectedFields, fields)
		assert.NoError(t, err)
	})
}

func Test_logLineParseAlternative (t *testing.T) {
	line1 := `testrail caller=testutils/investigation/testrail.go:126 ID=C5005 TestName=TestConstructor_SamePulse_AfterExecution TestPackage=github.com/insolar/assured-ledger/ledger-core/virtual/integration/deduplication testname=TestConstructor_SamePulse_WhileExecution`
	expectedFields1 := map[string]string{
		"caller":"testutils/investigation/testrail.go:126",
		"ID":"C5005",
		"TestName": "TestConstructor_SamePulse_AfterExecution",
		"TestPackage": "github.com/insolar/assured-ledger/ledger-core/virtual/integration/deduplication",
		"testname": "TestConstructor_SamePulse_WhileExecution",
	}
	expectedMessage1 := "testrail"

	t.Run("test basic", func(t *testing.T) {
		msg, fields, err := logLineParseAlternative(line1)

		assert.Equal(t, expectedMessage1, msg)
		assert.Equal(t, expectedFields1, fields)
		assert.NoError(t, err)
	})

	line2 := `Got Bootstrap request from host id: 0 ref: insolar:1GZ2ZjnsgEsQp49Lz0inGkDKoY2RrJzF3XH_n7gAAAAY addr: 127.0.0.1:10006; RequestID = 2 caller=network/hostnetwork/hostnetwork.go:136 loginstance=node testname=TestNodeLeave traceid=`
	expectedFields2 := map[string]string{
		"caller":"network/hostnetwork/hostnetwork.go:136",
		"loginstance":"node",
		"testname":"TestNodeLeave",
		"traceid":"",
	}
	expectedMessage2 := "Got Bootstrap request from host id: 0 ref: insolar:1GZ2ZjnsgEsQp49Lz0inGkDKoY2RrJzF3XH_n7gAAAAY addr: 127.0.0.1:10006; RequestID = 2"

	t.Run("test with =", func(t *testing.T) {
		msg, fields, err := logLineParseAlternative(line2)

		assert.Equal(t, expectedMessage2, msg)
		assert.Equal(t, expectedFields2, fields)
		assert.NoError(t, err)
	})

	line3 := `=== AddJoinCandidate id = 2483507232, address = 127.0.0.1:10006  caller=network/gateway/base.go:349 loginstance=node testname=TestNodeLeave traceid=`
	expectedFields3 := map[string]string{
		"caller":"network/gateway/base.go:349",
		"loginstance":"node",
		"testname":"TestNodeLeave",
		"traceid":"",
	}
	expectedMessage3 := "=== AddJoinCandidate id = 2483507232, address = 127.0.0.1:10006"

	t.Run("another test with =", func(t *testing.T) {
		msg, fields, err := logLineParseAlternative(line3)

		assert.Equal(t, expectedMessage3, msg)
		assert.Equal(t, expectedFields3, fields)
		assert.NoError(t, err)
	})
}