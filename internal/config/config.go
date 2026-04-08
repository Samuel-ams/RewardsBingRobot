package config

import (
	"errors"
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
	TypeTick: 150,
	// Number of searches the robot will perform
	QtdSearches: 30,
}

func Load() (*Config, error) {
	if defaultConfig.EdgePath != "" && defaultConfig.UserEdgeDir != "" {
		return defaultConfig, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	edgePath := os.Getenv("EDGE_PATH")
	if edgePath == "" {
		return nil, errors.New("EDGE_PATH is not set")
	}

	edgePath = filepath.Clean(edgePath)

	userEdgeDir := os.Getenv("USER_EDGE_DIR")
	if userEdgeDir == "" {
		return nil, errors.New("USER_EDGE_DIR is not set")
	}

	homeDir = filepath.Join(homeDir, userEdgeDir)
	homeDir = filepath.Clean(homeDir)

	// Load config from environment variables
	defaultConfig.EdgePath = edgePath
	defaultConfig.UserEdgeDir = homeDir

	return defaultConfig, nil
}
