package safaErrorGroup

import (
	"context"
	"fmt"
	"github.com/halalala222/try-safe-wait-group/safe/multiErrors"
	"golang.org/x/sync/semaphore"
	"sync"
)

type ErrorGroup struct {
	wg     sync.WaitGroup
	errors multiErrors.MultiErrors
	sem    *semaphore.Weighted
	ctx    context.Context
}

type taskFunc func(ctx context.Context) error

func (e *ErrorGroup) New(ctx context.Context, maxWorkers int64) *ErrorGroup {
	return &ErrorGroup{
		ctx:    ctx,
		sem:    semaphore.NewWeighted(maxWorkers),
		errors: multiErrors.Cap(maxWorkers),
	}
}

func (e *ErrorGroup) Do(task taskFunc) {
	e.wg.Add(1)

	if !e.checkSemAcquire() {
		return
	}

	go func() {
		e.work(task)
	}()
}

func (e *ErrorGroup) checkSemAcquire() bool {
	if err := e.sem.Acquire(e.ctx, 1); err != nil {
		defer e.wg.Done()
		e.errors = append(e.errors, fmt.Errorf("couldn't acquire semaphore : %s", err))
		return false
	}
	return true
}

func (e *ErrorGroup) work(task taskFunc) {
	defer func() {
		if recovery := recover(); recovery != nil {
			e.errors = append(e.errors, fmt.Errorf("panic erro : %v", recovery))
		}
		e.wg.Done()
		e.sem.Release(1)
	}()
	if err := task(e.ctx); err != nil {
		e.errors = append(e.errors, err)
	}
}

func (e *ErrorGroup) Wait() error {
	e.wg.Wait()
	return e.errors.ErrorOrNil()
}
