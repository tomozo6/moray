package aws

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type DBClusters []types.DBCluster

func NewDBClusters(c *rds.Client) (DBClusters, error) {
	ctx := context.Background()
	input := &rds.DescribeDBClustersInput{}

	output, err := c.DescribeDBClusters(ctx, input)
	if err != nil {
		return nil, err
	}

	return output.DBClusters, nil
}

// GetDBClusterNamesはクラスタ名のスライスを返します。
func (d *DBClusters) GetDBClusterNames() []string {
	DBClusterNames := []string{}

	for _, dbcluster := range *d {
		DBClusterNames = append(DBClusterNames, *dbcluster.DBClusterIdentifier)
	}

	return DBClusterNames
}

// クラスタ名からそのクラスタの詳細情報を返します。
func (d *DBClusters) GetDBClusterInfoFromName(name string) (types.DBCluster, error) {
	for _, dbcluster := range *d {
		if *dbcluster.DBClusterIdentifier == name {
			return dbcluster, nil
		}
	}
	return types.DBCluster{}, errors.New("該当のクラスタが存在しません。")
}

// for _, cluster := range output.DBClusters {
// fmt.Println(*cluster.DBClusterIdentifier)
// fmt.Println(*cluster.Endpoint)
// fmt.Println(*cluster.ReaderEndpoint)
// fmt.Println(*cluster.Port)
// fmt.Println(*cluster.Engine)
// }
//
// return ec2InstancesInfo, nil
// }
