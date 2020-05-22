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
		t.Log("C3605 Some testcase description"")
        ...
	})
}
```

If you want to skip test, you can add issue to skip description
```go
func TestExample3(t *testing.T) {
	t.Log("C3607 Some testcase description")
	t.Skip("https://example.net/browse/TASK-1 other description")
}
```

#### Run
| Param key     |    Env key    | Description                 |
| ------------- | ------------- | --------------------------- |
| --URL         |   TR_URL      | testrail url                |
| --USER        |   TR_USER     | testrail user               |
| --PASSWORD    |   TR_PASSWORD | testrail password           |
| --RUN_ID      |   TR_RUN_ID   | testrail run id             |
| --FILE        |   TR_FILE     | go test json file           |
| --SKIP-DESC   |   SKIP-DESC   | skip description check flag |

Use params
```
testrail-cli --URL=https://example.testrail.com/ --USER=example@gmail.com --PASSWORD=${pass} --RUN_ID=57 --FILE=example_test.log
```
Or env vars with TR prefix
```
TR_URL=https://example.testrail.com/ TR_USER=example@gmail.com TR_PASSWORD=${pass} TR_RUN_ID=57 TR_FILE=example_test_suite.json testrail-cli
```
Also you can pipe json in
```
go test ./... -json | testrail-cli --URL=https://example.testrail.com/ --USER=example@gmail.com --PASSWORD=${pass} --RUN_ID=57
```
Or save file using tee for debug
```
go test ./... -json | tee autotest.log | testrail-cli --URL=https://example.testrail.com/ --USER=example@gmail.com --PASSWORD=${pass} --RUN_ID=57
```
