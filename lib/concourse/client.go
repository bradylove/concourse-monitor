package concourse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	pipelinesPath = "/api/v1/teams/%s/pipelines"
	jobPath       = "/api/v1/teams/%s/pipelines/%s/jobs"
)

type Client struct {
	targets []Target
}

type Pipeline struct{}

func NewClient(targets []Target) (*Client, error) {
	return &Client{targets}, nil
}

func (c *Client) Pipelines() ([]Pipeline, error) {
	pipes := make([]Pipeline, 0)

	for _, t := range c.targets {
		p, err := c.requestPipeline(t.API, t.Team)
		if err != nil {
			return nil, err
		}
		pipes = append(pipes, p...)
	}

	return pipes, nil
}

func (c *Client) requestPipeline(host, team string) ([]Pipeline, error) {
	resp, err := http.Get(host + fmt.Sprintf(pipelinesPath, team))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 response code, got %d", resp.StatusCode)
	}

	var p []Pipeline
	err = json.NewDecoder(resp.Body).Decode(&p)

	return p, err
}
