# workerpool

[![codecov](https://codecov.io/gh/yarcode/workerpool/branch/main/graph/badge.svg)](https://codecov.io/gh/yarcode/workerpool)
[![golangci](https://golangci.com/badges/github.com/yarcode/workerpool.svg)](https://golangci.com/r/github.com/yarcode/workerpool)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/yarcode/workerpool)
[![Go Report Card](https://goreportcard.com/badge/github.com/yarcode/workerpool)](https://goreportcard.com/report/github.com/yarcode/workerpool)

Package workerpool provides a service for running small parts of code (called jobs) in a background.

Jobs could have contexts, timeouts, rich retry strategies.

## Examples

```golang

pool := New()
pool.Start()

job := func(ctx context.Context) error {
    fmt.Println("hello")
    return nil
}

pool.Run(job)

```

### AdvancedUsage

```golang

pool := New()
pool.Start()

job := func(ctx context.Context) error {
    //
    // some tricky logic goes here
    //

    return nil
}
// add 3 seconds timeout for a job execution
job = AddTimeout(job, time.Second*3)
// retry job execution withing 5 attempts
job = AddRetry(job, strategy.Limit(5))

pool.Run(job)

```
