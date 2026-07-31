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

type Template struct {
	Name   string
	Mat    gocv.Mat
	Width  int
	Height int
}

func NewTemplate(name string, data []byte) (*Template, error) {
	mat, err := gocv.IMDecode(data, gocv.IMReadColor)
	if err != nil {
		return nil, err
	}

	return &Template{
		Name:   name,
		Mat:    mat,
		Width:  mat.Cols(),
		Height: mat.Rows(),
	}, nil
}

func (t *Template) Close() {
	t.Mat.Close()
}

// MatchTemplate returns the location of the template in the screenshot
func MatchTemplate(tpl *Template) (image.Point, error) {
	screenshot := captureScreen()

	screenshotMat, err := gocv.ImageToMatRGB(screenshot)
	if err != nil {
		return image.Point{}, err
	}
	defer screenshotMat.Close()

	return matchTemplateOnMat(screenshotMat, tpl)
}

func matchTemplateOnMat(screen gocv.Mat, tpl *Template) (image.Point, error) {
	result := gocv.NewMat()
	defer result.Close()

	err := gocv.MatchTemplate(
		screen,
		tpl.Mat,
		&result,
		gocv.TmCcoeffNormed,
		gocv.NewMat(),
	)
	if err != nil {
		return image.Point{}, nil
	}

	_, maxValue, _, maxLoc := gocv.MinMaxLoc(result)

	if maxValue <= 0.85 {
		return image.Point{}, fmt.Errorf("template %s not found (%.2f)", tpl.Name, maxValue)
	}

	return maxLoc, nil
}

// MatchTemplates returns the location of the templates in the screenshot
func MatchTemplates(templates ...*Template) (image.Point, error) {
	screenshot := captureScreen()

	screenshotMat, err := gocv.ImageToMatRGB(screenshot)
	if err != nil {
		return image.Point{}, err
	}
	defer screenshotMat.Close()

	for _, tpl := range templates {
		loc, err := matchTemplateOnMat(screenshotMat, tpl)
		if err == nil {
			return loc, nil
		}
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

// MatchTemplateWithTimeout returns the location of the template in the screenshot with a timeout
func MatchTemplateWithTimeout(ctx context.Context, template *Template, timeout time.Duration) (image.Point, error) {
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

func MatchTemplateCenterWithTimeout(ctx context.Context, template *Template, timeout time.Duration) (image.Point, error) {
	start := time.Now()

	for {
		select {
		case <-ctx.Done():
			return image.Point{}, ctx.Err()
		default:
			if time.Since(start) > timeout {
				return image.Point{}, fmt.Errorf("template not found within timeout: %v", timeout)
			}

			maxLoc, err := MatchTemplate(template)
			if err != nil {
				continue
			}

			centerX := maxLoc.X + template.Width/2
			centerY := maxLoc.Y + template.Height/2

			return image.Point{X: centerX, Y: centerY}, nil
		}
	}
}

// MatchTemplatesWithTimeout returns the first location of the templates in the screenshot with a timeout
func MatchTemplatesWithTimeout(ctx context.Context, timeout time.Duration, templates ...*Template) (image.Point, error) {
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

func MatchTemplateAndClick(ctx context.Context, template *Template, timeout time.Duration) error {
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

func MatchTemplatesCenterWithTimeout(ctx context.Context, timeout time.Duration, templates ...*Template) (image.Point, error) {
	start := time.Now()
	idx := 0

	for {
		select {
		case <-ctx.Done():
			return image.Point{}, ctx.Err()
		default:
			if time.Since(start) > timeout {
				return image.Point{}, fmt.Errorf("templates not found within timeout: %v", timeout)
			}

			if idx >= len(templates) {
				idx = 0
			}

			maxLoc, err := MatchTemplate(templates[idx])
			if err != nil {
				idx++
				continue
			}

			centerX := maxLoc.X + templates[idx].Width/2
			centerY := maxLoc.Y + templates[idx].Height/2

			return image.Point{X: centerX, Y: centerY}, nil
		}
	}
}

func MatchTemplateAndClickCenter(ctx context.Context, template *Template, timeout time.Duration) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	maxLoc, err := MatchTemplateWithTimeout(ctx, template, timeout)
	if err != nil {
		return err
	}

	centerX := maxLoc.X + template.Width/2
	centerY := maxLoc.Y + template.Height/2

	robotgo.MoveSmooth(centerX, centerY, cfg.LowSpeed, cfg.HighSpeed)
	robotgo.Click()

	return nil
}
