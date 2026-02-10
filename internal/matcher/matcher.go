package matcher

import (
	"fmt"
	"image"

	"github.com/go-vgo/robotgo"
	"gocv.io/x/gocv"
)

// MatchTemplate returns the location of the template in the screenshot
func MatchTemplate(template []byte) (image.Point, error) {
	screenshot := captureScreen()

	screenshotMat, err := gocv.ImageToMatRGB(screenshot)
	if err != nil {
		return image.Point{}, err
	}
	defer screenshotMat.Close()

	templateMat, err := gocv.IMDecode(template, gocv.IMReadColor)
	if err != nil {
		return image.Point{}, err
	}
	defer templateMat.Close()

	result := gocv.NewMat()
	defer result.Close()

	err = gocv.MatchTemplate(screenshotMat, templateMat, &result, gocv.TmCcoeffNormed, gocv.NewMat())
	if err != nil {
		return image.Point{}, err
	}

	_, maxValue, _, maxLoc := gocv.MinMaxLoc(result)

	if maxValue <= 0.85 {
		return image.Point{}, fmt.Errorf("template not found")
	}

	return maxLoc, nil
}

// MatchTemplates returns the location of the templates in the screenshot
func MatchTemplates(templates ...[]byte) (image.Point, error) {
	for _, template := range templates {
		maxLoc, err := MatchTemplate(template)
		if err != nil {
			continue
		}

		return maxLoc, nil
	}

	return image.Point{}, fmt.Errorf("templates not found")
}

// FindTemplate returns true if the template is found in the screenshot
func FindTemplate(template []byte) (bool, error) {
	screenshot := captureScreen()

	screenshotMat, err := gocv.ImageToMatRGB(screenshot)
	if err != nil {
		return false, err
	}
	defer screenshotMat.Close()

	templateMat, err := gocv.IMDecode(template, gocv.IMReadColor)
	if err != nil {
		return false, err
	}
	defer templateMat.Close()

	result := gocv.NewMat()
	defer result.Close()

	err = gocv.MatchTemplate(screenshotMat, templateMat, &result, gocv.TmCcoeffNormed, gocv.NewMat())
	if err != nil {
		return false, err
	}

	_, maxValue, _, _ := gocv.MinMaxLoc(result)

	if maxValue <= 0.85 {
		return false, nil
	}

	return true, nil
}

// FindTemplates returns true if the templates are found in the screenshot
func FindTemplates(templates ...[]byte) (bool, error) {
	for _, template := range templates {
		isFind, err := FindTemplate(template)
		if err != nil {
			continue
		}

		return isFind, nil
	}

	return false, fmt.Errorf("%d templates not found", len(templates))
}

func captureScreen() image.Image {
	bitmapScreenshot := robotgo.CaptureScreen()
	defer robotgo.FreeBitmap(bitmapScreenshot)

	screenshotIMG := robotgo.ToImage(bitmapScreenshot)

	return screenshotIMG
}
