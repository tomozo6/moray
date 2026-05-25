package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func MakeEC2SVC(profile *string) (*ec2.Client, *string) {
	cfg, err := loadAWSConfig(context.TODO(), *profile)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := ec2.NewFromConfig(cfg)
	return client, &cfg.Region
}
