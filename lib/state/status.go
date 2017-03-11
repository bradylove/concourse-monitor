package state

import "github.com/bradylove/concourse-monitor/lib/concourse"

const (
	StatusPending   = "pending"
	StatusStarted   = "started"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
	StatusErrored   = "errored"
	StatusAborted   = "aborted"
	StatusPaused    = "paused"
	StatusUnknown   = "unknown"
)

var (
	failurePriority = []string{StatusAborted, StatusErrored, StatusFailed}
)

func PipelineToStatus(p *concourse.Pipeline) string {
	if p.Paused {
		return StatusPaused
	}

	failure := -1
	for _, j := range p.Jobs {
		priority := stateToPriority(JobToStatus(j))

		if priority > failure {
			failure = priority
		}
	}

	if failure >= 0 {
		return failurePriority[failure]
	}

	return StatusSucceeded
}

func JobToStatus(j *concourse.Job) string {
	if j.Paused {
		return StatusPaused
	}

	if j.NextBuild != nil {
		return j.NextBuild.Status
	}

	if j.FinishedBuild == nil {
		return StatusUnknown
	}

	return j.FinishedBuild.Status
}

func stateToPriority(status string) int {
	for i, s := range failurePriority {
		if s == status {
			return i
		}
	}

	return -1
}
