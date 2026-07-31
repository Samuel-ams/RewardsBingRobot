package config

import (
	"errors"
	"os"
	"path/filepath"
)

// Build-time variables
var (
	BuildEdgePath       string
	BuildUserDataDir    string
	BuildTmpUserDataDir string
)

type Config struct {
	EdgePath       string
	UserDataDir    string
	TmpUserDataDir string
	LowSpeed       float64
	HighSpeed      float64
	TypeTick       int
	QtdSearches    int
}

var defaultConfig = &Config{
	// EdgePath is the path to the Edge executable
	EdgePath: "",
	// UserEdgeDir is the path to the Edge user data directory, where the browser profile is stored
	UserDataDir: "",
	// TmpUserEdgeDir is the path to the copy Edge user data directory
	TmpUserDataDir: "",
	// Setting Low and High mouse speed
	LowSpeed:  0.1,
	HighSpeed: 1.2,
	// TypeTick is the time in ms between each keystroke when typing
	TypeTick: 150,
	// Number of searches the robot will perform
	QtdSearches: 20,
}

func Load() (*Config, error) {
	if defaultConfig.EdgePath != "" && defaultConfig.UserDataDir != "" && defaultConfig.TmpUserDataDir != "" {
		return defaultConfig, nil
	}

	edgePath := firstNonEmpty(BuildEdgePath, os.Getenv("EDGE_PATH"))
	if edgePath == "" {
		return nil, errors.New("EDGE_PATH is not set")
	}

	userEdgeDir := firstNonEmpty(BuildUserDataDir, os.Getenv("USER_DATA_DIR"))
	if userEdgeDir == "" {
		return nil, errors.New("USER_DATA_DIR is not set")
	}

	tmpUserDataDir := firstNonEmpty(BuildTmpUserDataDir, os.Getenv("TMP_USER_DATA_DIR"))
	if tmpUserDataDir == "" {
		return nil, errors.New("TMP_USER_DATA_DIR is not set")
	}

	programFilesX86 := os.Getenv("ProgramFiles(x86)")
	defaultConfig.EdgePath = filepath.Join(programFilesX86, edgePath)

	localAppDataDir := os.Getenv("LOCALAPPDATA")
	defaultConfig.UserDataDir = filepath.Join(localAppDataDir, userEdgeDir)
	defaultConfig.TmpUserDataDir = filepath.Join(localAppDataDir, tmpUserDataDir)

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
