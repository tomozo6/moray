package cmd

import (
	"fmt"
	"log"
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
	Use:   "moray",
	Short: "moray",
	Long: `moray is a CLI tool
created to facilitate port forwarding operations
to remote hosts with the AWS System Manager session manager.`,
	Version: Version,

	Run: func(cmd *cobra.Command, args []string) {

		// フラグの読み込み
		inputProfile, _ := cmd.PersistentFlags().GetString("profile")
		witerFlag, _ := cmd.PersistentFlags().GetBool("writer")
		inputLocalPort, _ := cmd.PersistentFlags().GetInt32("port")

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
		bastionName, err := ui.AskBastion(ec2Names)
		if err != nil {
			log.Fatalln(err)
		}

		// 選択されたBastion名からインスタンスIDを取得する
		bastionInfo, err := e.GetInstanceInfoFromName(bastionName)
		if err != nil {
			log.Fatalln(err)
		}
		bastionID := bastionInfo.InstanceId

		// RDSClientの生成
		rdsClient, _ := aws.MakeRDSSVC(&profile)

		// RDS&DocDBクラスタの情報が詰まったオブジェクトを生成
		d, _ := aws.NewDBClusters(rdsClient)

		// ユーザーに接続したいクラスタ名を選択させる
		dbNames := d.GetDBClusterNames()
		dbName, err := ui.AskDB(dbNames)
		if err != nil {
			log.Fatalln(err)
		}

		// 選択されたクラスタ名から必要な情報を取得する
		dbInfo, err := d.GetDBClusterInfoFromName(dbName)
		if err != nil {
			fmt.Println(err)
		}

		// Writerインスタンスに接続するかReaderインスタンスに接続するか
		var dbHost *string
		if witerFlag {
			dbHost = dbInfo.Endpoint
		} else {
			dbHost = dbInfo.ReaderEndpoint
		}

		dbPort := dbInfo.Port

		var localPort *int32
		if inputLocalPort == 0 {
			localPort = dbInfo.Port
		} else {
			localPort = &inputLocalPort
		}

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
	rootCmd.PersistentFlags().String("profile", "", "Use a specific profile from your credential file.")
	rootCmd.PersistentFlags().BoolP("writer", "w", false, "Connect to the writer instance. (Connects to the reader instance if not specified.)")
	rootCmd.PersistentFlags().Int32P("port", "p", 0, "Local port number. (If not set, the same number as the port number of the connection destination DB will be set as the local port.)")
}
