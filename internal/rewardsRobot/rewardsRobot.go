package rewardsrobot

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"rewardsAutomation/internal/assets"
	"rewardsAutomation/internal/config"
	"rewardsAutomation/internal/edge"
	"rewardsAutomation/internal/matcher"
	"strings"
	"time"
	"unicode"

	"github.com/gen2brain/beeep"
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
	startTime := time.Now()

	beeep.Notify(beeep.AppName, "Iniciando execução...", assets.RewardsLogoPNG.Data)

	defer func() {
		r := recover()
		if r != nil {
			beeep.Notify(beeep.AppName, fmt.Sprintf("Ocorreu erro na execução.\n%v", err), assets.RewardsLogoPNG.Data)
			return
		}
		slog.Info("Time elapsed", "time", time.Since(startTime))
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

	time.Sleep(time.Second)

	browser := rod.New().
		ControlURL(u).
		NoDefaultDevice().
		MustConnect()
	defer browser.MustClose()

	newsBingUrl := `https://www.bing.com/news/search?q=Fatos+Principais&nvaug=%5bNewsVertical+Category%3d%22rt_MaxClass%22%5d&FORM=Z9LH3`

	newsPage := browser.MustPage(newsBingUrl).MustWindowMaximize().MustWaitLoad()

	err = matcher.MatchTemplateAndClickCenter(r.ctx, assets.AceitarButton.Data, time.Second*20)
	if err != nil {
		slog.Error(assets.AceitarButton.Name+" button not found", "error", err)
	}

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

	time.Sleep(time.Second)

	robotgo.KeySleep = 300

	robotgo.KeyTap(robotgo.KeyT, robotgo.Ctrl)

	time.Sleep(time.Second)

	for _, ch := range snippetTitle {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			robotgo.Type(string(ch), 0, cfg.TypeTick)
		}
	}

	time.Sleep(time.Millisecond * 500)

	robotgo.KeyTap(robotgo.Enter)

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	for range cfg.QtdSearches - 1 {
		snippetTitleLength := len(snippetTitle)
		snippetTitle = snippetTitle[:snippetTitleLength-1]

		err = r.clickSearchBar()
		if err != nil {
			return err
		}

		robotgo.KeyTap(robotgo.End, robotgo.Ctrl)

		time.Sleep(time.Second)

		robotgo.KeyTap(robotgo.Backspace)

		time.Sleep(time.Second)

		if snippetTitle[snippetTitleLength-2] == ' ' {
			snippetTitle = snippetTitle[:snippetTitleLength-1]

			robotgo.KeyTap(robotgo.Backspace)

			time.Sleep(time.Millisecond * 500)
		}

		robotgo.KeyTap(robotgo.Enter)

		err = r.sleepOrCancel(time.Minute)
		if err != nil {
			return err
		}
	}

	rewardsUrl := "https://rewards.bing.com/"
	rewardsPage := browser.MustPage(rewardsUrl).MustWindowMaximize().MustWaitLoad()

	getCardsLength := `() => {
		const divs = document.querySelectorAll("span[mee-heading='heading']")
		return divs.length
	}`

	clickCard := `idx => {
		const divs = document.querySelectorAll("span[mee-heading='heading']")
		const parent = divs[idx]?.closest(".ds-card-sec")

		parent.querySelector("div.actionLink.x-hidden-vp1")?.children?.[0]?.click()

		return true
	}`

	time.Sleep(time.Second * 3)

	cardsLength := rewardsPage.MustEval(getCardsLength).Int()

	for idx := range cardsLength {
		rewardsPage.MustEval(clickCard, idx)

		err = r.sleepOrCancel(time.Second * 20)
		if err != nil {
			return err
		}

		rewardsPage.MustActivate().MustWaitLoad()

		time.Sleep(time.Second * 5)
	}

	return nil
}

func (r *RewardsRobot) clickSearchBar() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	searchBarPoint, err := matcher.MatchTemplatesWithTimeout(r.ctx, time.Second*20, assets.SearchBarDark.Data, assets.SearchBarLight.Data)
	if err != nil {
		return err
	}

	time.Sleep(500 * time.Millisecond)

	robotgo.MoveSmooth(searchBarPoint.X+60, searchBarPoint.Y, cfg.LowSpeed, cfg.HighSpeed)

	time.Sleep(500 * time.Millisecond)

	robotgo.Click()

	return nil
}

func (r *RewardsRobot) sleepOrCancel(d time.Duration) error {
	random := rand.IntN(int(time.Second * 30))
	t := time.NewTimer(d + time.Duration(random))
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
