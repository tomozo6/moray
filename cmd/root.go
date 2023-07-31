package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomozo6/moray/pkg/aws"
	"github.com/tomozo6/moray/pkg/ui"
)

var (
	// build時に-ldflagsで上書きする
	Version = "unset"
)

var rootCmd = &cobra.Command{
	Use:     "moray",
	Short:   "Short moray",
	Long:    "Long moray",
	Version: Version,

	Run: func(cmd *cobra.Command, args []string) {

		// フラグの読み込み
		inputProfile, _ := cmd.PersistentFlags().GetString("profile")

		// profileが指定されている場合はその値をプロファイル名とする。profileが指定されていない場合は環境変数AWS_PROFILEからプロファイル名を取得。AWS_PROFILEが存在しない場合はdefaultとする
		var profile string
		if len(inputProfile) > 0 {
			profile = inputProfile

		} else {
			envProfile := os.Getenv("AWS_PROFILE")
			if len(envProfile) > 0 {
				profile = envProfile
			} else {
				profile = "default"
			}
		}

		// ec2Clientの生成及びリージョン名の取得
		ec2Client, region := aws.MakeEC2SVC(&profile)

		// EC2インスタンスの情報が詰まったオブジェクトを生成
		e, err := aws.NewEC2Instances(ec2Client)
		if err != nil {
			fmt.Println(err)
		}

		// ユーザーにBastionを選択させる
		ec2Names := e.GetInstanceNames()
		bastionName := ui.AskBastion(ec2Names)

		// 選択されたBastion名からインスタンスIDを取得する
		bastionInfo, err := e.GetInstanceInfoFromName(bastionName)
		if err != nil {
			fmt.Println(err)
		}
		bastionID := bastionInfo.InstanceId

		// RDSClientの生成
		rdsClient, _ := aws.MakeRDSSVC(&profile)

		// RDS&DocDBクラスタの情報が詰まったオブジェクトを生成
		d, _ := aws.NewDBClusters(rdsClient)

		// ユーザーに接続したいクラスタ名を選択させる
		dbNames := d.GetDBClusterNames()
		dbName := ui.AskDB(dbNames)

		// 選択されたクラスタ名から必要な情報を取得する
		dbInfo, err := d.GetDBClusterInfoFromName(dbName)
		if err != nil {
			fmt.Println(err)
		}

		dbHost := dbInfo.ReaderEndpoint
		dbPort := dbInfo.Port
		localPort := dbInfo.Port

		// SSM Clientの生成
		ssmClient, _ := aws.MakeSSMSVC(&profile)

		// SSMを利用してポートフォワーディングをおこなう
		aws.SSMPortForwarding(ssmClient, bastionID, dbHost, dbPort, localPort, region, &profile)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("profile", "p", "", "Use a specific profile from your credential file.")
}
