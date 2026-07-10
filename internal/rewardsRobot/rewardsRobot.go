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
	"github.com/go-vgo/robotgo"
	"github.com/mxschmitt/playwright-go"
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
			beeep.Notify(beeep.AppName, fmt.Sprintf("Ocorreu erro na execução.\n%v", r), assets.RewardsLogoPNG.Data)
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

	err = edge.CopyUserData(cfg.UserEdgeDir, cfg.TmpUserDataDir)
	if err != nil {
		return err
	}

	defer edge.RemoveTempUserData(cfg.TmpUserDataDir)

	// l := launcher.New().
	// 	Bin(cfg.EdgePath).
	// 	Headless(false).
	// 	// UserDataDir(cfg.UserEdgeDir).
	// 	Leakless(false)

	// u := l.MustLaunch()

	// time.Sleep(time.Second)

	// browser := rod.New().
	// 	ControlURL(u).
	// 	NoDefaultDevice().
	// 	MustConnect()
	// defer browser.MustClose()

	pw, err := playwright.Run()
	if err != nil {
		return err
	}
	defer pw.Stop()

	contextBrowser, err := pw.Chromium.LaunchPersistentContext(cfg.TmpUserDataDir, playwright.BrowserTypeLaunchPersistentContextOptions{
		Channel:  new("msedge"),
		Headless: new(false),
		IgnoreDefaultArgs: []string{
			"--no-sandbox",
			"about:blank",
			// "--disable-field-trial-config",
			// "--disable-background-networking",
			// "--disable-background-timer-throttling",
			// "--disable-backgrounding-occluded-windows",
			// "--disable-back-forward-cache",
			// "--disable-breakpad",
			// "--disable-client-side-phishing-detection",
			// "--disable-component-extensions-with-background-pages",
			// "--disable-component-update",
			// "--no-default-browser-check",
			// "--disable-default-apps",
			// "--disable-dev-shm-usage",
			// "--disable-edgeupdater",
			// "--disable-extensions",
			// "--disable-features=AvoidUnnecessaryBeforeUnloadCheckSync,BoundaryEventDispatchTracksNodeRemoval,DestroyProfileOnBrowserClose,DialMediaRouteProvider,GlobalMediaControls,HttpsUpgrades,LensOverlay,MediaRouter,PaintHolding,ThirdPartyStoragePartitioning,Translate,AutoDeElevate,RenderDocument,OptimizationHints,msForceBrowserSignIn,msEdgeUpdateLaunchServicesPreferredVersion",
			// "--enable-features=CDPScreenshotNewSurface",
			// "--allow-pre-commit-input",
			// "--disable-hang-monitor",
			// "--disable-ipc-flooding-protection",
			// "--disable-popup-blocking",
			// "--disable-prompt-on-repost",
			// "--disable-renderer-backgrounding",
			// "--force-color-profile=srgb",
			// "--metrics-recording-only",
			// "--no-first-run",
			// "--password-store=basic",
			// "--use-mock-keychain",
			// "--no-service-autorun",
			// "--export-tagged-pdf",
			// "--disable-search-engine-choice-screen",
			// "--unsafely-disable-devtools-self-xss-warnings",
			// "--edge-skip-compat-layer-relaunch",
			// "--disable-infobars",
			// "--disable-search-engine-choice-screen",
			// "--disable-sync",
			// "--enable-unsafe-swiftshader",
			// "--remote-debugging-pipe",
		},
		Args: []string{
			"--start-maximized",
		},
		NoViewport: new(true),
		Locale:     new("pt-BR"),
	})
	if err != nil {
		return err
	}
	defer contextBrowser.Close()

	time.Sleep(time.Second * 5)

	pages := contextBrowser.Pages()
	var page playwright.Page

	if len(pages) > 0 {
		page = pages[0]
	} else {
		page, err = contextBrowser.NewPage()
		if err != nil {
			return err
		}
	}

	newsBingUrl := `https://www.bing.com/news/search?q=Fatos+Principais&nvaug=%5bNewsVertical+Category%3d%22rt_MaxClass%22%5d&FORM=Z9LH3`

	// newsPage := browser.MustPage(newsBingUrl).MustWaitLoad()
	_, err = page.Goto(newsBingUrl)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 5)

	// // Force the window to the OS foreground by PID — needed when launched by Task Scheduler.
	// edge.FocusPID(uint32(l.PID()))
	// time.Sleep(time.Millisecond * 500)

	// newsPage.MustWindowMaximize().MustActivate()

	// time.Sleep(time.Millisecond * 500)

	// robotgo.MoveSmooth(0, 0, cfg.LowSpeed, cfg.HighSpeed)
	// time.Sleep(time.Millisecond * 500)
	// robotgo.Click()
	// time.Sleep(time.Millisecond * 500)

	err = matcher.MatchTemplateAndClickCenter(r.ctx, assets.AgreeContinue.Data, time.Second*20)
	if err != nil {
		slog.Error(assets.AgreeContinue.Name+" button not found", "error", err)
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

	err = r.sleepOrCancel(time.Minute)
	if err != nil {
		return err
	}

	snippetTitle, err := page.Evaluate(snippetsJS)
	if err != nil {
		return err
	}

	snippetTitleStr := snippetTitle.(string)

	snippetTitleStr = keepAlphaNumeric(snippetTitleStr)

	time.Sleep(time.Second)

	robotgo.KeySleep = 300

	robotgo.KeyTap(robotgo.KeyT, robotgo.Ctrl)

	time.Sleep(time.Second)

	for _, ch := range snippetTitleStr {
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
		snippetTitleLength := len(snippetTitleStr)
		snippetTitle = snippetTitleStr[:snippetTitleLength-1]

		err = r.clickSearchBar()
		if err != nil {
			return err
		}

		robotgo.KeyTap(robotgo.End, robotgo.Ctrl)

		time.Sleep(time.Second)

		robotgo.KeyTap(robotgo.Backspace)

		time.Sleep(time.Second)

		if snippetTitleStr[snippetTitleLength-2] == ' ' {
			snippetTitle = snippetTitleStr[:snippetTitleLength-1]

			robotgo.KeyTap(robotgo.Backspace)

			time.Sleep(time.Millisecond * 500)
		}

		robotgo.KeyTap(robotgo.Enter)

		err = r.sleepOrCancel(time.Minute)
		if err != nil {
			return err
		}
	}

	// rewardsUrl := "https://rewards.bing.com/"
	// rewardsPage := browser.MustPage(rewardsUrl).MustWindowMaximize().MustWaitLoad().MustActivate()

	return nil
}

func (r *RewardsRobot) clickSearchBar() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	searchBarPoint, err := matcher.MatchTemplatesCenterWithTimeout(r.ctx, time.Second*20, assets.SearchBarDark.Data, assets.SearchBarLight.Data)
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
