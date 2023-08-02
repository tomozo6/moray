package aws

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Instances []types.Instance

// EC2Instancesオブジェクトを作成します。
func NewEC2Instances(c *ec2.Client) (EC2Instances, error) {
	ctx := context.Background()
	input := &ec2.DescribeInstancesInput{}

	output, err := c.DescribeInstances(ctx, input)
	if err != nil {
		log.Fatalf("DescribeInstancesに失敗しました: , %v", err)
	}

	ec2s := EC2Instances{}
	for _, r := range output.Reservations {
		ec2s = append(ec2s, r.Instances...)
	}

	return ec2s, nil
}

// インスタンス名(NameタグのValue)のスライスを返します。
func (e *EC2Instances) GetInstanceNames() []string {
	names := []string{}

	for _, instance := range *e {
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				names = append(names, *tag.Value)
			}
		}
	}

	return names
}

// 合致する名前のインスタンス詳細情報を返します。
func (e *EC2Instances) GetInstanceInfoFromName(name *string) (types.Instance, error) {
	for _, instance := range *e {
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" {
				if *tag.Value == *name {
					return instance, nil
				}
			}
		}
	}
	return types.Instance{}, errors.New("該当のインスタンスが存在しません。")
}
