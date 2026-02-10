package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	EdgePath    string
	UserEdgeDir string
	LowSpeed    float64
	HighSpeed   float64
	TypeTick    int
	QtdSearches int
}

var defaultConfig = &Config{
	// EdgePath is the path to the Edge executable
	EdgePath: "",
	// UserEdgeDir is the path to the Edge user data directory, where the browser profile is stored
	UserEdgeDir: "",
	// Setting Low and High mouse speed
	LowSpeed:  0.1,
	HighSpeed: 1.2,
	// TypeTick is the time in ms between each keystroke when typing
	TypeTick: 140,
	// Number of searches the robot will perform
	QtdSearches: 30,
}

func Load() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Edge user data directory is located in the user's home directory under "AppData\Local\Microsoft\Edge"
	homeDir = filepath.Join(homeDir, "AppData", "Local", "Microsoft", "Edge")

	defaultConfig.EdgePath = "C:\\Program Files (x86)\\Microsoft\\Edge\\Application\\msedge.exe"
	defaultConfig.UserEdgeDir = homeDir

	// Load the Edge executable path && user data directory from environment variables
	return defaultConfig, nil
}
