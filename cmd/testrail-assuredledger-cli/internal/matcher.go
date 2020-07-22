//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package internal

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/insolar/testrail-cli/parser"
	"github.com/insolar/testrail-cli/types"
)

var (
	testCaseIDRe    = regexp.MustCompile(`C(\d{1,8})`)
	testSkipIssueRe = regexp.MustCompile(`insolar\.atlassian\.net/browse/([A-Z]+-\d+)`)

)

type logLineParserState int
const (
	StateMessage logLineParserState = iota
	StateKey
	StateValue
)

func logLineParse(line string) (string, map[string]string, error) {
	var (
		state = StateMessage

		originalLine  = line
		messageFields = make(map[string]string)

		message   []byte
		lastKey   []byte
		lastValue []byte
		err       error

		skipEqualSign int
	)

	for {
		line = skipLeadingWhitespaces(line)
		if line == "" {
			return string(message), messageFields, nil
		}

		switch state {
		case StateMessage:
			line, message, err = parseMessage(line, skipEqualSign)
			if err != nil {
				return "", nil, err
			}

			state = StateKey
		case StateKey:
			line, lastKey, err = parseKeyString(line)
			if err == MalformedKeyHaveSpace {
				state = StateMessage
				skipEqualSign = len(messageFields) + 1
				messageFields = make(map[string]string)
				line = originalLine
				continue
			} else if err != nil {
				return "", nil, err
			}

			state = StateValue
		case StateValue:
			line, lastValue, err = parseJsonValueString(line)
			if err != nil {
				return "", nil, err
			}

			messageFields[string(lastKey)] = string(lastValue)

			state = StateKey
		default:
			panic("IllegalState")
		}
	}
}

func Reverse(line string) string {
	var (
		pos           = 0
		lineRuneCount = utf8.RuneCountInString(line)
		invertedLine  = make([]rune, lineRuneCount, lineRuneCount)
	)

	for _, r := range line {
		invertedLine[cap(invertedLine)-1 - pos] = r
		pos += 1
	}

	return string(invertedLine)
}

func ReverseBytesToString(val []byte) string {
	return Reverse(string(val))
}

func logLineParseAlternative(line string) (string, map[string]string, error) {
	var (
		lastKey   []byte
		lastValue []byte
		err       error

		messageFields = make(map[string]string)

		state = StateValue
	)

	line = Reverse(line)

	for {
		line = skipLeadingWhitespaces(line)
		if line == "" {
			return ReverseBytesToString(lastValue), messageFields, nil
		}

		switch state {
		case StateValue:
			line, lastValue, err = parseReverseJsonValueString(line)
			if err != nil {
				state = StateMessage
				continue
			}

			state = StateKey
		case StateKey:
			line, lastKey, err = parseReverseKeyString(line)
			if err != nil {
				state = StateMessage
				continue
			}

			messageFields[ReverseBytesToString(lastKey)] = ReverseBytesToString(lastValue)
			lastKey, lastValue = []byte(nil), []byte(nil)

			state = StateValue

		case StateMessage:
			return Reverse(line), messageFields, nil
		}
	}
}

type Converter struct {}

func (c Converter) ConvertEventsToMatcherObjectsPreload(events map[string][]parser.TestEvent) []*types.TestMatcher {
	reader := parser.NewStreamingEventReaderFromMap(events)
	return c.ConvertEventsToMatcherObjects(reader)
}

func (Converter) ConvertEventsToMatcherObjects(reader parser.EventReader) []*types.TestMatcher {
	matchers := make(map[string]*types.TestMatcher)

	for {
		_, event, err := reader.Next()
		if parser.IsEOF(err) {
			break
		} else if err != nil {
			log.Fatal(err)
		}


		if event.Action == "output" {
			if !strings.Contains(event.Output, "testrail ") {
				continue
			}

			lineMessage, lineFields, err := logLineParseAlternative(event.Output)
			if err != nil {
				panic(err)
			}

			pkgName, okPkg := lineFields["TestPackage"]
			testName, okTest := lineFields["TestName"]

			if lineMessage != "testrail" || !okPkg || !okTest {
				continue
			}

			name := parser.UniqueTestKeyFromFields(pkgName, testName)

			t, ok := matchers[name]
			if !ok {
				t = &types.TestMatcher{}
				matchers[name] = t
			}

			t.GoTestName = testName

			id, ok := lineFields["ID"]
			if !ok {
				continue
			}
			submatches := testCaseIDRe.FindStringSubmatch(id)
			if len(submatches) != 2 {
				continue
			}
			t.ID, err = strconv.Atoi(submatches[1])
			if err != nil {
				continue
			}

			issueURL, ok := lineFields["SkippedLink"]
			if ok {
				submatches := testSkipIssueRe.FindStringSubmatch(issueURL)
				if len(submatches) == 2 {
					t.IssueURL = submatches[1]
				}
			}

			t.Status = lineFields["Status"]
		}
	}

	matcherList := make([]*types.TestMatcher, 0, len(matchers))
	for _, val := range matchers {
		if val.ID == 0 {
			continue
		} else if val.Status == "" {
			val.Status = "FAIL"
		}
		matcherList = append(matcherList, val)
	}

	return matcherList
}