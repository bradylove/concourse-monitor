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
		o.Group("with successful response", func() {
			o.BeforeEach(func(t *testing.T) (*testing.T, *httptest.Server, chan *http.Request) {
				requests := make(chan *http.Request, 100)
				handler := func(w http.ResponseWriter, r *http.Request) {
					requests <- r

					w.WriteHeader(http.StatusOK)
					w.Write([]byte(getPipelinesRequest))
				}
				s := httptest.NewServer(http.HandlerFunc(handler))

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
				client, err := concourse.NewClient(s.URL, "main")
				Expect(t, err).To(Not(HaveOccurred()))

				_, err = client.Pipelines()
				Expect(t, err.Error()).To(Equal("expected 200 response code, got 500"))
			})
		})

		o.Group("with a request error", func() {
			o.Spec("it returns an error", func(t *testing.T) {
				client, err := concourse.NewClient("http://127.0.0.1:23223", "main")
				Expect(t, err).To(Not(HaveOccurred()))

				_, err = client.Pipelines()
				Expect(t, err).To(HaveOccurred())
			})
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
