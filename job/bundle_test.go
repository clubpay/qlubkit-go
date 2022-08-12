package job_test

import (
	"context"
	"strings"
	"sync"
	"testing"

	"github.com/clubpay/qlubkit-go/job"
	. "github.com/smartystreets/goconvey/convey"
)

type testJob struct {
	cg        *caseGen
	id        int64
	retry     int
	f         job.Func
	dependsOn []string
}

func (t *testJob) Retry(ctx context.Context, err error) bool {
	if t.retry > 0 {
		t.retry--

		return true
	}

	return false
}

func (t testJob) ID() int64 {
	return t.id
}

func (t testJob) Func() job.Func {
	return t.f
}

func (t testJob) RunAfter() []int64 {
	var out []int64
	for _, n := range t.dependsOn {
		out = append(out, t.cg.jobs[n].ID())
	}

	return out
}

type caseGen struct {
	c      C
	nextID int64
	doneL  sync.Mutex
	done   map[string]struct{}
	jobs   map[string]job.Job
	buf    strings.Builder
}

func newCaseGen(c C) *caseGen {
	return &caseGen{
		c:    c,
		done: map[string]struct{}{},
		jobs: map[string]job.Job{},
	}
}

func (cg *caseGen) Job(name string, dependsOn ...string) {
	cg.nextID++
	j := &testJob{
		cg:    cg,
		id:    cg.nextID,
		retry: 3,
		f: func(ctx *job.Context) error {
			cg.buf.WriteString(name)
			cg.doneL.Lock()
			cg.done[name] = struct{}{}
			cg.doneL.Unlock()

			return nil
		},
		dependsOn: dependsOn,
	}

	cg.jobs[name] = j
}

func (cg *caseGen) AddJobs(b *job.Bundle) {
	for _, j := range cg.jobs {
		b.AddJob(j)
	}
}

func TestBundle(t *testing.T) {
	Convey("Job Bundle", t, func(c C) {
		Convey("Success Case", func(c C) {
			cg := newCaseGen(c)
			cg.Job("J1")
			cg.Job("J2")
			cg.Job("J3", "J1", "J2")
			cg.Job("J4", "J3")
			cg.Job("J5", "J2", "J7", "J4")
			cg.Job("J6", "J3", "J5")
			cg.Job("J7", "J3", "J4")

			b := job.NewBundle(1)
			cg.AddJobs(b)
			b.Do(context.Background())
			c.So(cg.buf.String(), ShouldEqual, "J1J2J3J4J7J5J6")
		})
		Convey("Panic Case - Cyclic Dependency", func(c C) {
			cg := newCaseGen(c)
			cg.Job("J1")
			cg.Job("J2")
			cg.Job("J3", "J1", "J2")
			cg.Job("J4", "J1", "J5")
			cg.Job("J5", "J4", "J7")
			cg.Job("J6", "J3", "J5")
			cg.Job("J7", "J3", "J6")

			b := job.NewBundle(3)
			cg.AddJobs(b)
			f := func() { b.Do(context.Background()) }
			c.So(f, ShouldPanic)
		})
		Convey("Panic Case - Missing Job", func(c C) {
			cg := newCaseGen(c)
			cg.Job("J2")
			cg.Job("J3", "J1", "J2")
			cg.Job("J4", "J1", "J5")
			cg.Job("J5", "J4", "J7")
			cg.Job("J6", "J3", "J5")
			cg.Job("J7", "J3")

			b := job.NewBundle(3)
			cg.AddJobs(b)
			f := func() { b.Do(context.Background()) }
			c.So(f, ShouldPanic)
		})
	})
}
