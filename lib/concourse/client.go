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
	addr     string
	teamName string
}

type Pipeline struct{}

func NewClient(addr, teamName string) (*Client, error) {
	return &Client{addr: addr, teamName: teamName}, nil
}

func (c *Client) Pipelines() ([]*Pipeline, error) {
	resp, err := http.Get(c.addr + fmt.Sprintf(pipelinesPath, c.teamName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 response code, got %d", resp.StatusCode)
	}

	var pipes []*Pipeline
	err = json.NewDecoder(resp.Body).Decode(&pipes)
	if err != nil {
		return nil, err
	}

	return pipes, nil
}
