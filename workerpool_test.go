package workerpool

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rican7/retry/strategy"
	"github.com/rs/zerolog"
	"os"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	pool := New()
	pool.Start(1)
	pool.Stop()
}

func Example() {
	pool := New()
	pool.Start(1)

	job := func(ctx context.Context) error {
		fmt.Println("hello")
		return nil
	}

	pool.Run(job)
	pool.Stop()
	// Output: hello
}

func Example_panicHandling() {
	pool := New()
	pool.Start(1)

	job := func(ctx context.Context) error {
		panic("oops")
	}
	pool.Run(job)
	// that will work fine (and get to logger)

	pool.Stop()
}

func ExampleAddRetry() {
	pool := New()
	pool.Start(1)

	retryCount := 0

	job := func(ctx context.Context) error {
		retryCount++
		// example job will always fail
		return errors.New("unrecoverable error")
	}
	job = AddRetry(job, strategy.Limit(5))

	pool.Run(job)
	time.Sleep(time.Millisecond * 100)
	pool.Stop()

	fmt.Println(retryCount)
	// Output: 5
}

func ExampleAddPostRun() {
	pool := New()
	pool.Start(1)

	job := func(ctx context.Context) error {
		// example job will always fail
		return errors.New("unrecoverable error")
	}
	job = AddPostRun(job, func(err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	pool.Run(job)
	time.Sleep(time.Millisecond * 100)
	pool.Stop()

	// Output: unrecoverable error
}

func ExampleAddTimeout() {
	pool := New()
	pool.Start(1)

	job := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(200 * time.Millisecond):
			fmt.Println("job ok")
			return nil
		}
	}
	// adding timeout to job
	job = AddTimeout(job, time.Millisecond*100)
	// adding post run hook to job to output timeout error
	job = AddPostRun(job, func(err error) {
		if err != nil {
			fmt.Println(err.Error())
		}
	})

	pool.Run(job)
	time.Sleep(time.Millisecond * 300)
	pool.Stop()

	// Output: context deadline exceeded
}

func ExampleAddContext() {
	pool := New()
	pool.Start(1)

	ctxKey := struct{}{}

	job := func(ctx context.Context) error {
		if v, ok := ctx.Value(ctxKey).(string); ok {
			fmt.Println(v)
		}
		return nil
	}
	// adding custom context to job
	ctx := context.WithValue(context.Background(), ctxKey, "from context")
	job = AddContext(job, ctx)

	pool.Run(job)
	pool.Stop()

	// Output: from context
}

func ExampleAddLogger() {
	pool := New()
	pool.Start(1)

	logger := zerolog.New(os.Stderr).Level(zerolog.WarnLevel).With().Timestamp().Logger()

	job1 := func(ctx context.Context) error {
		zerolog.Ctx(ctx).Warn().Msg("hello from job1")
		return nil
	}
	job2 := func(ctx context.Context) error {
		zerolog.Ctx(ctx).Warn().Msg("hello from job2")
		return nil
	}
	// adding custom logger only to job1
	job1 = AddLogger(job1, logger)

	pool.Run(job1)
	// Will output with custom logger:
	// {"level":"warn","time":"2022-03-10T13:19:42+07:00","message":"hello from job1"}
	pool.Run(job2)
	// Will output with default logger:
	// {"level":"warn","worker_id":0,"job_id":"c8kpgvnefj8pmh7j67lg","time":"2022-03-10T13:19:42+07:00","message":"hello from job2"}

	pool.Stop()
}
