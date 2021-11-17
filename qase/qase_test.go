package qase

import (
	"testing"

	"github.com/insolar/testrail-cli/types"
)

func TestUploaderQase(t *testing.T) {
	uploader := NewUploader("3e4a3a04b911d277c4b70e7ac0a4b22b4b87d9b1")
	uploader.Init("SOVEREN", 4, "141__fc9b32f")
	matcher := types.TestMatcher{
		ID:     1,
		Status: "passed",
		// Description:         "blabla",
		// GoTestName:          "TestUploaderQase",
	}
	matcher2 := types.TestMatcher{
		ID:     2,
		Status: "failed",
		// Description:         "blabla",
		// GoTestName:          "TestUploaderQase",
	}
	uploader.AddTests([]*types.TestMatcher{&matcher, &matcher2})
	uploader.Upload()
}
