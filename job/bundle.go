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

func (s state) String() string {
	switch s {
	case initiated:
		return "initiated"
	case ready:
		return "ready"
	case running:
		return "running"
	case done:
		return "done"
	}

	panic("invalid state")
}

type Bundle struct {
	concurrency int
	relation    map[int64]map[int64]state
	jobs        []Job
	states      []state
}

func NewBundle(concurrency int) *Bundle {
	return &Bundle{
		concurrency: concurrency,
		relation:    map[int64]map[int64]state{},
	}
}

func (b *Bundle) AddJob(job ...Job) *Bundle {
	b.jobs = append(b.jobs, job...)

	return b
}

func (b *Bundle) Relate(job Job, r RelationMaker, jobs ...Job) {
	r(b).relates(job, jobs...)
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
			jobCtx := newBag(ctx)
			job := b.jobs[idx]
		Outer:
			for _, task := range job.Tasks() {
			Inner:
				for {
					err := task(jobCtx)
					if err == nil {
						break
					}

					switch job.OnError(ctx, err) {
					case Retry:
					case IgnoreAndContinue:
						break Inner
					case StopAndExit:
						break Outer
					}
				}
			}
			b.states[idx] = done

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
	for jobID, jobState := range b.relation[j.ID()] {
		for i := 0; i < len(b.jobs); i++ {
			if b.jobs[i].ID() == jobID && b.states[i] < jobState {
				return false
			}
		}
	}

	return true
}

type RelationMaker = func(b *Bundle) Relation

type Relation interface {
	relates(job Job, others ...Job)
}

type notBefore struct {
	b *Bundle
}

func (n notBefore) relates(job Job, jobs ...Job) {
	rel := n.b.relation[job.ID()]
	if rel == nil {
		rel = map[int64]state{}
	}

	for _, j := range jobs {
		rel[j.ID()] = running
	}

	n.b.relation[job.ID()] = rel
}

func NoBefore(b *Bundle) Relation {
	r := &notBefore{b: b}

	return r
}

type after struct {
	b *Bundle
}

func (n after) relates(job Job, jobs ...Job) {
	rel := n.b.relation[job.ID()]
	if rel == nil {
		rel = map[int64]state{}
	}

	for _, j := range jobs {
		rel[j.ID()] = done
	}

	n.b.relation[job.ID()] = rel
}

func After(b *Bundle) Relation {
	r := &after{b: b}

	return r
}
