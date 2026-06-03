package main

import (
	"context"
	"fmt"
	"log/slog"
	rewardsrobot "rewardsAutomation/internal/rewardsRobot"
	"rewardsAutomation/internal/tray"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env relative to the executable so the built binary works from any directory
	_ = godotenv.Overload()

	beeep.AppName = "RewardsRobot"

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error, 1)
	doneChan := make(chan struct{}, 1)

	sendErr := func(err error) {
		select {
		case errChan <- err:
		default:
		}
	}

	go func() {
		defer func() {
			recovered := recover()
			if recovered != nil {
				sendErr(fmt.Errorf("panic in main worker: %v", recovered))
				cancel()
			}
		}()

		robot := rewardsrobot.New(ctx)
		err := robot.Run()
		if err != nil && err != context.Canceled {
			slog.Error("robot exited with error", "error", err)
			sendErr(err)
			cancel()
			return
		}

		cancel()
		doneChan <- struct{}{}
	}()

	tray.Run(cancel, errChan, doneChan)
}
