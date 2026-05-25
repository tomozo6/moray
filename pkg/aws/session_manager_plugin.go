package aws

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func resolvePluginProfile(profile string) string {
	if profile != "" {
		return profile
	}

	if envProfile := os.Getenv("AWS_PROFILE"); envProfile != "" {
		return envProfile
	}
	if envDefaultProfile := os.Getenv("AWS_DEFAULT_PROFILE"); envDefaultProfile != "" {
		return envDefaultProfile
	}
	return "default"
}

func resolveSessionManagerPluginCommand() (string, error) {
	bin := "session-manager-plugin"
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}

	if p, err := exec.LookPath(bin); err == nil {
		return p, nil
	}

	if runtime.GOOS == "windows" {
		programFiles := os.Getenv("ProgramFiles")
		if programFiles != "" {
			fallback := filepath.Join(programFiles, "Amazon", "SessionManagerPlugin", "bin", bin)
			if _, err := os.Stat(fallback); err == nil {
				return fallback, nil
			}
		}
	}

	return "", errors.New("session-manager-plugin not found in PATH")
}
