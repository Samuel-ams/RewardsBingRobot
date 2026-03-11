package main

import (
	"log/slog"
	"rewardsAutomation/internal/tray"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		return
	}

	beeep.AppName = "RewardsRobot"

	tray.Run()
}
