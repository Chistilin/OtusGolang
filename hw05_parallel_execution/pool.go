package hw05parallelexecution

import (
	"context"
	"sync"
	"sync/atomic"
)

type limiter struct {
	count int64
	limit int64
}

func newLimiter(limit int64) *limiter {
	return &limiter{
		limit: limit,
	}
}

func (l *limiter) inc() {
	atomic.AddInt64(&l.count, 1)
}

func (l *limiter) isLimitExceeded() bool {
	return atomic.LoadInt64(&l.count) >= atomic.LoadInt64(&l.limit)
}

type pool struct {
	tasks       []Task
	concurrency int
	collector   chan Task
	wg          sync.WaitGroup
	limiter     *limiter
}

func newPool(tasks []Task, concurrency int, maxErrorCount int64) *pool {
	return &pool{
		tasks:       tasks,
		concurrency: concurrency,
		collector:   make(chan Task, len(tasks)),
		limiter:     newLimiter(maxErrorCount),
	}
}

func (p *pool) run(ctx context.Context) error {
	// Run workers
	for i := 0; i < p.concurrency; i++ {
		w := newWorker(i, p.collector)
		w.Start(&p.wg, p.limiter)
	}

	// Канал для сигнализации о завершении
	done := make(chan struct{})

	// Отправка задач в отдельной горутине
	go func() {
		defer close(p.collector)
		for _, task := range p.tasks {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			case p.collector <- task:
			}
		}
	}()

	p.wg.Wait()

	close(done)

	if p.limiter.isLimitExceeded() {
		return ErrErrorsLimitExceeded
	}

	return nil
}
