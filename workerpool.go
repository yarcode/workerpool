/*
Package workerpool provides a service for running small parts of code (called jobs) in a background.

Jobs could have contexts, timeouts, rich retry strategies.
*/
package workerpool

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Job is a function that receives context and being run asynchronously by the worker pool.
type Job func(ctx context.Context) error

type Pool struct {
	wg     *sync.WaitGroup
	logger zerolog.Logger

	jobs       chan Job
	stop       chan struct{}
	numWorkers int

	// DefaultContext is a context factory, generates default context for every job
	DefaultContext func() context.Context
}

// PoolOption is a functional parameter used in constructor
type PoolOption func(service *Pool)

// New Pool constructor
func New(opts ...PoolOption) *Pool {
	s := &Pool{
		wg:         &sync.WaitGroup{},
		logger:     log.Logger,
		numWorkers: runtime.GOMAXPROCS(0) * 2,

		DefaultContext: func() context.Context {
			return context.Background()
		},
	}

	for _, o := range opts {
		o(s)
	}

	return s
}

// WithLogger replaces the default logger with the specified one.
func WithLogger(zl zerolog.Logger) PoolOption {
	return func(service *Pool) {
		service.logger = zl
	}
}

// Start required num of workers.
func (s *Pool) Start() {

	s.logger.Debug().Int("worker_num", s.numWorkers).Msg("Starting workers")
	s.stop = make(chan struct{})
	s.jobs = make(chan Job)
	s.wg.Add(s.numWorkers)
	for i := 0; i < s.numWorkers; i++ {
		go func(workerID int) {
			defer s.wg.Done()
			l := s.logger.With().Int("worker_id", workerID).Logger()
			for {
				select {
				case <-s.stop:
					return
				case job := <-s.jobs:
					id := xid.New()
					now := time.Now()
					ll := l.With().Str("job_id", id.String()).Logger()
					ll.Debug().Msg("Running job")
					if err := AddLogger(AddPanicRecovery(job), ll)(s.DefaultContext()); err != nil {
						ll.Error().Err(err).Msg("Error running job")
					}
					ll.Info().Dur("job_duration", time.Since(now)).Msg("Done running job")
				}
			}
		}(i)
	}

	s.logger.Debug().Msg("Done starting workers")
}

// Stop all active workers.
// Don't forget to call this on your application shutdown.
//
// This call is blocking and waits for all the jobs to complete.
func (s *Pool) Stop() {
	s.logger.Debug().Msg("Shutting down workers")
	close(s.stop)
	s.wg.Wait()
	close(s.jobs)
	s.logger.Debug().Msg("Done shutting down workers")
}

func (s *Pool) Run(job Job) {
	s.jobs <- job
}
