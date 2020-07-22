// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"log"

	"github.com/insolar/testrail-cli/types"
)

type URLFormatter interface {
	FormatURL(int) string
}

type TestObjectSummary struct {
	Valid          []*types.TestMatcher
	NotFound       []*types.TestMatcher
	WrongDesc      []*types.TestMatcher
	SkippedNoIssue []*types.TestMatcher
}

func (s TestObjectSummary) LogInvalidTests(f URLFormatter) {
	if len(s.NotFound) > 0 {
		log.Println("Tests without testrail case ID:")
		for _, o := range s.NotFound {
			log.Printf("  %s", o.GoTestName)
		}
	}

	if len(s.WrongDesc) > 0 {
		log.Println("Test title discrepancy with testrail test-case title:")
		for _, o := range s.WrongDesc {
			log.Printf("  %s", o.GoTestName)
			log.Printf("    Test Description:     %s", o.Description)
			log.Printf("    Original Description: %s", o.OriginalDescription)
			log.Printf("    TestCase URL:         %s", f.FormatURL(o.ID))
		}
	}

	if len(s.SkippedNoIssue) > 0 {
		log.Println("Skipped tests without issue:")
		for _, o := range s.SkippedNoIssue {
			log.Printf("  %s", o.GoTestName)
		}
	}
}
