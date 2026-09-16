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
	"github.com/go-vgo/robotgo/clipboard"
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
			slog.Error("panic recovered", "error", r)
			beeep.Notify(beeep.AppName, "Ocorreu erro na execução.", assets.RewardsLogoPNG.Data)
			beeep.Notify(beeep.AppName, fmt.Sprintf("%v", r), assets.RewardsLogoPNG.Data)
			slog.Info("Time elapsed", "time", time.Since(startTime))
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

	err = edge.CopyUserData(cfg.UserDataDir, cfg.TmpUserDataDir)
	if err != nil {
		return err
	}

	defer edge.RemoveTempUserData(cfg.TmpUserDataDir)

	l := launcher.New().
		Bin(cfg.EdgePath).
		Headless(false).
		UserDataDir(cfg.TmpUserDataDir).
		NoSandbox(true).
		Leakless(false)
	u := l.MustLaunch()

	time.Sleep(time.Second)

	browser := rod.New().
		ControlURL(u).
		NoDefaultDevice().
		MustConnect()
	defer browser.MustClose()

	time.Sleep(time.Second * 5)

	pages := browser.MustPages()
	var page *rod.Page

	if len(pages) > 0 {
		page = pages[0]
	} else {
		page = browser.MustPage("edge://newtab")
	}
	defer page.MustClose()

	time.Sleep(time.Millisecond * 500)

	// Force the window to the OS foreground by PID — needed when launched by Task Scheduler.
	edge.FocusPID(uint32(l.PID()))

	time.Sleep(time.Second)

	page.MustWindowMaximize().MustActivate()

	bingUrl := "https://bing.com"

	page.MustNavigate(bingUrl).MustWaitLoad()

	time.Sleep(time.Second * 2)

	robotgo.MoveSmooth(0, 0, cfg.LowSpeed, cfg.HighSpeed)
	time.Sleep(time.Millisecond * 500)
	robotgo.Click()
	time.Sleep(time.Millisecond * 500)

	agreeContinueTpl, err := matcher.NewTemplate(assets.AgreeContinue.Name, assets.AgreeContinue.Data)
	if err != nil {
		return err
	}
	defer agreeContinueTpl.Close()

	err = matcher.MatchTemplateAndClickCenter(r.ctx, agreeContinueTpl, time.Second*5)
	if err != nil {
		slog.Error(assets.AgreeContinue.Name+" button not found", "error", err)
	}

	visualSearchDarkTpl, err := matcher.NewTemplate(assets.VisualSearchDark.Name, assets.VisualSearchDark.Data)
	if err != nil {
		return err
	}
	defer visualSearchDarkTpl.Close()

	visualSearchWhiteTpl, err := matcher.NewTemplate(assets.VisualSearchWhite.Name, assets.VisualSearchWhite.Data)
	if err != nil {
		return err
	}
	defer visualSearchWhiteTpl.Close()

	visualSearchPos, err := matcher.MatchTemplatesCenterWithTimeout(r.ctx, time.Second*5, visualSearchWhiteTpl, visualSearchDarkTpl)
	if err != nil {
		return err
	}

	robotgo.MoveSmooth(
		visualSearchPos.X,
		visualSearchPos.Y,
		cfg.LowSpeed,
		cfg.HighSpeed,
	)

	time.Sleep(time.Millisecond * 500)

	robotgo.Click()

	time.Sleep(time.Millisecond * 500)

	imageURLTpl, err := matcher.NewTemplate(assets.PasteImageURL.Name, assets.PasteImageURL.Data)
	if err != nil {
		return err
	}
	defer imageURLTpl.Close()

	messiImageURL := "https://upload.wikimedia.org/wikipedia/commons/thumb/c/c8/Leo_Messi_Argentina_v_Egypt_7_July_2026-1.jpg/250px-Leo_Messi_Argentina_v_Egypt_7_July_2026-1.jpg"

	err = clipboard.WriteAll(messiImageURL)
	if err != nil {
		return err
	}

	err = matcher.MatchTemplateAndClickCenter(r.ctx, imageURLTpl, time.Second*5)
	if err != nil {
		return err
	}

	time.Sleep(time.Millisecond * 500)

	robotgo.KeyTap(robotgo.KeyV, robotgo.Ctrl)

	time.Sleep(time.Second * 5)

	err = matcher.MatchTemplateAndClickCenter(r.ctx, agreeContinueTpl, time.Second*5)
	if err != nil {
		slog.Error(assets.AgreeContinue.Name+" button not found", "error", err)
	}

	r.sleepOrCancel(time.Minute)

	robotgo.KeyTap(robotgo.F4)

	time.Sleep(time.Millisecond * 500)

	newsBingUrl := `https://www.bing.com/news/search?q=Fatos+Principais&nvaug=%5bNewsVertical+Category%3d%22rt_MaxClass%22%5d&FORM=Z9LH3`

	for _, c := range newsBingUrl {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			robotgo.Type(string(c), 0, cfg.TypeTick)
		}
	}

	time.Sleep(time.Millisecond * 500)

	robotgo.KeyTap(robotgo.Enter)

	snippetsJS := fmt.Sprintf(`() => {
		let snippets = [...document.querySelectorAll(".news_title")]

		return snippets
			.map(s => s.title)
			.filter(Boolean)
			.sort((a, b) => b.length - a.length)
			.slice(0,%d)
	}`, cfg.QtdSearches+1)

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	snippetTitles := page.MustEval(snippetsJS).Arr()

	titles := make([]string, 0, len(snippetTitles))

	for i, title := range snippetTitles {
		titles = append(titles, title.String())
		titles[i] = keepAlphaNumeric(titles[i])
	}

	time.Sleep(time.Second)

	robotgo.KeySleep = 300

	robotgo.KeyTap(robotgo.KeyT, robotgo.Ctrl)

	time.Sleep(time.Second)

	for _, ch := range titles[0] {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			robotgo.Type(string(ch), 0, cfg.TypeTick)
		}
	}

	titles = titles[1:]

	time.Sleep(time.Millisecond * 500)

	robotgo.KeyTap(robotgo.Enter)

	time.Sleep(time.Second)

	pages = browser.MustPages()
	searchPage := pages[1]
	defer searchPage.MustClose()

	searchPage.MustWaitLoad()

	err = matcher.MatchTemplateAndClickCenter(r.ctx, agreeContinueTpl, time.Second*5)
	if err != nil {
		slog.Error(assets.AgreeContinue.Name+" button not found", "error", err)
	}

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	for i := range cfg.QtdSearches - 1 {
		err = r.clickSearchBar()
		if err != nil {
			return err
		}

		time.Sleep(time.Millisecond * 500)

		robotgo.KeyTap(robotgo.KeyA, robotgo.Ctrl)

		time.Sleep(time.Second)

		for _, ch := range titles[i] {
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
	}

	robotgo.KeyTap(robotgo.F4)

	time.Sleep(time.Millisecond * 500)

	rewardsUrl := "https://rewards.bing.com/dashboard"

	for _, c := range rewardsUrl {
		select {
		case <-r.ctx.Done():
			return r.ctx.Err()
		default:
			robotgo.Type(string(c), 0, cfg.TypeTick)
		}
	}

	time.Sleep(time.Millisecond * 500)

	robotgo.KeyTap(robotgo.Enter)

	time.Sleep(time.Millisecond * 500)

	searchPage.MustWaitLoad()

	time.Sleep(time.Second * 10)

	readyToClaimTpl, err := matcher.NewTemplate(assets.ReadyToClaim.Name, assets.ReadyToClaim.Data)
	if err != nil {
		return err
	}
	defer readyToClaimTpl.Close()

	err = matcher.MatchTemplateAndClickCenter(r.ctx, readyToClaimTpl, time.Second*5)
	if err != nil {
		return fmt.Errorf(assets.AgreeContinue.Name+" not found: %v", err)
	}

	time.Sleep(time.Second * 5)

	claimPointsTpl, err := matcher.NewTemplate(readyToClaimTpl.Name, assets.ClaimPoints.Data)
	if err != nil {
		return err
	}
	defer claimPointsTpl.Close()

	err = matcher.MatchTemplateAndClickCenter(r.ctx, claimPointsTpl, time.Second*5)
	if err != nil {
		slog.Error(assets.ClaimPoints.Name+" not found", "error", err)

		closeClaimPointsTpl, err := matcher.NewTemplate(assets.CloseClaimPoints.Name, assets.CloseClaimPoints.Data)
		if err != nil {
			return err
		}
		defer closeClaimPointsTpl.Close()

		err = matcher.MatchTemplateAndClickCenter(r.ctx, closeClaimPointsTpl, time.Second*5)
		if err != nil {
			return fmt.Errorf("%s not found: %v", assets.ClaimPoints.Name, err)
		}
	}

	time.Sleep(time.Second * 5)

	anchors := searchPage.MustElements("#dailyset div.grid > a")

	for _, a := range anchors {
		a.MustEval(`() => {this.click()}`)

		err = r.sleepOrCancel(time.Minute)
		if err != nil {
			return err
		}

		robotgo.KeyTap(robotgo.KeyW, robotgo.Ctrl)

		time.Sleep(time.Second * 2)
	}

	ganharLabelTpl, err := matcher.NewTemplate(assets.GanharLabel.Name, assets.GanharLabel.Data)
	if err != nil {
		return err
	}
	defer ganharLabelTpl.Close()

	err = matcher.MatchTemplateAndClickCenter(r.ctx, ganharLabelTpl, time.Second*5)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 10)

	anchors = searchPage.MustElementsByJS(`
		() => [...document.querySelectorAll("#moreactivities a")].filter(a => 
			a.textContent.includes("+")
		)
	`)

	for _, a := range anchors {
		a.MustEval(`() => {this.click()}`)

		err = r.sleepOrCancel(time.Minute)
		if err != nil {
			return err
		}

		robotgo.KeyTap(robotgo.KeyW, robotgo.Ctrl)

		time.Sleep(time.Second * 2)
	}

	return nil
}

func (r *RewardsRobot) clickSearchBar() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	searchBarDarkTpl, err := matcher.NewTemplate(assets.SearchBarDark.Name, assets.SearchBarDark.Data)
	if err != nil {
		return err
	}
	defer searchBarDarkTpl.Close()

	searchBarLightTpl, err := matcher.NewTemplate(assets.SearchBarLight.Name, assets.SearchBarLight.Data)
	if err != nil {
		return err
	}
	defer searchBarLightTpl.Close()

	searchPlusLightTpl, err := matcher.NewTemplate(assets.SearchPlusLight.Name, assets.SearchPlusLight.Data)
	if err != nil {
		return err
	}
	defer searchPlusLightTpl.Close()

	searchPlusDarkTpl, err := matcher.NewTemplate(assets.SearchPlusDark.Name, assets.SearchPlusDark.Data)
	if err != nil {
		return err
	}
	defer searchPlusDarkTpl.Close()

	searchBarPoint, err := matcher.MatchTemplatesCenterWithTimeout(r.ctx, time.Second*20, searchBarDarkTpl, searchBarLightTpl, searchPlusLightTpl, searchBarDarkTpl)
	if err != nil {
		slog.Error("template not found", "error", err)
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
