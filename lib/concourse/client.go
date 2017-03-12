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
	Paused        bool   `json:"paused"`
	NextBuild     *Build `json:"next_build"`
	FinishedBuild *Build `json:"finished_build"`
}

type Build struct {
	Status       string `json:"status"`
	JobName      string `json:"job_name"`
	PipelineName string `json:"pipeline_name"`
}

func NewClient(targets []Target) *Client {
	return &Client{targets}
}

func (c *Client) Pipelines() ([]*Pipeline, error) {
	pipes := make([]*Pipeline, 0)

	for _, t := range c.targets {
		p, err := c.requestPipeline(t.API, t.Team, t.Token)
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, p...)
	}

	return pipes, nil
}

func (c *Client) requestPipeline(host, team string, token Token) ([]*Pipeline, error) {
	req, err := http.NewRequest(http.MethodGet, host+fmt.Sprintf(pipelinesPath, team), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.Value))

	resp, err := http.DefaultClient.Do(req)
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
		j, err := c.requestJobs(host, p, token)
		if err != nil {
			return nil, err
		}

		p.Jobs = j
	}

	return pipelines, nil
}

func (c *Client) requestJobs(host string, p *Pipeline, token Token) ([]*Job, error) {
	req, err := http.NewRequest(http.MethodGet, host+fmt.Sprintf(jobsPath, p.TeamName, p.Name), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.Type, token.Value))

	resp, err := http.DefaultClient.Do(req)
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
