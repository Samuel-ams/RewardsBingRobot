package config

import (
	"errors"
	"os"
	"path/filepath"
)

// Build-time variables
var (
	BuildUserEdgeDir    string
	BuildTmpUserDataDir string
)

type Config struct {
	// EdgePath    string
	UserEdgeDir    string
	TmpUserDataDir string
	LowSpeed       float64
	HighSpeed      float64
	TypeTick       int
	QtdSearches    int
}

var defaultConfig = &Config{
	// EdgePath is the path to the Edge executable
	// EdgePath: "",
	// UserEdgeDir is the path to the Edge user data directory, where the browser profile is stored
	UserEdgeDir: "",
	// TmpUserEdgeDir is the path to the copy Edge user data directory
	TmpUserDataDir: "",
	// Setting Low and High mouse speed
	LowSpeed:  0.1,
	HighSpeed: 1.2,
	// TypeTick is the time in ms between each keystroke when typing
	TypeTick: 150,
	// Number of searches the robot will perform
	QtdSearches: 30,
}

func Load() (*Config, error) {
	// if defaultConfig.EdgePath != "" && defaultConfig.UserEdgeDir != "" {
	// 	return defaultConfig, nil
	// }

	if defaultConfig.UserEdgeDir != "" && defaultConfig.TmpUserDataDir != "" {
		return defaultConfig, nil
	}

	// homeDir, err := os.UserHomeDir()
	// if err != nil {
	// 	return nil, err
	// }

	// edgePath := firstNonEmpty(BuildEdgePath, os.Getenv("EDGE_PATH"))
	// if edgePath == "" {
	// 	return nil, errors.New("EDGE_PATH is not set")
	// }
	// edgePath = filepath.Clean(edgePath)

	userEdgeDir := firstNonEmpty(BuildUserEdgeDir, os.Getenv("USER_EDGE_DIR"))
	if userEdgeDir == "" {
		return nil, errors.New("USER_EDGE_DIR is not set")
	}

	tmpUserDataDir := firstNonEmpty(BuildTmpUserDataDir, os.Getenv("TMP_USER_DATA_DIR"))
	if tmpUserDataDir == "" {
		return nil, errors.New("TMP_USER_DATA_DIR is not set")
	}

	localAppDataDir := os.Getenv("LOCALAPPDATA")
	defaultConfig.UserEdgeDir = filepath.Join(localAppDataDir, userEdgeDir)
	defaultConfig.TmpUserDataDir = filepath.Join(localAppDataDir, tmpUserDataDir)

	// homeDir = filepath.Join(homeDir, userEdgeDir)
	// homeDir = filepath.Clean(homeDir)

	// Load config from environment variables
	// defaultConfig.EdgePath = edgePath
	// defaultConfig.UserEdgeDir = homeDir

	// slog.Info("Edge path:", "value", defaultConfig.EdgePath)
	// slog.Info("User edge directory:", "value", defaultConfig.UserEdgeDir)

	return defaultConfig, nil
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
