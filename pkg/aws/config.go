package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	awssdk "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func loadConfigOptions(profile string) []func(*config.LoadOptions) error {
	if profile == "" {
		return nil
	}

	return []func(*config.LoadOptions) error{
		config.WithSharedConfigProfile(profile),
	}
}

type cliExportCredentials struct {
	AccessKeyID     string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken    string `json:"SessionToken"`
}

func loadAWSConfig(ctx context.Context, profile string) (awssdk.Config, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		loadConfigOptions(profile)...,
	)
	if err != nil {
		return cfg, err
	}

	if _, err = cfg.Credentials.Retrieve(ctx); err == nil {
		return cfg, nil
	}

	cliCreds, cliErr := exportCredentialsWithAWSCLI(ctx, profile)
	if cliErr != nil {
		return cfg, fmt.Errorf("failed to retrieve credentials with SDK (%v) and AWS CLI fallback (%v)", err, cliErr)
	}

	cfg.Credentials = awssdk.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
		cliCreds.AccessKeyID,
		cliCreds.SecretAccessKey,
		cliCreds.SessionToken,
	))

	return cfg, nil
}

func exportCredentialsWithAWSCLI(ctx context.Context, profile string) (*cliExportCredentials, error) {
	args := []string{"configure", "export-credentials", "--format", "process"}
	if profile != "" {
		args = append(args, "--profile", profile)
	}

	cmd := exec.CommandContext(ctx, "aws", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var cred cliExportCredentials
	if err := json.Unmarshal(out, &cred); err != nil {
		return nil, err
	}

	if strings.TrimSpace(cred.AccessKeyID) == "" || strings.TrimSpace(cred.SecretAccessKey) == "" {
		return nil, fmt.Errorf("AWS CLI returned empty credential fields")
	}

	return &cred, nil
}
