package workerpool

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Rican7/retry/strategy"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Example() {
	pool := New()
	pool.Start()

	job := func(ctx context.Context) error {
		fmt.Println("hello")
		return nil
	}

	pool.Run(job)
}

func Example_advancedUsage() {
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
}

func TestAddPanicRecovery(t *testing.T) {
	job := func(ctx context.Context) error {
		panic("oops")
	}
	job = AddPanicRecovery(job)
	err := job(context.Background())
	assert.EqualError(t, err, "panic recovered: oops")
}

func ExampleAddPanicRecovery() {
	pool := New()
	pool.Start()

	job := func(ctx context.Context) error {
		panic("oops")
	}
	pool.Run(job)
	// that will work fine (and get to logger)

	pool.Stop()
}

func ExampleAddRetry() {
	pool := New()
	pool.Start()

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
	pool.Start()

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
	pool.Start()

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
	pool.Start()

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
	// mock time for test purposes
	zerolog.TimestampFunc = func() time.Time {
		t, _ := time.Parse("2006-01-02", "2021-01-01")
		return t
	}

	l := zerolog.New(os.Stdout).Level(zerolog.WarnLevel).With().Timestamp().Logger()
	l1 := l.With().Str("logger", "logger1").Logger()
	l2 := l.With().Str("logger", "logger2").Logger()

	pool := New(WithLogger(l))
	pool.Start()

	job1 := func(ctx context.Context) error {
		zerolog.Ctx(ctx).Warn().Msg("hello from job1")
		return nil
	}
	job2 := func(ctx context.Context) error {
		zerolog.Ctx(ctx).Warn().Msg("hello from job2")
		return nil
	}
	// adding custom jobs
	job1 = AddLogger(job1, l1)
	job2 = AddLogger(job2, l2)

	pool.Run(job1)
	pool.Run(job2)

	pool.Stop()

	// Output: {"level":"warn","logger":"logger1","time":"2021-01-01T00:00:00Z","message":"hello from job1"}
	// {"level":"warn","logger":"logger2","time":"2021-01-01T00:00:00Z","message":"hello from job2"}
}
