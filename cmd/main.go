package main

import (
	"context"
	"fmt"
	rewardsrobot "rewardsAutomation/internal/rewardsRobot"
	"rewardsAutomation/internal/tray"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Overload()
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

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
			sendErr(err)
			cancel()
			return
		}

		cancel()
		doneChan <- struct{}{}
	}()

	tray.Run(cancel, errChan, doneChan)
}
