package concourse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pipelinesPath = "/api/v1/teams/%s/pipelines"
	jobsPath      = "/api/v1/teams/%s/pipelines/%s/jobs"
)

type Client struct {
	targets []Target
}

type Pipeline struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	TeamName string `json:"team_name"`
	Paused   bool   `json:"paused"`
	Jobs     []*Job
}

type Job struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	NextBuild     *Build `json:"next_build"`
	FinishedBuild *Build `json:"finished_build"`
}

type Build struct {
	Status string `json:"status"`
}

func NewClient(targets []Target) (*Client, error) {
	return &Client{targets}, nil
}

func (c *Client) Pipelines() ([]*Pipeline, error) {
	pipes := make([]*Pipeline, 0)

	for _, t := range c.targets {
		p, err := c.requestPipeline(t.API, t.Team)
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, p...)
	}

	return pipes, nil
}

func (c *Client) requestPipeline(host, team string) ([]*Pipeline, error) {
	resp, err := http.Get(host + fmt.Sprintf(pipelinesPath, team))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 response code, got %d", resp.StatusCode)
	}

	var pipelines []*Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, err
	}

	for _, p := range pipelines {
		j, err := c.requestJobs(host, p)
		if err != nil {
			return nil, err
		}

		p.Jobs = j
	}

	return pipelines, nil
}

func (c *Client) requestJobs(host string, p *Pipeline) ([]*Job, error) {
	resp, err := http.Get(host + fmt.Sprintf(jobsPath, p.TeamName, p.Name))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 response code, got %d", resp.StatusCode)
	}

	var jobs []*Job
	if err := json.NewDecoder(resp.Body).Decode(&jobs); err != nil {
		return nil, err
	}

	return jobs, nil
}
