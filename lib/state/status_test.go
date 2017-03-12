package state_test

import (
	"testing"

	"github.com/apoydence/onpar"
	"github.com/bradylove/concourse-monitor/lib/concourse"
	"github.com/bradylove/concourse-monitor/lib/state"

	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

func TestStatus(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("JobToStatus()", func() {
		o.Group("with no next build", func() {
			o.Spec("it returns the status of finished build", func(t *testing.T) {
				testset := []string{
					state.StatusPending,
					state.StatusStarted,
					state.StatusSucceeded,
					state.StatusFailed,
					state.StatusErrored,
					state.StatusAborted,
				}

				for _, st := range testset {
					job := &concourse.Job{
						FinishedBuild: &concourse.Build{Status: st},
					}

					Expect(t, state.JobStatus(job)).To(Equal(st))
				}
			})
		})

		o.Group("with next build", func() {
			o.Spec("it returns the status of next build", func(t *testing.T) {
				testset := []string{
					state.StatusPending,
					state.StatusStarted,
					state.StatusSucceeded,
					state.StatusFailed,
					state.StatusErrored,
					state.StatusAborted,
				}

				for _, st := range testset {
					job := &concourse.Job{
						NextBuild:     &concourse.Build{Status: st},
						FinishedBuild: &concourse.Build{Status: state.StatusSucceeded},
					}

					Expect(t, state.JobStatus(job)).To(Equal(st))
				}
			})
		})

		o.Group("with a paused job", func() {
			o.Spec("it returns the status of finished build", func(t *testing.T) {
				testset := []string{
					state.StatusPending,
					state.StatusStarted,
					state.StatusSucceeded,
					state.StatusFailed,
					state.StatusErrored,
					state.StatusAborted,
				}

				for _, st := range testset {
					job := &concourse.Job{
						Paused:        true,
						FinishedBuild: &concourse.Build{Status: st},
					}

					Expect(t, state.JobStatus(job)).To(Equal(state.StatusPaused))
				}
			})
		})

		o.Group("with no next or finished build", func() {
			o.Spec("it returns the status of next build", func(t *testing.T) {
				job := &concourse.Job{}

				Expect(t, state.JobStatus(job)).To(Equal(state.StatusUnknown))
			})
		})
	})

	o.Group("PipelineToStatus()", func() {
		o.Group("with no failed, aborted or errored jobs", func() {
			o.Spec("it returns succeeded", func(t *testing.T) {
				p := &concourse.Pipeline{
					Jobs: []*concourse.Job{
						{FinishedBuild: &concourse.Build{Status: state.StatusSucceeded}},
						{NextBuild: &concourse.Build{Status: state.StatusStarted}},
						{NextBuild: &concourse.Build{Status: state.StatusPending}},
					},
				}

				Expect(t, state.PipelineStatus(p)).To(Equal(state.StatusSucceeded))
			})
		})

		o.Group("with a pipeline that is paused", func() {
			o.Spec("it returns paused", func(t *testing.T) {
				p := &concourse.Pipeline{
					Paused: true,
					Jobs: []*concourse.Job{
						{FinishedBuild: &concourse.Build{Status: state.StatusSucceeded}},
					},
				}

				Expect(t, state.PipelineStatus(p)).To(Equal(state.StatusPaused))
			})
		})

		o.Group("with a failed job", func() {
			o.Spec("it returns failed", func(t *testing.T) {
				p := &concourse.Pipeline{
					Jobs: []*concourse.Job{
						{FinishedBuild: &concourse.Build{Status: state.StatusSucceeded}},
						{FinishedBuild: &concourse.Build{Status: state.StatusAborted}},
						{FinishedBuild: &concourse.Build{Status: state.StatusFailed}},
						{FinishedBuild: &concourse.Build{Status: state.StatusErrored}},
					},
				}

				Expect(t, state.PipelineStatus(p)).To(Equal(state.StatusFailed))
			})
		})

		o.Group("with an errored job", func() {
			o.Spec("it returns errored", func(t *testing.T) {
				p := &concourse.Pipeline{
					Jobs: []*concourse.Job{
						{FinishedBuild: &concourse.Build{Status: state.StatusSucceeded}},
						{FinishedBuild: &concourse.Build{Status: state.StatusAborted}},
						{FinishedBuild: &concourse.Build{Status: state.StatusErrored}},
					},
				}

				Expect(t, state.PipelineStatus(p)).To(Equal(state.StatusErrored))
			})
		})

		o.Group("with an aborted job", func() {
			o.Spec("it returns aborted", func(t *testing.T) {
				p := &concourse.Pipeline{
					Jobs: []*concourse.Job{
						{FinishedBuild: &concourse.Build{Status: state.StatusSucceeded}},
						{FinishedBuild: &concourse.Build{Status: state.StatusAborted}},
					},
				}

				Expect(t, state.PipelineStatus(p)).To(Equal(state.StatusAborted))
			})
		})
	})
}
