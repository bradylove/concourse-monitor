package main

import (
	"fmt"

	"github.com/0xAX/notificator"
	"github.com/bradylove/concourse-monitor/lib/concourse"
	"github.com/bradylove/concourse-monitor/lib/state"
)

const (
	notifyMsg = "%s/%s status was %s and is now %s."
)

type Notifier interface {
	Push(string, string, string, string) error
}

type Cache struct {
	store    map[string]*concourse.Job
	notifier Notifier
}

func NewCache(n Notifier) *Cache {
	return &Cache{
		store:    make(map[string]*concourse.Job),
		notifier: n,
	}
}

func (c *Cache) Update(pipelines []*concourse.Pipeline) {
	for _, p := range pipelines {
		for _, j := range p.Jobs {
			oldJob, ok := c.store[j.URL]
			if !ok {
				c.store[j.URL] = j
				continue
			}

			c.NotifyOnBuildDiff(oldJob.FinishedBuild, j.FinishedBuild)

			c.store[j.URL] = j
		}
	}
}

func (c *Cache) Notify(msg string) {
	c.notifier.Push(
		"Concourse Monitor",
		msg,
		"",
		notificator.UR_NORMAL,
	)
}

func (c *Cache) NotifyOnBuildDiff(old, new *concourse.Build) {
	// TODO: Remove this and handle gracefully
	if old == nil {
		return
	}

	if new == nil {
		return
	}

	if old.Status != state.StatusSucceeded && new.Status == state.StatusSucceeded {
		c.Notify(fmt.Sprintf(notifyMsg, new.PipelineName, new.JobName, old.Status, new.Status))
	}

	if old.Status == state.StatusSucceeded && new.Status != state.StatusSucceeded {
		c.Notify(fmt.Sprintf(notifyMsg, new.PipelineName, new.JobName, old.Status, new.Status))
	}
}
