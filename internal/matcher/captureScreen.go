package matcher

import (
	"image"

	"github.com/go-vgo/robotgo"
)

func captureScreen() image.Image {
	bitmapScreenshot := robotgo.CaptureScreen()
	defer robotgo.FreeBitmap(bitmapScreenshot)

	screenshotIMG := robotgo.ToImage(bitmapScreenshot)

	return screenshotIMG
}
