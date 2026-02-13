package rewardsrobot

import (
	"context"
	"fmt"
	"rewardsAutomation/internal/assets"
	"rewardsAutomation/internal/config"
	"rewardsAutomation/internal/edge"
	"rewardsAutomation/internal/matcher"
	"strings"
	"time"
	"unicode"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-vgo/robotgo"
)

type RewardsRobot struct {
	ctx context.Context
}

func New(ctx context.Context) *RewardsRobot {
	return &RewardsRobot{
		ctx: ctx,
	}
}

func (r *RewardsRobot) Run() (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	err = edge.Kill()
	if err != nil {
		return err
	}

	u := launcher.New().
		Bin(cfg.EdgePath).
		Headless(false).
		UserDataDir(cfg.UserEdgeDir).
		Leakless(false).
		MustLaunch()

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	browser := rod.New().
		ControlURL(u).
		NoDefaultDevice().
		MustConnect()
	defer browser.MustClose()

	newsBingUrl := `https://www.bing.com/news/search?q=Fatos+Principais&nvaug=%5bNewsVertical+Category%3d%22rt_MaxClass%22%5d&FORM=Z9LH3`

	newsPage := browser.MustPage(newsBingUrl).MustWindowMaximize().MustWaitLoad()

	snippetsJS := `() => {
		let snippets = document.querySelectorAll(".snippet")
		let title = ""

		snippets.forEach((snippet) => {
			if (snippet.title.length > title.length) {
				title = snippet.title
			}
		})

		return title
	}`

	snippetTitle := newsPage.MustEval(snippetsJS).String()

	snippetTitle = keepAlphaNumeric(snippetTitle)

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	robotgo.KeySleep = 300

	robotgo.KeyTap(robotgo.KeyT, robotgo.Ctrl)

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	for _, ch := range snippetTitle {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			robotgo.Type(string(ch), 0, cfg.TypeTick)
		}
	}

	err = r.sleepOrCancel(time.Millisecond * 500)
	if err != nil {
		return err
	}
	robotgo.KeyTap(robotgo.Enter)

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	for range cfg.QtdSearches - 1 {
		snippetTitleLength := len(snippetTitle)
		snippetTitle = snippetTitle[:snippetTitleLength-1]

		err = r.clickSearchBar(cfg.LowSpeed, cfg.HighSpeed)
		if err != nil {
			return err
		}

		robotgo.KeyTap(robotgo.End, robotgo.Ctrl)

		err = r.sleepOrCancel(time.Second)
		if err != nil {
			return err
		}

		robotgo.KeyTap(robotgo.Backspace)

		err = r.sleepOrCancel(time.Second)
		if err != nil {
			return err
		}

		if snippetTitle[snippetTitleLength-2] == ' ' {
			snippetTitle = snippetTitle[:snippetTitleLength-1]

			robotgo.KeyTap(robotgo.Backspace)

			err = r.sleepOrCancel(time.Millisecond * 500)
			if err != nil {
				return err
			}
		}

		robotgo.KeyTap(robotgo.Enter)

		err = r.sleepOrCancel(time.Minute)
		if err != nil {
			return err
		}
	}

	bingUrl := "bing.com"

	robotgo.KeyTap(robotgo.KeyT, robotgo.Ctrl)

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	for _, ch := range bingUrl {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			robotgo.Type(string(ch), 0, cfg.TypeTick)
		}
	}

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	robotgo.KeyTap(robotgo.Delete)

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	robotgo.KeyTap(robotgo.Enter)

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	robotgo.KeyTap(robotgo.KeyS, robotgo.Alt, robotgo.Shift)

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	robotgo.MouseDown()
	axisX, axisY := robotgo.Location()

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	robotgo.MoveSmooth(axisX+400, axisY+400, cfg.LowSpeed, cfg.HighSpeed)

	err = r.sleepOrCancel(time.Second)
	if err != nil {
		return err
	}

	robotgo.MouseUp()

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	return nil
}

func (r *RewardsRobot) clickSearchBar(lowSpeed, highSpeed float64) error {
	searchBarPoint, err := matcher.MatchTemplates(assets.SearchBarDark, assets.SearchBarLight)
	if err != nil {
		return err
	}

	err = r.sleepOrCancel(500 * time.Millisecond)
	if err != nil {
		return err
	}

	robotgo.MoveSmooth(searchBarPoint.X+60, searchBarPoint.Y, lowSpeed, highSpeed)

	err = r.sleepOrCancel(500 * time.Millisecond)
	if err != nil {
		return err
	}
	robotgo.Click()

	err = r.sleepOrCancel(time.Second * 2)
	if err != nil {
		return err
	}

	return nil
}

func (r *RewardsRobot) sleepOrCancel(d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-r.ctx.Done():
		return r.ctx.Err()
	case <-t.C:
		return nil
	}
}

func keepAlphaNumeric(s string) string {
	var newString strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) || unicode.IsLetter(r) || unicode.IsSpace(r) {
			newString.WriteRune(r)
		}
	}

	return newString.String()
}
