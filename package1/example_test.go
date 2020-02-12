package testrail_cli

import (
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
	t.Log("C3605 Pass testsdf")
}

func TestExample2(t *testing.T) {
	t.Parallel()
	time.Sleep(300 * time.Millisecond)
	t.Log("C3606 Fail test")
	t.Fail()
}

func TestExample3(t *testing.T) {
	time.Sleep(500 * time.Millisecond)
	t.Log("C3607 Skip test")
	t.Skip("issue: https://insolar.atlassian.net/browse/OPS-8 some bad desc")
}
