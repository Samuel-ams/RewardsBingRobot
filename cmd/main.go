package main

import (
	"rewardsAutomation/internal/tray"

	"github.com/gen2brain/beeep"
)

func main() {
	beeep.AppName = "RewardsRobot"

	tray.Run()
}
