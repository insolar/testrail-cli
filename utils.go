package testrail_cli

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"strings"
)

func TicketFromURL(url string) string {
	if strings.HasPrefix(url, "https") || strings.HasPrefix(url, "http") {
		s := strings.Split(url, "/")
		return s[len(s)-1]
	}
	return url
}

func Read(stream io.Reader) []*TestEvent {
	testEvents := make([]*TestEvent, 0)
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "\n" {
			break
		}
		var te *TestEvent
		if err := json.Unmarshal([]byte(line), &te); err != nil {
			log.Fatalf("failed to unmarshal test event json: %s\n", err)
		}
		testEvents = append(testEvents, te)
	}
	return testEvents
}
