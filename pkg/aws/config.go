package aws

import "github.com/aws/aws-sdk-go-v2/config"

func loadConfigOptions(profile string) []func(*config.LoadOptions) error {
	if profile == "" {
		return nil
	}

	return []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(profile),
	}
}
