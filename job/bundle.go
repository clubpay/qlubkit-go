package job

import (
	"context"
	"sync"
)

type state int

const (
	initiated state = iota
	ready
	running
	done
)

type Func func(ctx *Context) error

type Bundle struct {
	concurrency int

	jobs   []Job
	states []state
}

func NewBundle(concurrency int) *Bundle {
	return &Bundle{
		concurrency: concurrency,
	}
}

func (b *Bundle) AddJob(j Job) *Bundle {
	b.jobs = append(b.jobs, j)

	return b
}

func (b *Bundle) Do(ctx context.Context) {
	totalJobs := len(b.jobs)
	totalDone := 0

	// if there is no job just return quickly.
	if totalJobs == 0 {
		return
	}

	b.states = make([]state, totalJobs)
	rateLimit := make(chan struct{}, b.concurrency)
	wg := sync.WaitGroup{}
	jobCtx := newCtx(ctx)

MainLoop:
	for cur := 0; cur < totalJobs; cur++ {
		switch b.states[cur] {
		default:
			panic("BUG!! we should never have unknown state")
		case initiated:
			// check if state is ready to run if it is runnable
			if b.isRunnable(b.jobs[cur]) {
				b.states[cur] = ready
			} else {
				continue
			}
		case done, running:
			continue
		case ready:
		}

		// run the job in background, it will retry if it fails.
		// we set the state to running, and will set the state done
		// if job is finished successfully or maximum retries reached.
		rateLimit <- struct{}{}
		b.states[cur] = running
		wg.Add(1)

		go func(idx int) {
			for {
				err := b.jobs[idx].Func()(jobCtx)
				if err == nil || !b.jobs[idx].Retry(ctx, err) {
					b.states[idx] = done

					break
				} else {
					b.states[idx] = ready
				}
			}

			<-rateLimit
			wg.Done()
		}(cur)
	}
	wg.Wait()

	newTotalDone := b.totalDone()
	if newTotalDone == totalDone && totalDone != totalJobs {
		panic("BUG!! there is cyclic dependency between jobs")
	}

	totalDone = newTotalDone
	for cur := 0; cur < totalJobs; cur++ {
		switch b.states[cur] {
		case ready, initiated:
			goto MainLoop
		}
	}

}

func (b *Bundle) totalDone() int {
	doneCnt := 0
	for cur := 0; cur < len(b.jobs); cur++ {
		if b.states[cur] == done {
			doneCnt++
		}
	}

	return doneCnt
}

func (b *Bundle) isRunnable(j Job) bool {
	for _, jobID := range j.RunAfter() {
		for i := 0; i < len(b.jobs); i++ {
			if b.jobs[i].ID() == jobID && b.states[i] != done {

				return false
			}
		}

	}

	return true
}
