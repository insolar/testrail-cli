#### Go test2json integration to TestRail

```
go get github.com/insolar/testrail-cli
go install cmd/testrail-cli/testrail-cli.go
```

#### Usage
Every test MUST log testrail case in format
```go
func TestExample(t *testing.T) {
	t.Log("C3605 Some testcase description")
	...
}
```

```go
func TestSuite(t *testing.T) {
	t.Run("TestExample", func(t *testing.T) {
		t.Log("C3605 Pass testsdf")
        ...
	})
}
```

If you want to skip test, you can add issue to skip description
```go
func TestExample3(t *testing.T) {
	t.Log("C3607 Some testcase description")
	t.Skip("issue: https://insolar.atlassian.net/browse/OPS-8 other description")
}
```

#### Run
| Param key     | Value             |
| ------------- | ----------------- |
| --URL         | testrail url      |
| --USER        | testrail user     |
| --PASSWORD    | testrail password |
| --RUN_ID      | testrail run id   |
| --FILE        | go test json file |

Use params
```
testrail-cli --URL=https://insolar.testrail.io/ --USER=autotest@insolar.io --PASSWORD=${pass} --RUN_ID=57 --FILE=example_test.json
```
Or env vars with TR prefix
```
TR_URL=https://insolar.testrail.io/ TR_USER=autotest@insolar.io TR_PASSWORD=${pass} TR_RUN_ID=57 TR_FILE=example_test_suite.json testrail-cli
```
Also you can pipe json in
```
go test ./... -json | testrail-cli --URL=https://insolar.testrail.io/ --USER=autotest@insolar.io --PASSWORD=${pass} --RUN_ID=57
```
