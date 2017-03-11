package concourse_test

import (
	"io/ioutil"
	"testing"

	"github.com/apoydence/onpar"
	"github.com/bradylove/concourse-monitor/lib/concourse"

	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

func TestTargets(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("LoadTargets()", func() {
		o.Group("with a valid .flyrc", func() {
			o.BeforeEach(func(t *testing.T) (*testing.T, string) {
				f, err := ioutil.TempFile("", "flyrc")
				Expect(t, err).To(Not(HaveOccurred()))
				defer f.Close()

				_, err = f.Write([]byte(flyRCWithTargets))
				Expect(t, err).To(Not(HaveOccurred()))

				return t, f.Name()
			})

			o.Spec("it returns a list of targets", func(t *testing.T, filename string) {
				targets, err := concourse.LoadTargets(filename)
				Expect(t, err).To(Not(HaveOccurred()))
				Expect(t, targets).To(HaveLen(2))

				t1 := targets[0]
				Expect(t, t1.API).To(Equal("http://127.0.0.1:8080"))
				Expect(t, t1.Team).To(Equal("awesome"))
				Expect(t, t1.Token.Type).To(Equal("Bearer"))
				Expect(t, t1.Token.Value).To(Equal("a-token"))
			})
		})

		o.Group("with an invalid file path", func() {
			o.Spec("it returns an error", func(t *testing.T) {
				_, err := concourse.LoadTargets("unknown-file")
				Expect(t, err).To(HaveOccurred())
			})
		})
	})
}

const flyRCWithTargets = `targets:
  targetone:
    api: http://127.0.0.1:8080
    team: awesome
    token:
      type: Bearer
      value: a-token
  targettwo:
    api: https://example-domain.tld
    team: win
    token:
      type: Bearer
      value: a-token
`
