package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	OSGetwd = os.Getwd
	OSStat  = os.Stat
)

func GetEnvFilePath(envFile string) (string, error) {
	currentDir, err := OSGetwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(currentDir, envFile)
		if _, sErr := OSStat(envPath); sErr == nil {
			return envPath, nil
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			break
		}
		currentDir = parentDir
	}

	return "", fmt.Errorf("could not find envrionment file")
}
