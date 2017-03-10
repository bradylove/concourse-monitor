package concourse_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/apoydence/onpar"
	"github.com/bradylove/concourse-monitor/lib/concourse"

	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

func TestConcourseClient(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("GetPipelines()", func() {
		o.BeforeEach(func(t *testing.T) (*testing.T, *httptest.Server, chan *http.Request) {
			requests := make(chan *http.Request, 100)
			httpResponseSuccess := func(w http.ResponseWriter, r *http.Request) {
				requests <- r

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(getPipelinesRequest))
			}
			s := httptest.NewServer(http.HandlerFunc(httpResponseSuccess))

			return t, s, requests
		})

		o.AfterEach(func(t *testing.T, s *httptest.Server, r chan *http.Request) {
			s.Close()
		})

		o.Spec("it returns a list of pipelines", func(t *testing.T, s *httptest.Server, r chan *http.Request) {
			client, err := concourse.NewClient(s.URL, "main")
			Expect(t, err).To(Not(HaveOccurred()))

			pipes, err := client.Pipelines()
			Expect(t, err).To(Not(HaveOccurred()))

			req := <-r
			Expect(t, req.URL.Path).To(Equal("/api/v1/teams/main/pipelines"))

			Expect(t, pipes).To(HaveLen(2))
		})
	})
}

var ()

const getPipelinesRequest = `[
  {
    "name": "pipeline-1",
    "url": "/teams/main/pipelines/pipeline-1",
    "paused": false,
    "public": true,
    "team_name": "main"
  },
  {
    "name": "pipeline-1",
    "url": "/teams/main/pipelines/pipeline-1",
    "paused": false,
    "public": true,
    "team_name": "main"
  }
]`
