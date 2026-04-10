package tray

import (
	"context"
	"rewardsAutomation/internal/assets"

	"github.com/gen2brain/beeep"
	"github.com/getlantern/systray"
)

func onReady(cancel context.CancelFunc, errChan <-chan error, doneChan <-chan struct{}) func() {
	return func() {
		systray.SetIcon(assets.RewardsLogoICO.Data)
		systray.SetTitle("Rewards Robot")
		systray.SetTooltip("Automação Bing Rewards")

		mQuit := systray.AddMenuItem("Fechar", "Finaliza o robô")

		go func() {
			<-mQuit.ClickedCh
			cancel()
		}()

		go func() {
			select {
			case <-errChan:
				beeep.Notify(beeep.AppName, "Ocorreu um erro durante a execução.", assets.RewardsLogoPNG.Data)
				systray.Quit()
			case <-doneChan:
				beeep.Notify(beeep.AppName, "Execução finalizada.", assets.RewardsLogoPNG.Data)
				systray.Quit()
			}
		}()
	}
}
