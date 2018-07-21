package concourse_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apoydence/onpar"
	"github.com/bradylove/concourse-monitor/pkg/concourse"

	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

func TestConcourseClient(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("Pipelines()", func() {
		o.Group("with successful response", func() {
			o.BeforeEach(func(t *testing.T) (*testing.T, *FakeConcourse) {
				fc := NewFakeConcourse()

				return t, fc
			})

			o.AfterEach(func(t *testing.T, fc *FakeConcourse) {
				fc.Close()
			})

			o.Spec("it returns all pipelines and jobs", func(t *testing.T, fc *FakeConcourse) {
				targets := []concourse.Target{
					{Name: "one", API: fc.URL, Team: "main", Token: concourse.Token{Type: "Bearer", Value: "main-token"}},
					{Name: "two", API: fc.URL, Team: "awesome", Token: concourse.Token{Type: "Bearer", Value: "awesome-token"}},
				}
				client := concourse.NewClient(targets)

				pipes, err := client.Pipelines()
				Expect(t, err).To(Not(HaveOccurred()))
				Expect(t, pipes).To(HaveLen(2))

				p1 := pipes[0]
				Expect(t, p1.DisplayName).To(Equal("one/pipeline-1"))
				Expect(t, p1.Name).To(Equal("pipeline-1"))
				Expect(t, p1.URL).To(Equal("/teams/main/pipelines/pipeline-1"))
				Expect(t, p1.Paused).To(Equal(false))
				Expect(t, p1.TeamName).To(Equal("main"))
				Expect(t, p1.Jobs).To(HaveLen(1))
				Expect(t, p1.Target).To(Equal(&targets[0]))

				j1 := p1.Jobs[0]
				Expect(t, j1.Name).To(Equal("hello-world"))
				Expect(t, j1.Paused).To(BeTrue())
				Expect(t, j1.URL).To(Equal("/teams/main/pipelines/pipeline-1/jobs/hello-world"))

				b1 := j1.FinishedBuild
				Expect(t, b1.Status).To(Equal("failed"))

				req := <-fc.requests
				Expect(t, req.URL.Path).To(Equal("/api/v1/teams/main/pipelines"))
				Expect(t, req.Header.Get("Authorization")).To(Equal("Bearer main-token"))

				req = <-fc.requests
				Expect(t, req.URL.Path).To(Equal("/api/v1/teams/main/pipelines/pipeline-1/jobs"))
				Expect(t, req.Header.Get("Authorization")).To(Equal("Bearer main-token"))

				req = <-fc.requests
				Expect(t, req.URL.Path).To(Equal("/api/v1/teams/awesome/pipelines"))
				Expect(t, req.Header.Get("Authorization")).To(Equal("Bearer awesome-token"))

				req = <-fc.requests
				Expect(t, req.URL.Path).To(Equal("/api/v1/teams/main/pipelines/pipeline-1/jobs"))
				Expect(t, req.Header.Get("Authorization")).To(Equal("Bearer awesome-token"))
			})
		})

		o.Group("with non 200 response code", func() {
			o.BeforeEach(func(t *testing.T) (*testing.T, *httptest.Server) {
				handler := func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}
				s := httptest.NewServer(http.HandlerFunc(handler))

				return t, s
			})

			o.AfterEach(func(t *testing.T, s *httptest.Server) {
				s.Close()
			})

			o.Spec("it returns an error", func(t *testing.T, s *httptest.Server) {
				targets := []concourse.Target{
					{API: s.URL, Team: "main"},
				}
				client := concourse.NewClient(targets)

				_, err := client.Pipelines()
				Expect(t, err.Error()).To(Equal("expected 200 response code, got 500"))
			})
		})

		o.Group("with a request error", func() {
			o.Spec("it returns an error", func(t *testing.T) {
				targets := []concourse.Target{
					{API: "http://127.0.0.1:23223", Team: "main"},
				}
				client := concourse.NewClient(targets)

				_, err := client.Pipelines()
				Expect(t, err).To(HaveOccurred())
			})
		})
	})
}

type FakeConcourse struct {
	*httptest.Server

	requests chan *http.Request
}

func NewFakeConcourse() *FakeConcourse {
	f := &FakeConcourse{
		requests: make(chan *http.Request, 100),
	}

	f.Server = httptest.NewServer(f)

	return f
}

func (f *FakeConcourse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.requests <- r

	w.WriteHeader(http.StatusOK)
	if r.URL.Path == "/api/v1/teams/main/pipelines" ||
		r.URL.Path == "/api/v1/teams/awesome/pipelines" {

		w.Write([]byte(getPipelinesResponse))
		return
	}

	if r.URL.Path == "/api/v1/teams/main/pipelines/pipeline-1/jobs" ||
		r.URL.Path == "/api/v1/teams/awesome/pipelines/pipeline-1/jobs" {

		w.Write([]byte(getJobsResponse))
		return
	}
}

const (
	getPipelinesResponse = `[
  {
    "name": "pipeline-1",
    "url": "/teams/main/pipelines/pipeline-1",
    "paused": false,
    "public": true,
    "team_name": "main"
  }
]`

	getJobsResponse = `[
  {
    "name": "hello-world",
    "url": "/teams/main/pipelines/pipeline-1/jobs/hello-world",
    "next_build": null,
	"paused": true,
    "finished_build": {
      "id": 4,
      "team_name": "main",
      "name": "4",
      "status": "failed",
      "job_name": "hello-world",
      "url": "/teams/main/pipelines/pipeline-1/jobs/hello-world/builds/4",
      "api_url": "/api/v1/builds/4",
      "pipeline_name": "pipeline-1",
      "start_time": 1488167005,
      "end_time": 1488167007
    },
    "inputs": [],
    "outputs": [],
    "groups": []
  }
]`
)
