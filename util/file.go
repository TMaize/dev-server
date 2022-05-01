package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GetConfigFile(subPath string) (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(dir, ".dev-server")

	_, err = os.Stat(configDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}
		err = os.MkdirAll(configDir, os.ModeDir)
		if err != nil {
			return "", err
		}
	}

	configFile := FmtFilePath(filepath.Join(configDir, subPath))

	return configFile, nil
}

func FmtFilePath(p string) string {
	directory, _ := filepath.Abs(p)
	return strings.ReplaceAll(directory, "\\", "/")
}

func FmtFileSize(size int64) string {
	var pretty float64
	var unit string
	if size >= 1048576 {
		pretty = float64(size) / 1048576
		unit = "M"
	} else {
		pretty = float64(size) / 1024
		unit = "K"
	}

	return fmt.Sprintf("%.2f%s", pretty, unit)
}