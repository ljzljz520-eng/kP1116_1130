# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
?   	showroom/cmd/showroom	[no test files]
ok  	showroom/internal/analytics	0.019s
?   	showroom/internal/audit	[no test files]
?   	showroom/internal/catalog	[no test files]
ok  	showroom/internal/config	0.026s
?   	showroom/internal/diagnostics	[no test files]
ok  	showroom/internal/display	0.165s
ok  	showroom/internal/gesture	0.008s
ok  	showroom/internal/httpapi	0.170s
ok  	showroom/internal/model	0.004s
ok  	showroom/internal/particles	0.007s
ok  	showroom/internal/persistence	0.212s
ok  	showroom/internal/rehearsal	0.010s
ok  	showroom/internal/scheduler	0.008s
?   	showroom/internal/session	[no test files]
--- FAIL: TestWelcomeUsesLatestPhrase (0.15s)
    service_test.go:41: ending phrase = "First Confirmed", want latest confirmation
FAIL
FAIL	showroom/internal/welcome	0.264s
ok  	showroom/internal/workflow	0.493s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/showroom): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/showroom): exit `0`
