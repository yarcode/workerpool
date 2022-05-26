package workerpool

import (
	"context"
	"fmt"
	"github.com/Rican7/retry"
	"github.com/Rican7/retry/strategy"
	"github.com/rs/zerolog"
	"time"
)

// AddPanicRecovery to your job. If panic occurs, it converts panic into error.
func AddPanicRecovery(job Job) Job {
	return func(ctx context.Context) (retErr error) {
		defer func() {
			p := recover()
			if p != nil {
				retErr = fmt.Errorf("panic recovered: %v", p)
			}
		}()
		return job(ctx)
	}
}

// AddLogger replaces job context zerolog logger with the specified one. Could be useful for http handlers.
func AddLogger(job Job, l zerolog.Logger) Job {
	return func(ctx context.Context) error {
		ctx = l.WithContext(ctx)
		return job(ctx)
	}
}

// AddRetry middleware allows you to apply retry strategies to your job.
//
// See https://github.com/Rican7/retry for details.
func AddRetry(job Job, strategies ...strategy.Strategy) Job {
	return func(ctx context.Context) error {
		return retry.Retry(
			func(attempt uint) error {
				l := zerolog.Ctx(ctx).With().Uint("attempt", attempt).Logger()
				return AddLogger(job, l)(ctx)
			},
			strategies...,
		)
	}
}

// AddTimeout middleware allows you to add timeout to your job.
// It will be set on the job context.
func AddTimeout(job Job, timeout time.Duration) Job {
	return func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return job(ctx)
	}
}

// AddPostRun middleware allows you to add some logic once was completed or failed.
// The err parameter in hook will contain error received from job.
func AddPostRun(job Job, hook func(err error)) Job {
	return func(ctx context.Context) error {
		err := job(ctx)
		hook(err)
		return err
	}
}

// AddContext replaces the default context (context.Background()) with the specified one only for this job.
func AddContext(job Job, c context.Context) Job {
	return func(ctx context.Context) error {
		return job(c)
	}
}
