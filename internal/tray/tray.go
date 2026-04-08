package tray

import (
	"context"

	"github.com/getlantern/systray"
)

func Run(cancel context.CancelFunc, errChan <-chan error, doneChan <-chan struct{}) {
	systray.Run(onReady(cancel, errChan, doneChan), nil)
}
