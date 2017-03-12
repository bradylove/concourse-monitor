package state

import (
	"image"
	"image/color"
)

var (
	ColorPending   = color.RGBA{189, 195, 199, 255}
	ColorStarted   = color.RGBA{241, 196, 15, 255}
	ColorSucceeded = color.RGBA{46, 204, 113, 255}
	ColorFailed    = color.RGBA{231, 76, 60, 255}
	ColorErrored   = color.RGBA{230, 126, 34, 255}
	ColorAborted   = color.RGBA{143, 75, 45, 255}
	ColorPaused    = color.RGBA{52, 152, 219, 255}
	ColorUnknown   = color.RGBA{96, 81, 163, 255}
)

func StateIcon(status string) image.Image {
	return image.NewUniform(statusColor(status))
}

func statusColor(status string) color.Color {
	switch status {
	case StatusPending:
		return ColorPending
	case StatusStarted:
		return ColorStarted
	case StatusSucceeded:
		return ColorSucceeded
	case StatusFailed:
		return ColorFailed
	case StatusErrored:
		return ColorErrored
	case StatusAborted:
		return ColorAborted
	case StatusPaused:
		return ColorPaused
	default:
		return ColorUnknown
	}
}
