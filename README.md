# workerpool

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
