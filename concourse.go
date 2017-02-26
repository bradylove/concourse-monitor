package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net/http"

	"github.com/bradylove/concourse-monitor/assets"
)

const (
	pipelinesPath = "/api/v1/teams/%s/pipelines"
	jobPath       = "/api/v1/teams/%s/pipelines/%s/jobs"
)

var (
	statusIconSucceeded = assets.Image("icons/success.png")
	statusIconStarted   = assets.Image("icons/started.png")
	statusIconFailed    = assets.Image("icons/failed.png")
	statusIconPaused    = assets.Image("icons/paused.png")
	statusIconAborted   = assets.Image("icons/aborted.png")
	statusIconPending   = assets.Image("icons/pending.png")
	statusIconErrored   = assets.Image("icons/errored.png")

	statusSucceeded = "succeeded"
	statusStarted   = "started"
	statusFailed    = "failed"
	statusPaused    = "paused"
	statusAborted   = "aborted"
	statusPending   = "pending"
	statusErrored   = "errored"
)

type Team struct {
	ID        uint        `json:"id"`
	Name      string      `json:"name"`
	Pipelines []*Pipeline `json:"-"`
}

type Pipeline struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	TeamName string `json:"team_name"`
	Paused   bool   `json:"paused"`
	Public   bool   `json:"public"`
	Jobs     []*Job `json:"-"`
}

func (p *Pipeline) StatusIcon() image.Image {
	for _, j := range p.Jobs {
		if j.Status() == statusFailed {
			return statusIconFailed
		}
	}

	return statusIconSucceeded
}

type Job struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Paused        bool   `json:"paused"`
	FinishedBuild *Build `json:"finished_build"`
	NextBuild     *Build `json:"next_build"`
}

func (j *Job) Status() string {
	if j.NextBuild != nil {
		return j.NextBuild.Status
	}

	return j.FinishedBuild.Status
}

func (j *Job) StatusIcon() image.Image {
	if j.Paused {
		return statusIconPaused
	}

	switch j.Status() {
	case statusSucceeded:
		return statusIconSucceeded
	case statusFailed:
		return statusIconFailed
	case statusPaused:
		return statusIconPaused
	case statusAborted:
		return statusIconAborted
	case statusPending:
		return statusIconPending
	case statusErrored:
		return statusIconErrored
	case statusStarted:
		return statusIconStarted
	}

	log.Fatalf("Unhandled state: %s", j.Status())

	return nil
}

type Build struct {
	ID           uint   `json:"id"`
	TeamName     string `json:"team_name"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	JobName      string `json:"job_name"`
	URL          string `json:"url"`
	APIURL       string `json:"api_url"`
	PipelineName string `json:"pipeline_name"`
	StartTime    int64  `json:"start_time"`
	EndTime      int64  `json:"end_time"`
}

type ConcourseClient struct {
	addr string
}

func NewConcourseClient(addr string) *ConcourseClient {
	return &ConcourseClient{addr: addr}
}

func (c *ConcourseClient) GetPipelines(teamName string) ([]*Pipeline, error) {
	resp, err := http.Get(fmt.Sprintf(c.addr+pipelinesPath, teamName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var pipelines []*Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	for _, p := range pipelines {
		jobs, err := c.GetJobs(teamName, p.Name)
		if err != nil {
			log.Println("Failed to get jobs: %s", err)
			continue
		}

		p.Jobs = jobs
	}

	return pipelines, err
}

func (c *ConcourseClient) GetJobs(teamName, pipelineName string) ([]*Job, error) {
	resp, err := http.Get(fmt.Sprintf(c.addr+jobPath, teamName, pipelineName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jobs []*Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return jobs, err
}
