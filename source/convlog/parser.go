//  Copyright 2020 Insolar Network Ltd.
//  All rights reserved.
//  This material is licensed under the Insolar License version 1.0,
//  available at https://github.com/insolar/testrail-cli/LICENSE.md.

package convlog

import (
	"bytes"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/insolar/testrail-cli/source"
)

var _ source.Parser = (*Parser)(nil)

type Parser struct{}

var (
	// PASS|FAIL|ok  | all have 4 character len
	// and 4 spaces is common intendation
	commonPrefixLen = 4
	// printed by test on successful run.
	bigPass = []byte("PASS\n")

	// printed by test after a normal test failure.
	bigFail = []byte("FAIL\n")

	okBytes = []byte("ok  \t")

	// printed by 'go test' along with an error if the test binary terminates
	// with an error.
	bigFailErrorPrefix = []byte("FAIL\t")

	updates = [][]byte{
		[]byte("=== RUN   "),
		[]byte("=== PAUSE "),
		[]byte("=== CONT  "),
	}

	reports = [][]byte{
		[]byte("--- PASS: "),
		[]byte("--- FAIL: "),
		[]byte("--- SKIP: "),
		[]byte("--- BENCH: "),
	}

	fourSpace = []byte("    ")

	skipLinePrefix = []byte("?   \t")
	skipLineSuffix = []byte("\t[no test files]\n")

	convLogPrefixRe = regexp.MustCompile(`(?m)^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+\+\d{2}:\d{2}\s.{3}\s`)
)

func (Parser) Parse(input io.Reader) []source.TestEvent {
	scanner := NewLineScanner(input)

	converter := pkgConverter{}

	var res []source.TestEvent
	chunk := make([]source.TestEvent, 0)

	for scanner.Scan() {
		chunk = append(chunk, converter.handleInputLine(scanner.Text())...)

		if converter.finished {
			for i := range chunk {
				if chunk[i].Package == "" {
					chunk[i].Package = converter.pkg
				}
			}
			chunk = append(chunk, converter.flushReport(0)...)
			if len(chunk) == 0 {
				continue
			}
			chunk = append(chunk, source.TestEvent{
				Action:  converter.result,
				Package: converter.pkg,
				Elapsed: converter.elapsed,
			})
			res = append(res, chunk...)
			chunk = make([]source.TestEvent, 0)
			converter = pkgConverter{}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	if converter.result != "" {
		res = append(res, converter.flushReport(0)...)
		res = append(res, source.TestEvent{
			Action:  converter.result,
			Package: converter.pkg,
			Elapsed: converter.elapsed,
		})
	}

	return res
}

type pkgConverter struct {
	pkg      string             // package to name in events
	elapsed  float64            // duration in seconds
	testName string             // name of current test, for output attribution
	report   []source.TestEvent // pending test result reports (nested for subtests)
	result   string             // overall test result if seen
	finished bool
}

func (c *pkgConverter) flushReport(depth int) []source.TestEvent {
	res := make([]source.TestEvent, 0, len(c.report)-depth)
	for len(c.report) > depth {
		e := c.report[len(c.report)-1]
		c.report = c.report[:len(c.report)-1]
		if e.Test == "" {
			e.Test = c.testName
		}
		e.Package = c.pkg
		res = append(res, e)
	}
	c.testName = ""
	return res
}

func (c *pkgConverter) handleInputLine(line []byte) []source.TestEvent {
	// Final PASS or FAIL.
	if bytes.Equal(line, bigPass) {
		return []source.TestEvent{
			{
				Action: "output",
				Output: string(line),
			},
		}
	}
	if bytes.Equal(line, bigFail) {
		if c.pkg == "" {
			return []source.TestEvent{}
		}
		return []source.TestEvent{
			{
				Action: "output",
				Output: string(line),
			},
		}
	}

	if bytes.HasPrefix(line, bigFailErrorPrefix) || bytes.HasPrefix(line, okBytes) {
		c.result = "pass"
		if bytes.Equal(line[:commonPrefixLen], []byte("FAIL")) {
			c.result = "fail"
		}
		info := line[5:]
		tabSepIndex := bytes.Index(info, []byte{0x09})
		c.pkg = string(info[:tabSepIndex])
		c.finished = true

		if !bytes.Contains(info, []byte("(cached)")) {
			elapsedPart := info[tabSepIndex+1:]

			sIndex := bytes.Index(elapsedPart, []byte{'s'})

			elapsedString := string(elapsedPart[:sIndex+1])

			c.elapsed = parseSeconds(elapsedString)
		}

		res := c.flushReport(0)
		res = append(res, source.TestEvent{
			Action: "output",
			Output: string(line),
		})

		return res
	}

	res := make([]source.TestEvent, 0)

	// Special case for entirely skipped test binary: "?   \tpkgname\t[no test files]\n" is only line.
	// Report it as plain output but remember to say skip in the final summary.
	if bytes.HasPrefix(line, skipLinePrefix) && bytes.HasSuffix(line, skipLineSuffix) && len(c.report) == 0 {
		// 5 is ? + 3 spaces + 1 tab -- see example above
		info := line[5:]
		tabSepIndex := bytes.Index(info, []byte{0x09})
		res = c.flushReport(0)
		c.pkg = string(info[:tabSepIndex])
		c.result = "skip"
		c.finished = true
	}

	// "=== RUN   "
	// "=== PAUSE "
	// "=== CONT  "
	actionColon := false
	origLine := line
	ok := false
	indent := 0
	for _, magic := range updates {
		if bytes.HasPrefix(line, magic) {
			ok = true
			break
		}
	}
	if !ok {
		// "--- PASS: "
		// "--- FAIL: "
		// "--- SKIP: "
		// "--- BENCH: "
		// but possibly indented.
		for bytes.HasPrefix(line, fourSpace) {
			line = line[len(fourSpace):]
			indent++
		}
		for _, magic := range reports {
			if bytes.HasPrefix(line, magic) {
				actionColon = true
				ok = true
				break
			}
		}
	}

	// Not a special test output line.
	if !ok {
		// Lookup the name of the test which produced the output using the
		// indentation of the output as an index into the stack of the current
		// subtests.
		// If the indentation is greater than the number of current subtests
		// then the output must have included extra indentation. We can't
		// determine which subtest produced this output, so we default to the
		// old behaviour of assuming the most recently run subtest produced it.
		if indent > 0 && indent <= len(c.report) {
			c.testName = c.report[indent-1].Test
		}

		origLine = []byte(convLogPrefixRe.ReplaceAllString(string(origLine), ""))

		name := c.testName

		res = append(res, source.TestEvent{
			Action: "output",
			Output: string(origLine),
			Test:   name,
		})
		return res
	}

	// Parse out action and test name.
	i := 0
	if actionColon {
		i = bytes.IndexByte(line, ':') + 1
	}
	if i == 0 {
		i = len(updates[0])
	}
	action := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(string(line[commonPrefixLen:i])), ":"))
	name := strings.TrimSpace(string(line[i:]))

	e := source.TestEvent{Action: action}
	if line[0] == '-' { // PASS or FAIL report
		// Parse out elapsed time.
		if i := strings.Index(name, " ("); i >= 0 {
			if strings.HasSuffix(name, "s)") {
				t, err := strconv.ParseFloat(name[i+2:len(name)-2], 64)
				if err == nil {
					e.Elapsed = t
				}
			}
			name = name[:i]
		}

		if len(c.report) < indent {
			// Nested deeper than expected.
			// Treat this line as plain output.
			return append(res, source.TestEvent{
				Action: "output",
				Output: string(origLine),
				Test:   name,
			})
		}
		// Flush reports at this indentation level or deeper.
		res = append(res, c.flushReport(indent)...)
		e.Test = name
		c.testName = name
		c.report = append(c.report, e)
		res = append(res, source.TestEvent{
			Action: "output",
			Output: string(origLine),
			Test:   c.testName,
		})
		return res
	}
	// === update.
	// Finish any pending PASS/FAIL reports.
	res = append(res, c.flushReport(0)...)
	c.testName = name
	e.Test = name

	if action == "pause" {
		// For a pause, we want to write the pause notification before
		// delivering the pause event, just so it doesn't look like the test
		// is generating output immediately after being paused.
		res = append(res, source.TestEvent{
			Action: "output",
			Output: string(line),
			Test:   c.testName,
		})
	}
	res = append(res, e)
	if action != "pause" {
		res = append(res, source.TestEvent{
			Action: "output",
			Output: string(line),
			Test:   c.testName,
		})
	}

	return res
}

func parseSeconds(line string) float64 {
	res, _ := time.ParseDuration(line)
	return res.Seconds()
}
