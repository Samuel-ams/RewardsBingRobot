package tray

import (
	"context"
	"fmt"
	"rewardsAutomation/internal/assets"
	rewardsrobot "rewardsAutomation/internal/rewardsRobot"
	"sync"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

func onReady() {
	systray.SetIcon(assets.RewardsLogoICO)
	systray.SetTitle("Rewards Robot")
	systray.SetTooltip("Automação Bing Rewards")

	var errChan = make(chan error, 1)
	var doneChan = make(chan struct{}, 1)
	var wg sync.WaitGroup

	beeep.Notify(beeep.AppName, "Iniciando execução...", assets.RewardsLogoPNG)

	mQuit := systray.AddMenuItem("Fechar", "Finaliza o robô")

	ctx, cancel := context.WithCancel(context.Background())

	wg.Go(func() {
		robot := rewardsrobot.New(ctx)
		err := robot.Run()
		if err != nil && err != context.Canceled {
			errChan <- err
			cancel()
			return
		}
		cancel()
		doneChan <- struct{}{}
	})

	go func() {
		<-mQuit.ClickedCh
		cancel()
		wg.Wait()
		doneChan <- struct{}{}
	}()

	go func() {
		select {
		case err := <-errChan:
			beeep.Notify(beeep.AppName, fmt.Sprintf("Ocorreu erro na execução.\n%v", err), assets.RewardsLogoPNG)
			systray.Quit()
		case <-doneChan:
			beeep.Notify(beeep.AppName, "Execução finalizada.", assets.RewardsLogoPNG)
			systray.Quit()
		}
	}()
}
