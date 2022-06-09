# workerpool

[![codecov](https://codecov.io/gh/.//branch/master/graph/badge.svg)](https://codecov.io/gh/./)
[![golangci](https://golangci.com/badges/./.svg)](https://golangci.com/r/./)
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/./)
[![Go Report Card](https://goreportcard.com/badge/./)](https://goreportcard.com/report/./)

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
pool.Stop()
```

 Output:

```
hello
```
