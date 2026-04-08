package matcher

import (
	"context"
	"fmt"
	"image"
	"time"

	"rewardsAutomation/internal/config"

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

func MatchTemplateWithTimeout(ctx context.Context, template []byte, timeout time.Duration) (image.Point, error) {
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return image.Point{}, ctx.Err()
		default:
			maxLoc, err := MatchTemplate(template)
			if err == nil {
				return maxLoc, nil
			}

			if time.Since(start) > timeout {
				return image.Point{}, fmt.Errorf("template not found within timeout: %v", timeout)
			}
		}
	}
}

func MatchTemplatesWithTimeout(ctx context.Context, timeout time.Duration, templates ...[]byte) (image.Point, error) {
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return image.Point{}, ctx.Err()
		default:
			maxLoc, err := MatchTemplates(templates...)
			if err == nil {
				return maxLoc, nil
			}

			if time.Since(start) > timeout {
				return image.Point{}, fmt.Errorf("templates not found within timeout: %v", timeout)
			}
		}
	}
}

func MatchTemplateAndClick(ctx context.Context, template []byte, timeout time.Duration) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	maxLoc, err := MatchTemplateWithTimeout(ctx, template, timeout)
	if err != nil {
		return err
	}

	robotgo.MoveSmooth(maxLoc.X, maxLoc.Y, cfg.LowSpeed, cfg.HighSpeed)
	robotgo.Click()

	return nil
}

func MatchTemplateAndClickCenter(ctx context.Context, template []byte, timeout time.Duration) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	maxLoc, err := MatchTemplateWithTimeout(ctx, template, timeout)
	if err != nil {
		return err
	}

	// Decode template to get its size
	templateMat, err := gocv.IMDecode(template, gocv.IMReadColor)
	if err != nil {
		return err
	}
	defer templateMat.Close()

	templateWidth := templateMat.Cols()
	templateHeight := templateMat.Rows()

	centerX := maxLoc.X + templateWidth/2
	centerY := maxLoc.Y + templateHeight/2

	robotgo.MoveSmooth(centerX, centerY, cfg.LowSpeed, cfg.HighSpeed)
	robotgo.Click()

	return nil
}
