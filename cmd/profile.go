package cmd

import "os"

func resolveProfile(flagProfile string) string {
	if flagProfile != "" {
		return flagProfile
	}

	if envProfile := os.Getenv("AWS_PROFILE"); envProfile != "" {
		return envProfile
	}

	// AWS Tools for PowerShell often sets AWS_DEFAULT_PROFILE instead of AWS_PROFILE.
	if envDefaultProfile := os.Getenv("AWS_DEFAULT_PROFILE"); envDefaultProfile != "" {
		return envDefaultProfile
	}

	// Keep empty to let AWS SDK default credential chain resolve credentials
	// (env vars, shared config default profile, SSO cache, etc).
	return ""
}
