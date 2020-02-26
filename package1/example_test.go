// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/testrail-cli/LICENSE.md.

package testrail_cli

import (
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	t.Parallel()
	time.Sleep(100 * time.Millisecond)
	t.Log("C9999 Pass test")
}

func TestExample2(t *testing.T) {
	t.Parallel()
	time.Sleep(300 * time.Millisecond)
	t.Log("C3606 Fail testsdf")
	t.Fail()
}

func TestExample3(t *testing.T) {
	time.Sleep(500 * time.Millisecond)
	t.Log("C3607 Skip test")
	t.Skip("sdfs")
	//t.Skip("sdfs https://insolar.atlassian.net/browse/OPS-1 some bad desc")
}

func TestMGRGroupCreateWith2Members(t *testing.T) {
	t.Parallel()
	time.Sleep(1 * time.Second)
	t.Log("C3703 [MGR] Error creating group of 2 members")
}

func TestMGRGroupCreateCheckEmptySequence(t *testing.T) {
	t.Parallel()
	checkGroup(t, "200", []string{"100", "100", "100"}, "C3702 [MGR] Create group of 3 members with groupGoal=200 and check empty sequence")
	checkGroup(t, "300", []string{"100", "100", "100"}, "C3704 [MGR] Create group of 3 members with groupGoal=300 and check empty sequence")
	checkGroup(t, "400", []string{"100", "100", "100"}, "C3696 [MGR] Create group of 3 members with groupGoal=400 and check empty sequence")
}

func checkGroup(t *testing.T, groupGoal string, userGoals []string, testName string) {
	t.Run("groupGoal="+groupGoal, func(t *testing.T) {
		t.Log(testName)
	})
}

//
//func TestMGR(t *testing.T) {
//	t.Parallel()
//	checkGroup(t, "200", []string{"100", "100", "100"}, "C3694 [MGR] Create group of 3 members with groupGoal=200 and check empty sequence")
//	time.Sleep(290*time.Millisecond)
//	t.Log("C3702 [MGR] Create group of 3 members with groupGoal=200 and check empty sequence")
//}
