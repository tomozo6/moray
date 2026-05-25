package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/rds"
)

func MakeRDSSVC(profile *string) (*rds.Client, *string) {
	cfg, err := loadAWSConfig(context.TODO(), *profile)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := rds.NewFromConfig(cfg)
	return client, &cfg.Region
}
