//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"log"
	"regexp"
	"strconv"

	"github.com/insolar/testrail-cli/parser"
	"github.com/insolar/testrail-cli/types"
)

var (
	testStatusRe    = regexp.MustCompile(`--- (.*):`)
	testSkipIssueRe = regexp.MustCompile(`insolar\.atlassian\.net/browse/([A-Z]+-\d+)`)
	testCaseIdRe    = regexp.MustCompile(`C(\d{1,8})\s(.*)`)
)

type Converter struct {}

func (c Converter) ConvertEventsToMatcherObjectsPreload(events map[string][]parser.TestEvent) []*types.TestMatcher {
	reader := parser.NewStreamingEventReaderFromMap(events)
	return c.ConvertEventsToMatcherObjects(reader)
}

func (Converter) ConvertEventsToMatcherObjects(reader parser.EventReader) []*types.TestMatcher {
	matchers := make(map[string]*types.TestMatcher)

	for {
		name, event, err := reader.Next()
		if parser.IsEOF(err) {
			break
		} else if err != nil {
			log.Fatal(err)
		}


		if event.Action == "output" {
			if event.Test == "" {
				continue
			}

			t, ok := matchers[name]
			if !ok {
				t = &types.TestMatcher{}
				matchers[name] = t
			}

			t.GoTestName = event.Test

			if res := testCaseIdRe.FindStringSubmatch(event.Output); len(res) == 3 {
				d, err := strconv.Atoi(res[1])
				if err != nil {
					log.Fatal(err)
				}
				// TODO: Bad, should harden regex instead, debatable for now
				if t.ID != 0 {
					continue
				}
				t.ID = d
				t.Description = res[2]
			} else if res := testStatusRe.FindStringSubmatch(event.Output); len(res) == 2 {
				t.Status = res[1]
			} else if res := testSkipIssueRe.FindStringSubmatch(event.Output); len(res) == 2 {
				t.IssueURL = res[1]
			}
		}
	}


	matcherList := make([]*types.TestMatcher, 0, len(matchers))
	for _, val := range matchers {
		matcherList = append(matcherList, val)
	}

	return matcherList
}