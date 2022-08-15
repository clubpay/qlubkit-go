package job_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/clubpay/qlubkit-go/job"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBundle(t *testing.T) {
	dummyTask := func(txt string, w io.Writer, latency time.Duration) job.Task {
		return func(ctx *job.Context) error {
			time.Sleep(latency)

			_, err := io.WriteString(w, txt)
			
			return err
		}
	}
	Convey("Job Bundle", t, func(c C) {
		Convey("Success Case 1", func(c C) {
			buf := &strings.Builder{}
			var jobs []job.Job
			for i := 0; i < 7; i++ {
				jobs = append(jobs,
					job.NewJob(
						fmt.Sprintf("J%d", i+1),
						job.WithMaxRetry(3),
					).AddTask(
						dummyTask(
							fmt.Sprintf("J%d", i+1),
							buf,
							time.Millisecond*100,
						),
					),
				)
			}

			b := job.NewBundle(1)
			b.AddJob(jobs...)
			b.Relate(jobs[2], job.After, jobs[1], jobs[0])
			b.Relate(jobs[3], job.After, jobs[2])
			b.Relate(jobs[4], job.After, jobs[1], jobs[6], jobs[3])
			b.Relate(jobs[5], job.After, jobs[2], jobs[4])
			b.Relate(jobs[6], job.After, jobs[2], jobs[3])

			b.Do(context.Background())
			c.So(buf.String(), ShouldEqual, "J1J2J3J4J7J5J6")
		})

		Convey("Success Case 1 (Concurrent)", func(c C) {
			buf := &strings.Builder{}
			var jobs []job.Job
			for i := 0; i < 7; i++ {
				jobs = append(jobs,
					job.NewJob(
						fmt.Sprintf("J%d", i+1),
						job.WithMaxRetry(3),
					).AddTask(
						dummyTask(
							fmt.Sprintf("J%d", i+1),
							buf,
							time.Millisecond*time.Duration(100*(i+1)),
						),
					),
				)
			}

			b := job.NewBundle(3)
			b.AddJob(jobs...)
			b.Relate(jobs[2], job.After, jobs[1], jobs[0])
			b.Relate(jobs[3], job.After, jobs[2])
			b.Relate(jobs[4], job.After, jobs[1], jobs[6], jobs[3])
			b.Relate(jobs[5], job.After, jobs[2], jobs[4])
			b.Relate(jobs[6], job.After, jobs[2], jobs[3])

			b.Do(context.Background())
			c.So(buf.String(), ShouldEqual, "J1J2J3J4J7J5J6")
		})

		Convey("Success Case 2", func(c C) {
			buf := &strings.Builder{}
			var jobs []job.Job
			for i := 0; i < 7; i++ {
				jobs = append(jobs,
					job.NewJob(
						fmt.Sprintf("J%d", i+1),
						job.WithMaxRetry(3),
					).AddTask(
						dummyTask(
							fmt.Sprintf("J%d", i+1),
							buf,
							time.Millisecond,
						),
					),
				)
			}

			b := job.NewBundle(1)
			b.AddJob(jobs...)
			b.Relate(jobs[0], job.After, jobs[1])
			b.Relate(jobs[2], job.After, jobs[3])
			b.Relate(jobs[3], job.After, jobs[1])
			b.Relate(jobs[4], job.After, jobs[1], jobs[6], jobs[3])
			b.Relate(jobs[5], job.After, jobs[2], jobs[4])
			b.Relate(jobs[6], job.After, jobs[2], jobs[3])

			b.Do(context.Background())
			c.So(buf.String(), ShouldEqual, "J2J1J4J3J7J5J6")
		})

		Convey("Success Case 2 (Concurrent)", func(c C) {
			buf := &strings.Builder{}
			var jobs []job.Job
			for i := 0; i < 7; i++ {
				jobs = append(jobs,
					job.NewJob(
						fmt.Sprintf("J%d", i+1),
						job.WithMaxRetry(3),
					).AddTask(
						dummyTask(
							fmt.Sprintf("J%d", i+1),
							buf,
							time.Millisecond,
						),
					),
				)
			}

			b := job.NewBundle(3)
			b.AddJob(jobs...)
			b.Relate(jobs[0], job.After, jobs[1])
			b.Relate(jobs[2], job.After, jobs[3])
			b.Relate(jobs[3], job.After, jobs[1])
			b.Relate(jobs[4], job.After, jobs[1], jobs[6], jobs[3])
			b.Relate(jobs[5], job.After, jobs[2], jobs[4])
			b.Relate(jobs[6], job.After, jobs[2], jobs[3])

			b.Do(context.Background())
			c.So(buf.String(), ShouldEqual, "J2J1J4J3J7J5J6")
		})

		Convey("Panic Case - Cyclic Dependency", func(c C) {
			buf := &strings.Builder{}
			var jobs []job.Job
			for i := 0; i < 7; i++ {
				jobs = append(jobs,
					job.NewJob(
						fmt.Sprintf("J%d", i+1),
						job.WithMaxRetry(3),
					).AddTask(
						dummyTask(
							fmt.Sprintf("J%d", i+1),
							buf,
							time.Millisecond,
						),
					),
				)
			}

			b := job.NewBundle(1)
			b.AddJob(jobs...)
			b.Relate(jobs[0], job.After, jobs[1])
			b.Relate(jobs[1], job.After, jobs[2])
			b.Relate(jobs[2], job.After, jobs[0])

			c.So(func() { b.Do(context.Background()) }, ShouldPanic)
		})
	})
}
