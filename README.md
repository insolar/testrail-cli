#### go test json integration to TestRail

```
go get github.com/insolar/testrail-cli
```

```
go test ./... -json | testrail-cli -u ${testrail_user} -p ${testraul_token} -run_id 57
```