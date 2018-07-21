package state_test

import (
	"image/color"
	"testing"

	"github.com/apoydence/onpar"
	"github.com/bradylove/concourse-monitor/pkg/state"

	. "github.com/apoydence/onpar/expect"
	. "github.com/apoydence/onpar/matchers"
)

func TestIcons(t *testing.T) {
	o := onpar.New()
	defer o.Run(t)

	o.Group("StateIcon()", func() {
		o.Spec("it returns the icon matching the state", func(t *testing.T) {
			testset := []struct {
				status string
				color  color.Color
			}{
				{state.StatusPending, state.ColorPending},
				{state.StatusStarted, state.ColorStarted},
				{state.StatusSucceeded, state.ColorSucceeded},
				{state.StatusFailed, state.ColorFailed},
				{state.StatusErrored, state.ColorErrored},
				{state.StatusAborted, state.ColorAborted},
				{state.StatusPaused, state.ColorPaused},
				{state.StatusUnknown, state.ColorUnknown},
			}

			for _, ts := range testset {
				_ = ts
				Expect(t, state.StatusIcon(ts.status).At(1, 1)).To(Equal(ts.color))
			}
		})
	})
}
