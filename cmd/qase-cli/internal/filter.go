// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"log"

	"github.com/insolar/testrail-cli/types"
)

// FilterTestObjects split test objects into groups: valid/not found/wrong description
func FilterTestObjects(objectList []*types.TestMatcher, caseList types.TestCasesWithDescription, skipDesc bool) *TestObjectSummary {
	var (
		summary   TestObjectSummary
		objectMap = make(map[int]*types.TestMatcher)
		caseMap   = make(map[int]types.TestCaseWithDescription)
	)

	for _, c := range caseList {
		caseMap[c.ID] = c
	}

	for _, object := range objectList {
		if object.Status == "SKIP" && object.IssueURL == "" {
			summary.SkippedNoIssue = append(summary.SkippedNoIssue, object)
		}

		if object.ID != 0 {
			if _, ok := objectMap[object.ID]; ok {
				log.Printf("duplicate testcase: %d\n", object.ID)
			}
		}

		objectMap[object.ID] = object

		c, ok := caseMap[object.ID]
		if !ok {
			summary.NotFound = append(summary.NotFound, object)
			continue
		}

		object.OriginalDescription = c.Description
		if !skipDesc && object.Description != c.Description {
			summary.WrongDesc = append(summary.WrongDesc, object)
			continue
		}
		summary.Valid = append(summary.Valid, object)
	}

	return &summary
}
