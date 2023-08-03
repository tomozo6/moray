package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tomozo6/moray/pkg/aws"
	"github.com/tomozo6/moray/pkg/ui"
)

var ec2loginCmd = &cobra.Command{
	Use:   "ec2login",
	Short: "ec2login",
	Long:  "Login to EC2 securely using SSM Session Manager.",

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

		// ユーザーにEC2を選択させる
		ec2Names := e.GetInstanceNames()
		ec2Name, err := ui.AskEC2(ec2Names)
		if err != nil {
			log.Fatalln(err)
		}

		// 選択されたEC2名からインスタンスIDを取得する
		ec2Info, err := e.GetInstanceInfoFromName(ec2Name)
		if err != nil {
			log.Fatalln(err)
		}
		ec2ID := ec2Info.InstanceId

		// SSM Clientの生成
		ssmClient, _ := aws.MakeSSMSVC(&profile)

		// SSMを利用してEC2にログインをおこなう
		aws.SSMEC2Login(ssmClient, ec2ID, region, &profile)
	},
}

func init() {
	rootCmd.AddCommand(ec2loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ec2loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ec2loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
