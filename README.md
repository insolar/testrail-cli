#### go test json integration to TestRail

```
go get github.com/insolar/testrail-cli
go install cmd/testrail-cli/testrail-cli.go
```
| Param key     | Value             |
| ------------- | ----------------- |
| -u            | testrail user     |
| -p            | testrail password |
| -url          | testrail url      |
| -r            | testrail run id   |

```
go test ./... -json | testrail-cli -u ${testrail_user} -p ${testrail_token} -r 57
```