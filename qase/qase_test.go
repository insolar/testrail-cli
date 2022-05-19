package qase

import (
	"testing"

	"github.com/insolar/testrail-cli/types"
)

func TestUploaderQase(t *testing.T) {
	uploader := NewUploader("3e4a3a04b911d277c4b70e7ac0a4b22b4b87d9b1")
	uploader.Init("SOVEREN", 4, "Backend::Stats__141__fc9b32f")
	matcher := types.TestMatcher{
		ID:          1,
		Status:      "passed",
		Description: "blabla",
		GoTestName:  "TestUploaderQase",
	}
	matcher2 := types.TestMatcher{
		ID:         2,
		Status:     "failed",
		GoTestName: "TestUploaderQase",
		IssueURL:   "https://soveren.atlassian.net/browse/SMAT-98",
	}
	matcher3 := types.TestMatcher{
		ID:         3,
		Status:     "skipped",
		GoTestName: "TestUploaderQase",
		IssueURL:   "https://soveren.atlassian.net/browse/SMAT-98",
	}
	uploader.AddTests([]*types.TestMatcher{&matcher, &matcher2, &matcher3})
	uploader.Upload()
}
